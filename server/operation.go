package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func operationQueryHandle(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Hanlding operation query request.")
	log.Println("Parsing the body of request.")

	// Parse body of signature request, and make a http request to sign.
	requestToSign, requestKeys, err := buildRequestToSign(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	headersInterface, ok := requestKeys["headers"]
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Empty headers.")
		return
	}
	headersMap, ok := headersInterface.(map[string]interface{})
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid headers.")
		return
	}

	expiresInterface, ok := requestKeys["expires"]
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Empty expires.")
		return
	}
	stringExpires, ok := expiresInterface.(string)
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid expires.")
		return
	}
	intExpires, err := strconv.Atoi(stringExpires)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	err = addHeadersToRequest(requestToSign, headersMap)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// Build string to sign.
	err = signer.WriteQuerySignature(requestToSign, intExpires)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// Respond to client.
	requestToSign.ParseForm()
	signature := requestToSign.Form.Get("signature")
	data := queryRespData{config.AccessKeyID, signature, intExpires}
	respData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.Write(respData)
	log.Println("Respond operation query request success.")
	return
}

func operationHeaderHandle(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Hanlding operation header request.")
	log.Println("Parsing the body of request.")

	// Parse body of signature request, and make a http request to sign.
	requestToSign, requestKeys, err := buildRequestToSign(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	headersInterface, ok := requestKeys["headers"]
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Empty headers.")
		return
	}

	headersMap, ok := headersInterface.(map[string]interface{})
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid headers.")
		return
	}

	err = addHeadersToRequest(requestToSign, headersMap)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// Sign the request to sign.
	stringToSign, err := signer.BuildStringToSign(requestToSign)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	h := hmac.New(sha256.New, []byte(config.SecretAccessKey))
	h.Write([]byte(stringToSign))
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	// Respond to client.
	authorization := "QS " + config.AccessKeyID + ":" + signature
	respDataJSON := headerRespData{authorization}
	respData, err := json.Marshal(respDataJSON)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.Write(respData)
	log.Println("Respond operation header request success.")
	return
}

func buildRequestToSign(request *http.Request) (*http.Request, map[string]interface{}, error) {
	requestBody, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Body of signature request:")
	log.Println(string(requestBody))

	var tempInterface interface{}
	err = json.Unmarshal(requestBody, &tempInterface)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	requestKeys, _ := tempInterface.(map[string]interface{})

	methodInterface, ok := requestKeys["method"]
	if ok != true {
		return nil, nil, errors.New("Empty method.")
	}
	methodString, ok := methodInterface.(string)
	if ok != true {
		return nil, nil, errors.New("Invalid method.")
	}

	hostInterface, ok := requestKeys["host"]
	if ok != true {
		return nil, nil, errors.New("Empty host.")
	}
	hostString, ok := hostInterface.(string)
	if ok != true {
		return nil, nil, errors.New("Invalid host.")
	}

	portInterface, ok := requestKeys["port"]
	if ok != true {
		return nil, nil, errors.New("Empty port.")
	}
	portString, ok := portInterface.(string)
	if ok != true {
		return nil, nil, errors.New("Invalid port.")
	}

	pathInterface, ok := requestKeys["path"]
	if ok != true {
		return nil, nil, errors.New("Empty path.")
	}
	pathString, ok := pathInterface.(string)
	if ok != true {
		return nil, nil, errors.New("Invalid path.")
	}

	queryString := ""
	queryInterface, ok := requestKeys["query"]
	if ok {
		queryMap, ok := queryInterface.(map[string]interface{})
		if ok != true {
			return nil, nil, errors.New("Invalid query.")
		}
		queryString, err = generateStringQuery(queryMap)
		if err != nil {
			return nil, nil, err
		}
	}

	protocolInterface, ok := requestKeys["protocol"]
	if ok != true {
		return nil, nil, errors.New("Empty protocol.")
	}
	protocolString, ok := protocolInterface.(string)
	if ok != true {
		return nil, nil, errors.New("Invalid protocol.")
	}

	uriString := protocolString + "://" + hostString + ":" + portString + pathString
	if queryString != "" {
		uriString += "?" + queryString
	}

	requestToSign, err := http.NewRequest(methodString, uriString, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}

	log.Println("request to sign:")
	log.Println(requestToSign)

	return requestToSign, requestKeys, nil
}
