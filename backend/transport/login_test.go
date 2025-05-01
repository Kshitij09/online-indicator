package transport

import (
	"bytes"
	"encoding/json"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler_Success(t *testing.T) {
	body := LoginRequest{Name: "test", Token: "123"}
	req, err := createLoginRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: body.Token}
	fakeClock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGenerator, staticGenerator, fakeClock)
	handler := NewHttpHandler(LoginHandler(storage))

	existing := domain.Account{Name: body.Name}
	_, err = storage.Auth().Create(existing)
	if err != nil {
		t.Error(err)
	}

	handler(recorder, req)

	result := recorder.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("request failed: expected %d, got %d", http.StatusOK, result.StatusCode)
	}
	var resp LoginResponse
	err = json.NewDecoder(result.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.SessionId != staticGenerator.StubValue {
		t.Errorf("session incorrect, expected %s, got %s", staticGenerator.StubValue, resp.SessionId)
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	body := LoginRequest{Name: "test", Token: "123"}
	req, err := createLoginRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: "random"}
	fakeClock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGenerator, staticGenerator, fakeClock)
	handler := NewHttpHandler(LoginHandler(storage))

	existing := domain.Account{Name: body.Name}
	_, err = storage.Auth().Create(existing)
	if err != nil {
		t.Error(err)
	}

	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusUnauthorized
	if result.StatusCode != expectedStatusCode {
		t.Errorf("request failed: expected %d, got %d", expectedStatusCode, result.StatusCode)
	}
}

func TestLoginHandler_AccountNotFound(t *testing.T) {
	body := LoginRequest{Name: "test"}
	req, err := createLoginRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: "1"}
	fakeClock := clockwork.NewFakeClock()
	handler := loginHandler(staticGenerator, staticGenerator, fakeClock)

	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusNotFound
	if result.StatusCode != expectedStatusCode {
		t.Errorf("request failed: expected %d, got %d", expectedStatusCode, result.StatusCode)
	}
	var resp apierror.APIError
	err = json.NewDecoder(result.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("invalid status code in the response body: expected %d, got %d", expectedStatusCode, resp.StatusCode)
	}
	errMsg := resp.Errors[0]["msg"]
	expected := "account does not exist"
	if errMsg != expected {
		t.Errorf("invalid error message: expected %s, got %s", expected, errMsg)
	}
}

func TestLoginHandler_NameRequired(t *testing.T) {
	body := LoginRequest{}
	req, err := createLoginRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: "1"}
	fakeClock := clockwork.NewFakeClock()
	handler := loginHandler(staticGenerator, staticGenerator, fakeClock)

	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusBadRequest
	if result.StatusCode != expectedStatusCode {
		t.Errorf("request failed: expected %d, got %d", expectedStatusCode, result.StatusCode)
	}
	var resp apierror.APIError
	err = json.NewDecoder(result.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("invalid status code in the response body: expected %d, got %d", expectedStatusCode, resp.StatusCode)
	}
	errMsg := resp.Errors[0]["msg"]
	expected := "name is required"
	if errMsg != expected {
		t.Errorf("invalid error message: expected %s, got %s", expected, errMsg)
	}
}

func createLoginRequest(req LoginRequest) (*http.Request, error) {
	serialized, err := json.Marshal(req)
	if err != nil {
		return nil, nil
	}
	body := bytes.NewBuffer(serialized)
	httpReq := httptest.NewRequest(http.MethodPost, "/register", body)
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func loginHandler(tokenGen domain.TokenGenerator, sessionGen domain.SessionGenerator, clock clockwork.Clock) http.HandlerFunc {
	storage := inmem.NewStorage(tokenGen, sessionGen, clock)
	register := LoginHandler(storage)
	return NewHttpHandler(register)
}
