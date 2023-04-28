package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func InvokeMethodOnUrl(method, url string, headers map[string]string, data interface{}) (interface{}, error) {
	method = strings.ToUpper(method)

	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respBody interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, errors.Wrapf(err, "error decoding response body")
	}
	return respBody, nil
}
