package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func HttpPost(httpClient *http.Client, url string, header map[string]string, body interface{}, result interface{}) ([]byte, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return io.ReadAll(resp.Body)
}
func HttpGet(url string, header map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
