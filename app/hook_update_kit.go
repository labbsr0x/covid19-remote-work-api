package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func updateRemoteKitURL(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Fazendo a requisição de criação do kit para %s", storageAPIURL)

	params := mux.Vars(r)
	key := strings.ToLower(params["key"])

	if key == "" {
		logrus.Error("Error processing your request. Your request URL must have the attribute `key`")
		w.WriteHeader(400)
		w.Write([]byte("Erro ao processar sua requisição. O body do request precisa ter o atributo `key` especificado"))
	}

	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logrus.Errorf("Error decoding Body of request to remote work kit %s. Details: %s", storageAPIURL, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestBody["kitDownloadURL"] == "" {
		logrus.Error("Error processing your request. Your request body must have the attribute `kitDownloadURL`")
		w.WriteHeader(400)
		w.Write([]byte("Erro ao processar sua requisição. O body do request precisa ter o atributo `kitDownloadURL` especificado"))
	}
	requestBody["kitDownloadURL"] = strings.ToLower(requestBody["kitDownloadURL"])

	// Update the current User "kitDownloadURL"
	resp, err := registerRemoteWorkRequest(key, requestBody)
	bytes, err := json.Marshal(resp)

	if err != nil {
		logrus.Errorf("Error marshelling JSON response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Erro ao atualizar a requisição para kits no storage. Chave do usuário: %s. Error: %s", key, err)))
		return
	}

	w.WriteHeader(201)
	w.Write(bytes)

}
