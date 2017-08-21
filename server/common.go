package main

import (
	"net/http"
	"strings"
)

type queryRespData struct {
	AccessKeyID string `json:"access_key_id"`
	Signature   string `json:"signature"`
	Expires     int    `json:"expires"`
}

type headerRespData struct {
	Authorization string `json:"authorization"`
}

func addHeadersToRequest(request *http.Request, headersMap map[string]interface{}) (err error) {
	// Add headers to request.
	for key, item := range headersMap {
		value, ok := item.(string)
		if ok != true {
			continue
		}
		if value != "" {
			request.Header.Set(key, value)
		}
	}
	return nil
}

func generateStringQuery(queryMap map[string]interface{}) (string, error) {
	stringValue := ""
	for key, item := range queryMap {
		value, ok := item.(string)
		if ok != true {
			continue
		}
		queryItem := (key + "=" + value + "&")
		stringValue += queryItem
	}
	if stringValue == "" {
		return "", nil
	}

	return strings.Trim(stringValue, "&"), nil
}
