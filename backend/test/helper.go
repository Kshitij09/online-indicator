package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func CreateRequest(method, target string, req interface{}) (*http.Request, error) {
	serialized, err := json.Marshal(req)
	if err != nil {
		return nil, nil
	}
	body := bytes.NewBuffer(serialized)
	httpReq := httptest.NewRequest(method, target, body)
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}
