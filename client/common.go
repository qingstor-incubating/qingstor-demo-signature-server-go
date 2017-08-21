package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/pengsrc/go-shared/rest"
	"github.com/yunify/qingstor-sdk-go/request"
)

type queryResp struct {
	AccessKeyID string `json:"access_key_id"`
	Signature   string `json:"signature"`
	Expires     int    `json:"expires"`
}

type headerResp struct {
	Authorization string `json:"authorization"`
}

func buildStringToSign(request *request.Request, signType string, expires int) (stringToSign string, err error) {
	objectHTTPRequest, err := builder.BuildHTTPRequest(request.Operation, request.Input)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	switch signType {
	case "header":
		stringToSign, err = signer.BuildStringToSign(objectHTTPRequest)
	case "query":
		stringToSign, err = signer.BuildQueryStringToSign(objectHTTPRequest, expires)

	default:
		return "", errors.New("Unknow sign type.")
	}

	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return stringToSign, nil
}

func buildSignatureRequest(requestPath string, requestBody string) rest.Request {
	hostMap, paramMap := buildParametersOfRequest()
	baseURL := serverURL + requestPath
	return rest.Request{Method: rest.Post, BaseURL: baseURL, Headers: hostMap, QueryParams: paramMap, Body: []byte(requestBody)}
}

func buildParametersOfRequest() (hostMap map[string]string, paramMap map[string]string) {
	hostMap = make(map[string]string)
	hostMap["Host"] = config.SignatureServerHost + ":" + config.SignatureServerPort
	paramMap = make(map[string]string)
	return
}

func addAuthToHeader(request *request.Request, resp *rest.Response) error {
	var responseJSON headerResp
	err := json.Unmarshal([]byte(resp.Body), &responseJSON)
	if err != nil {
		return err
	}

	request.ApplySignature(responseJSON.Authorization)
	return nil
}

func testSignature(handleRequest *request.Request, handleString string) {
	err := handleRequest.Do()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(handleString + ": " + handleRequest.HTTPResponse.Status)
}

func createTempFile(fileName string) (*os.File, error) {
	tempDir := os.TempDir()
	tempFilePath := tempDir + "/" + fileName
	file, err := os.Create(tempFilePath)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	stringContent := "put-test-file-content"
	_, err = file.Write([]byte(stringContent))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return file, nil
}
