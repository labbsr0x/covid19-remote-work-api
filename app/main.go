package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var remoteKitBuilderAPI string

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	remoteKitBuilderAPI = *flag.String("remoteKitBuilderAPI", "http://localhost:3000", "Virtual Machine Image Builder URL")
	logrus.Infof("Starting Remote Kit Builder API")
	logrus.Infof("remoteKitBuilderAPI=%s", remoteKitBuilderAPI)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.HandleFunc("/kit", requestKit).Methods("POST")
	router.HandleFunc("/kit/{id}", getKit).Methods("GET")
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
	logrus.Debugf("Fazendo a requisição de criação do kit para %s", remoteKitBuilderAPI)
	resp, err := http.Post(remoteKitBuilderAPI, "application/json", r.Body)
	if err != nil {
		logrus.Errorf("Error requesting remote work kit to %s. Details: %s", remoteKitBuilderAPI, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		logrus.Errorf("Error requesting kit creation. Code: %s", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte("Erro na requisição de criação do kit."))

	} else {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(201)
			w.Write(result)
		} else {
			w.WriteHeader(500)
			w.Write([]byte("Erro na conversão da resposta da requisição de criação do kit"))
		}
	}
}

func getKit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	getRemoteKitURL := remoteKitBuilderAPI + "/" + id

	logrus.Debugf("Fazendo a requisição de recuperação do kit para %s", getRemoteKitURL)
	resp, err := http.Get(getRemoteKitURL)
	if err != nil {
		logrus.Errorf("Error getting remote work kit to %s. Details: %s", getRemoteKitURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		logrus.Errorf("Error getting kit. Code: %s", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte("Erro na requisição de recuperação do kit."))

	} else {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(201)
			w.Write(result)
		} else {
			w.WriteHeader(500)
			w.Write([]byte("Erro na conversão da resposta da requisição de criação do kit"))
		}
	}
}
