package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

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

	// 1 - checks whether a kit has been already issued to the user
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

	// 2 - Registers the remote work kit request made
	requestBody := map[string]string{
		"kitDownloadURL": "",
	}
	_, err = registerRemoteWorkRequest(user.Key, requestBody)

	// 3 - Invoking services that issues the certificates

	// 4 - Retrieving the users roles from Corporate Services

	// 5 - Update the users roles on the OpenVPN Servers

	// 6 - Invoking Conductor to create the kit

	// 7 - Building the kitDownloadURL
	kitDownloadURL := baseURL + "/kit/" + strings.ToLower(user.Key)
	resp := map[string]string{
		"key":            user.Key,
		"kitDownloadURL": kitDownloadURL,
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
