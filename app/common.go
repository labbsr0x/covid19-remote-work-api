package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func registerRemoteWorkRequest(key string, requestBody map[string]string) (map[string]string, error) {

	client := &http.Client{}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		logrus.Errorf("Error marshelling JSON response: %s", err)
		return nil, fmt.Errorf("Erro ao atualizar URL para kit. Chave do usuário: %s. Error: %s", key, err)
	}

	remoteKitRequestURL := storageAPIURL + fileServerVMUsersPath + key
	req, err := http.NewRequest("PUT", remoteKitRequestURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		logrus.Errorf("Error preparing remote work kit request to %s. Details: %s", storageAPIURL, err)
		return nil, fmt.Errorf("Erro ao atualizar a URL do kit")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error requesting file server. Server: %s. Details: %s", storageAPIURL, err)
		return nil, fmt.Errorf("Erro ao atualizar a URL do kit")
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		logrus.Errorf("Error requesting kit creation. Code: %d", resp.StatusCode)
		return nil, fmt.Errorf("Erro na requisição de criação do kit %s", key)

	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Erro na conversão da resposta da requisição de criação do kit. Details: %s", err)
		return nil, fmt.Errorf("Erro ao atualizar a URL do kit")
	}

	kitURL := baseURL + "/kit/" + strings.ToLower(key)
	ret := map[string]string{
		"key":            key,
		"kitURL":         kitURL,
		"kitDownloadURL": string(result),
	}

	return ret, nil
}
