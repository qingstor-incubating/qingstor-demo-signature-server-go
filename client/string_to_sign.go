package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yunify/qingstor-sdk-go/request"
)

type stringToSignQueryObject struct {
	StringToSign string `json:"string_to_sign"`
	Expires      int    `json:"expires"`
}

type stringToSignHeaderObject struct {
	StringToSign string `json:"string_to_sign"`
}

func stringToSignQuery() {
	getObjectRequest, _, err := bucket.GetObjectRequest("put-test-file", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	expires := int(time.Now().Unix() + interval)
	// Build the body of a signature request.
	requestBody, err := buildRequestBody(getObjectRequest, "query", expires)
	if err != nil {
		log.Println(err.Error())
		return
	}

	signatureRequest := buildSignatureRequest("/string-to-sign/query", requestBody)
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

	// Add query parameters to getobject request to QingStor Object Storage.
	getObjectRequest.ApplyQuerySignature(responseJSON.AccessKeyID, responseJSON.Expires, responseJSON.Signature)

	testSignature(getObjectRequest, "Get Object")
	return
}

func stringToSignHeader() {
	deleteObjectRequest, _, err := bucket.DeleteObjectRequest("signature-test-file")
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Build the body of signature request.
	requestBody, err := buildRequestBody(deleteObjectRequest, "header", 0)
	if err != nil {
		log.Println(err.Error())
		return
	}

	signatureRequest := buildSignatureRequest("/string-to-sign/header", requestBody)
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

	// Add authorization to header of deleteobject request to QingStor Object Storage.
	err = addAuthToHeader(deleteObjectRequest, signatureResponse)
	if err != nil {
		log.Println(err.Error())
		return
	}

	testSignature(deleteObjectRequest, "Delete Object")
	return
}

func buildRequestBody(handleObjectRequest *request.Request, signType string, expires int) (string, error) {
	stringToSign, err := buildStringToSign(handleObjectRequest, signType, expires)
	if err != nil {
		return "", err
	}

	var requestBodyByte []byte
	statusBool := false
	switch signType {
	case "query":
		requestBodyJSON := stringToSignQueryObject{stringToSign, expires}
		requestBodyByte, err = json.Marshal(requestBodyJSON)
	case "header":
		requestBodyJSON := stringToSignHeaderObject{stringToSign}
		requestBodyByte, err = json.Marshal(requestBodyJSON)
	default:
		statusBool = true
	}

	if (err != nil) || statusBool {
		return "", err
	}

	return string(requestBodyByte), nil
}
