package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string, method string, body io.Reader) (interface{}, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to form request: %v", err))
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to send request: %v", err))
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response body: %v", err))
	}
	var in map[string]interface{}
	err = json.Unmarshal(responseBytes, &in)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse json response: %v", err))
	}

	if in["status-code"] == nil {
		return nil, errors.New(fmt.Sprintf("malformed response from request"))
	}

	if in["data"] == nil && method == "GET" {
		return nil, errors.New(fmt.Sprintf("malformed response from GET request"))
	}

	if in["status-code"].(float64) >= 400 {
		return nil, errors.New(
			fmt.Sprintf("code: %v - got an error response from server - error: %v", in["status-code"], in["data"]))
	}

	return in["data"], nil
}
