package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func getKit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := strings.ToLower(params["key"])

	user, err := getKitRequest(key)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Erro ao recuperar o kit para a chave: %s. Details: %s", key, err)))
		return
	}

	if user.Key == "" {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Kit not found for user %s", key)))
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)

}

func getKitRequest(key string) (User, error) {
	getRemoteKitURL := storageAPIURL + fileServerVMUsersPath + strings.ToLower(key)

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
		logrus.Errorf("Error getting kit. Code: %d", resp.StatusCode)
		return User{}, nil
	}

	var retMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&retMap)
	logrus.Debugf("Received from file-server: %s", retMap)
	if err != nil {
		logrus.Errorf("Error decoding Body of request to remote work kit %s. Details: %s", storageAPIURL, err)
		return User{}, err
	}

	return User{
		Key:    key,
		KitURL: retMap["kitDownloadURL"],
	}, nil
}
