package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// User is someone that has access to VPN
type User struct {
	Key      string `json:"key,omitempty"`
	ImageURL string `json:"imageURL,empty"`
}

var storageAPIURL string
var baseURL string
var fileServerVMMetadataPath string

func main() {

	fileServerVMMetadataPath = "/kits/vm/meta/"

	logrus.Infof("Starting Remote Kit Builder API")
	logrus.SetLevel(logrus.DebugLevel)
	flag.StringVar(&storageAPIURL, "storageAPIURL", "http://localhost:3000", "Simple File storage URL to store Virtual Machine images and configurations")
	flag.StringVar(&baseURL, "baseURL", "http://localhost:8000", "Base URL of this service, used to print responses")
	flag.Parse()

	logrus.Infof("storageAPIURL=%s", storageAPIURL)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.HandleFunc("/kit", requestKit).Methods("POST")
	router.HandleFunc("/kit/{key}", updateRemoteKitURL).Methods("PUT")
	router.HandleFunc("/kit/{key}", getKit).Methods("GET")
	logrus.Infof("Banco do Brasil COVID-19 Remote Work Kit API up and running!")

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	logrus.Fatal(srv.ListenAndServe())
}

func requestKit(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Fazendo a requisição de criação do kit para %s", storageAPIURL)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logrus.Errorf("Error decoding Body of request to remote work kit %s. Details: %s", storageAPIURL, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Key = strings.ToLower(user.Key)

	// checks whether a kit has been already issued to the user
	fetchedUser, err := getKitRequest(user.Key)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Erro na checagem de usuário já existente. Chave do usuário: %s. Error: %s", user.Key, err)))
		return
	}

	if fetchedUser.Key != "" {
		w.WriteHeader(409)
		w.Write([]byte(fmt.Sprintf("Erro ao requisitar criação do kit: kit já existente para usuario %s. ", fetchedUser.Key)))
		return
	}

	// Invoking Conductor to create the kit
	downloadURL := baseURL + "/kit/" + strings.ToLower(user.Key)
	resp := map[string]string{
		"key":         user.Key,
		"downloadURL": downloadURL,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		logrus.Errorf("Error marshelling JSON response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Erro ao criar requisição para kit. Chave do usuário: %s. Error: %s", user.Key, err)))
		return
	}

	w.WriteHeader(201)
	w.Write(bytes)
}

func updateRemoteKitURL(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Fazendo a requisição de criação do kit para %s", storageAPIURL)
	client := &http.Client{}

	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logrus.Errorf("Error decoding Body of request to remote work kit %s. Details: %s", storageAPIURL, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestBody["key"] = strings.ToLower(requestBody["key"])

	json, err := json.Marshal(requestBody)
	if err != nil {
		logrus.Errorf("Error marshelling JSON response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Erro ao atualizar URL para kit. Chave do usuário: %s. Error: %s", requestBody["key"], err)))
	}

	remoteKitRequestURL := storageAPIURL + fileServerVMMetadataPath + requestBody["key"]
	req, err := http.NewRequest("PUT", remoteKitRequestURL, bytes.NewBuffer(json))
	if err != nil {
		logrus.Errorf("Error preparing remote work kit request to %s. Details: %s", storageAPIURL, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error requesting file server. Server: %s. Details: %s", storageAPIURL, err)
		w.WriteHeader(500)
		w.Write([]byte("Erro ao atualizar a URL do kit"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		logrus.Errorf("Error requesting kit creation. Code: %s", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte("Erro na requisição de criação do kit."))

	} else {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500)
			logrus.Errorf("Erro na conversão da resposta da requisição de criação do kit. Details: %s", err)
			w.Write([]byte("Erro ao atualizar a URL do kit"))
			return
		}

		w.WriteHeader(201)
		w.Write(result)
	}
}

func getKit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := strings.ToLower(params["id"])

	user, err := getKitRequest(id)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Erro ao recuperar o kit para o id: %s. Details: %s", id, err)))
		return
	}

	if user.Key == "" {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Kit not found for user %s", id)))
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)

}

func getKitRequest(key string) (User, error) {
	getRemoteKitURL := storageAPIURL + fileServerVMMetadataPath + strings.ToLower(key)

	logrus.Debugf("Fazendo a requisição de recuperação do kit para %s", getRemoteKitURL)
	resp, err := http.Get(getRemoteKitURL)
	if err != nil {
		logrus.Errorf("Error getting remote work kit to %s. Details: %s", getRemoteKitURL, err)
		return User{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return User{}, nil
	}
	if resp.StatusCode > 299 {
		logrus.Errorf("Error getting kit. Code: %s", resp.StatusCode)
		var user User
		err := json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			logrus.Errorf("Error decoding Body of request to remote work kit %s. Details: %s", storageAPIURL, err)
			return User{}, err
		}

		return user, nil

	}

	logrus.Debugf("User not found with key: %s", key)
	return User{}, nil
}
