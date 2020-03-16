package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// User is someone that has access to VPN
type User struct {
	Key    string `json:"key,omitempty"`
	KitURL string `json:"kitDownloadURL,empty"`
}

var storageAPIURL string
var baseURL string
var fileServerVMUsersPath string

func main() {

	fileServerVMUsersPath = "/kit/users/"

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
