package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	_ "strconv"
	"strings"
)

func stringToSignQueryHandle(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Hanlding string_to_sign query request.")
	log.Println("Parsing the body of request.")
	log.Println(request)
	// Parse body of signature request, and get string to sign.
	stringToSign, requestKeys, err := getStringToSign(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// Generate a signature.
	h := hmac.New(sha256.New, []byte(config.SecretAccessKey))
	h.Write([]byte(stringToSign))
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	expiresInterface, ok := requestKeys["expires"]
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Empty expires.")
		return
	}

	expires, ok := expiresInterface.(float64)
	if ok != true {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid expires.")
		return
	}

	// Respond to client.
	data := queryRespData{config.AccessKeyID, signature, int(expires)}
	respData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
	} else {
		w.Write(respData)
	}
	log.Println("Respond string_to_sign query request success.")
}

func stringToSignHeaderHandle(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Hanlding string_to_sign header request.")
	log.Println("Parsing the body of request.")

	// Parse body of signature request, and get string to sign.
	stringToSign, _, err := getStringToSign(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// Generate a signature.
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
	log.Println("Respond string_to_sign header request success.")
}

func getStringToSign(request *http.Request) (string, map[string]interface{}, error) {
	requestBody, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return "", nil, err
	}
	log.Println("Body of signature request:")
	log.Println(string(requestBody))

	var tempInterface interface{}
	err = json.Unmarshal(requestBody, &tempInterface)
	if err != nil {
		log.Println(err.Error())
		return "", nil, err
	}

	requestKeys, ok := tempInterface.(map[string]interface{})
	if ok != true {
		log.Println("Type error.")
		return "", nil, errors.New("Empty string_to_sign.")
	}

	stringToSignInterface, ok := requestKeys["string_to_sign"]
	if ok == false {
		log.Println("Empty string to sign.")
		return "", nil, errors.New("Empty string to sign.")
	}
	stringToSign, ok := stringToSignInterface.(string)
	if ok != true {
		log.Println("Invalid string to sign.")
		return "", nil, errors.New("Invalid string to sign.")
	}

	log.Println("string to sign:")
	log.Println(stringToSign)

	return stringToSign, requestKeys, nil
}
