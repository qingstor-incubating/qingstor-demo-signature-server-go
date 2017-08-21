package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yunify/qingstor-sdk-go/request"
	qs "github.com/yunify/qingstor-sdk-go/service"
)

func operationQuery() {
	listObjectRequest, _, err := bucket.ListObjectsRequest(nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Build the body of signature request.
	requestBody, err := buildOperationRequestBody(listObjectRequest, "query")
	if err != nil {
		log.Println(err.Error())
		return
	}
	signatureRequest := buildSignatureRequest("/operation/query", requestBody)
	signatureResponse, err := client.API(&signatureRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if signatureResponse == nil {
		log.Println("Empty response.")
		return
	}

	if signatureResponse.StatusCode != 200 {
		log.Println("Response status code:")
		log.Println(signatureResponse.StatusCode)
		return
	}

	var responseJSON queryResp
	err = json.Unmarshal([]byte(signatureResponse.Body), &responseJSON)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Add query parameters to listobject request to QingStor Object Storage.
	listObjectRequest.ApplyQuerySignature(responseJSON.AccessKeyID, responseJSON.Expires, responseJSON.Signature)

	testSignature(listObjectRequest, "List Object")
	return
}

func operationHeader() {
	// Create a temp file to put.
	fileName := "put-test-file"
	file, err := createTempFile(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	putObjectRequest, _, err := bucket.PutObjectRequest(fileName, &qs.PutObjectInput{Body: file})
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Build the body of signature request.
	requestBody, err := buildOperationRequestBody(putObjectRequest, "header")
	if err != nil {
		return
	}

	signatureRequest := buildSignatureRequest("/operation/header", requestBody)
	signatureResponse, err := client.API(&signatureRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if signatureResponse == nil {
		log.Println("Empty response.")
		return
	}
	if signatureResponse.StatusCode != 200 {
		log.Println("Response status code:")
		log.Println(signatureResponse.StatusCode)
		return
	}

	// Add authorization to header of putobject request to QingStor Object Storage.
	err = addAuthToHeader(putObjectRequest, signatureResponse)
	if err != nil {
		log.Println(err.Error())
		return
	}

	testSignature(putObjectRequest, "Put Object")
	return
}

func buildOperationRequestBody(handleObjectRequest *request.Request, signType string) (string, error) {
	err := handleObjectRequest.Build()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	if signType == "query" {
		queryValues := handleObjectRequest.HTTPRequest.URL.Query()
		queryValues.Set("prefix", "test")
		handleObjectRequest.HTTPRequest.URL.RawQuery = queryValues.Encode()
	}
	requestBodyByte, err := makeOperationRequestBody(handleObjectRequest.HTTPRequest, signType)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return string(requestBodyByte), nil
}

func makeOperationRequestBody(request *http.Request, signType string) (requestBodyByte []byte, err error) {
	headerKeys := []string{"Date", "Content-Type", "Content-Length", "User-Agent"}
	headerMap := make(map[string]string)
	for _, headerKey := range headerKeys {
		headerValue := request.Header.Get(headerKey)
		if "" != headerValue {
			headerMap[headerKey] = headerValue
		}
	}

	requestMap := make(map[string]interface{})
	requestMap["method"] = request.Method
	requestMap["host"] = request.URL.Hostname()
	requestMap["port"] = request.URL.Port()
	requestMap["path"] = request.URL.Path
	requestMap["protocol"] = request.URL.Scheme
	requestMap["headers"] = headerMap

	queryValues := request.URL.Query()
	queryMap := make(map[string]string)
	for key, value := range queryValues {
		queryMap[key] = value[0]
	}

	if len(queryMap) != 0 {
		requestMap["query"] = queryMap
	}

	switch signType {
	case "query":
		expires := int(time.Now().Unix() + interval)
		requestMap["expires"] = strconv.Itoa(expires)
		requestBodyByte, err = json.Marshal(requestMap)
	case "header":
		requestBodyByte, err = json.Marshal(requestMap)
	default:
		return nil, errors.New("Unkown type.")
	}
	if err != nil {
		return nil, err
	}

	return requestBodyByte, nil
}
