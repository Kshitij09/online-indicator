package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type pingRequest struct {
	SessionId string `json:"sessionId"`
}

func Ping(sessionId string, baseUrl string) error {
	req := pingRequest{
		SessionId: sessionId,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseUrl+"/ping", "application/json", &body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
