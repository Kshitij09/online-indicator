package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type registerRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	ApiKey string `json:"apiKey"`
	Id     string `json:"id"`
}

func Register(name string, baseUrl string) (RegisterResponse, error) {
	req := registerRequest{
		Name: name,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(req)
	if err != nil {
		return RegisterResponse{}, err
	}

	resp, err := http.Post(baseUrl+"/register", "application/json", &body)
	if err != nil {
		return RegisterResponse{}, err
	}
	defer resp.Body.Close()

	var registerResponse RegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	if err != nil {
		return RegisterResponse{}, err
	}

	return registerResponse, nil
}
