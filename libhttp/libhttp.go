package libhttp

import (
	"encoding/json"
	"io"

	"net/http"
	"strings"

	"github.com/helloferdie/golib/liblogger"
)

// Request - request HTTP and expect response in JSON map[string]interface{}
func Request(address string, method string, payloadData map[string]interface{}, headerData map[string]string) (map[string]interface{}, int, error) {
	payloadBytes, _ := json.Marshal(payloadData)
	payload := strings.NewReader(string(payloadBytes))
	request, err := http.NewRequest(method, address, payload)
	request.Header.Add("Content-Type", "application/json")
	for k, v := range headerData {
		request.Header.Add(k, v)
	}
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return nil, 0, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return nil, 0, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return nil, response.StatusCode, err
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(responseBody, &responseJSON)
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return nil, response.StatusCode, err
	}
	return responseJSON, response.StatusCode, nil
}

// RequestRaw - request HTTP and expect response in raw string
func RequestRaw(address string, method string, payloadData map[string]interface{}, headerData map[string]string) (string, int, error) {
	payloadBytes, _ := json.Marshal(payloadData)
	payload := strings.NewReader(string(payloadBytes))
	request, err := http.NewRequest(method, address, payload)
	request.Header.Add("Content-Type", "application/json")
	for k, v := range headerData {
		request.Header.Add(k, v)
	}
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return "", 0, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return "", 0, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err)
		return "", response.StatusCode, err
	}

	return string(responseBody), response.StatusCode, err
}
