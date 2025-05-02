package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type LoginResponse struct {
	SessionID string `json:"sessionId"`
}

func Login(id string, token string, baseUrl string) (LoginResponse, error) {
	req := loginRequest{
		Id:    id,
		Token: token,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(req)
	if err != nil {
		return LoginResponse{}, err
	}

	resp, err := http.Post(baseUrl+"/login", "application/json", &body)
	if err != nil {
		return LoginResponse{}, err
	}
	defer resp.Body.Close()

	var loginResponse LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return LoginResponse{}, err
	}

	return loginResponse, nil
}
