package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/pengsrc/go-shared/rest"
	qscfg "github.com/yunify/qingstor-sdk-go/config"
	qsbld "github.com/yunify/qingstor-sdk-go/request/builder"
	qss "github.com/yunify/qingstor-sdk-go/request/signer"
	qs "github.com/yunify/qingstor-sdk-go/service"
	"gopkg.in/yaml.v2"
)

type configInfo struct {
	Host                    string `yaml:"host"`
	Port                    string `yaml:"port"`
	Protocol                string `yaml:"protocol"`
	Zone                    string `yaml:"zone"`
	BucketName              string `yaml:"bucket_name"`
	SignatureServerHost     string `yaml:"signature_server_host"`
	SignatureServerPort     string `yaml:"signature_server_port"`
	SignatureServerProtocol string `yaml:"signature_server_protocol"`
}

// Global localClient.
var client rest.Client
var interval int64
var serverURL string
var bucket *qs.Bucket
var builder *qsbld.QingStorBuilder
var config configInfo
var signer qss.QingStorSigner

func configInit() (err error) {
	confFileData, err := ioutil.ReadFile("./config_client.yaml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(confFileData, &config)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return
	}
	qsConfig, err := qscfg.NewDefault()
	qsConfig.AccessKeyID = "EXAMPLE_ACCESS_KEY_ID"
	qsConfig.SecretAccessKey = "EXAMPLE_SECRET_ACCESS_KEY"
	qsConfig.Host = config.Host
	qsConfig.Port = port
	qsConfig.Protocol = config.Protocol
	service, err := qs.Init(qsConfig)
	if err != nil {
		return
	}

	bucket, err = service.Bucket(config.BucketName, config.Zone)
	if bucket == nil {
		log.Println("Bucket is not initilized.")
	}

	if err != nil {
		return
	}

	builder = &qsbld.QingStorBuilder{}
	client.HTTPClient = http.DefaultClient
	interval = 120

	serverURL = config.SignatureServerProtocol + "://" + config.SignatureServerHost + ":" + config.SignatureServerPort

	signer = qss.QingStorSigner{AccessKeyID: "EXAMPLE_ACCESS_KEY_ID", SecretAccessKey: "EXAMPLE_SECRET_ACCESS_KEY"}

	log.Println("Configure Success.")
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

	operationQuery()
	operationHeader()
	stringToSignQuery()
	stringToSignHeader()

	return
}
