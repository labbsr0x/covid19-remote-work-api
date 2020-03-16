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

// storageAPIURL is the URL of a storage provider used to persist the certificates and other information
var storageAPIURL string

// certificateIssuerAPIURL is the URL of the service that issues a user certificate
var certificateIssuerAPIURL string

// userRolesAPIURL is the URL of the service used to fetch users roles
var userRolesAPIURL string

// conductorAPIURL is the URL of the service that triggers and manages a conductor used to build and store the virtual machines
var conductorAPIURL string

// baseURL is this service base URL. Used to build output and return URLs
var baseURL string

// fileServerVMUsersPath is the path in the file server to persist
var fileServerVMUsersPath string

func main() {

	fileServerVMUsersPath = "/kit/users/"

	logrus.Infof("Starting Remote Kit Builder API")
	logrus.SetLevel(logrus.DebugLevel)
	flag.StringVar(&storageAPIURL, "storageAPIURL", "http://localhost:3000", "Simple File storage URL to store Virtual Machine images and configurations")
	flag.StringVar(&certificateIssuerAPIURL, "certificateIssuerAPIURL", "http://localhost:3005", "Certificate issuer service URL that will issue VPN certificates for users")
	flag.StringVar(&userRolesAPIURL, "userRolesAPIURL", "http://localhost:3006", "URL of the service used to fetch users roles")
	flag.StringVar(&conductorAPIURL, "conductorAPIURL", "http://localhost:3007", "URL of the service that triggers and manages a conductor used to build and store the virtual machines")
	flag.StringVar(&baseURL, "baseURL", "http://localhost:8000", "Base URL of this service, used to print responses")
	flag.Parse()

	logrus.Infof("storageAPIURL=%s", storageAPIURL)
	logrus.Infof("certificateIssuerAPIURL=%s", certificateIssuerAPIURL)
	logrus.Infof("userRolesAPIURL=%s", userRolesAPIURL)
	logrus.Infof("conductorAPIURL=%s", conductorAPIURL)
	logrus.Infof("baseURL=%s", baseURL)

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
