package main

import (
	"io/ioutil"
	"log"
	"net/http"

	qss "github.com/yunify/qingstor-sdk-go/request/signer"
	"gopkg.in/yaml.v2"
)

type configInfo struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
}

var config configInfo
var signer *qss.QingStorSigner

func configInit() (err error) {
	confFileData, err := ioutil.ReadFile("./config_server.yaml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(confFileData, &config)
	if err != nil {
		return
	}

	signer = &qss.QingStorSigner{AccessKeyID: config.AccessKeyID, SecretAccessKey: config.SecretAccessKey}
	return
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Configure global variables.
	err := configInit()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Configure server success.")

	// Handle signature requests.
	http.HandleFunc("/operation/query", operationQueryHandle)
	http.HandleFunc("/operation/header", operationHeaderHandle)
	http.HandleFunc("/string-to-sign/query", stringToSignQueryHandle)
	http.HandleFunc("/string-to-sign/header", stringToSignHeaderHandle)

	log.Println("Server is running...")
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, nil))

	return
}
