package transport

import (
	"bytes"
	"encoding/json"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/apierror"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler_Success(t *testing.T) {
	body := RegisterRequest{Name: "test"}
	req, err := createRegisterRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: "1"}
	fakeClock := clockwork.NewFakeClock()
	handler := registerHandler(staticGenerator, staticGenerator, fakeClock, staticGenerator)
	handler(recorder, req)

	result := recorder.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("request failed: expected %d, got %d", http.StatusOK, result.StatusCode)
	}
	var resp RegisterResponse
	err = json.NewDecoder(result.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Token != staticGenerator.StubValue {
		t.Errorf("token incorrect, expected %s, got %s", staticGenerator.StubValue, resp.Token)
	}
}

func TestRegisterHandler_AccountExists(t *testing.T) {
	body := RegisterRequest{Name: "test"}
	req, err := createRegisterRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: "1"}
	fakeClock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGenerator, staticGenerator, fakeClock, staticGenerator)
	register := RegisterHandler(storage)
	handler := NewHttpHandler(register)

	existing := domain.Account{Name: "test"}
	_, err = storage.Auth().Create(existing)
	if err != nil {
		t.Error(err)
	}
	handler(recorder, req)

	result := recorder.Result()
	if result.StatusCode != http.StatusConflict {
		t.Errorf("request failed: expected %d, got %d", http.StatusConflict, result.StatusCode)
	}
	var resp apierror.APIError
	err = json.NewDecoder(result.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("invalid status code in the response body: expected %d, got %d", http.StatusConflict, resp.StatusCode)
	}
	errMsg := resp.Errors[0]["msg"]
	expected := "account already exists"
	if errMsg != expected {
		t.Errorf("invalid error message: expected %s, got %s", expected, errMsg)
	}
}

func TestRegisterHandler_NameRequired(t *testing.T) {
	body := RegisterRequest{}
	req, err := createRegisterRequest(body)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	staticGenerator := stubs.StaticGenerator{StubValue: "1"}
	fakeClock := clockwork.NewFakeClock()
	handler := registerHandler(staticGenerator, staticGenerator, fakeClock, staticGenerator)

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

func createRegisterRequest(req RegisterRequest) (*http.Request, error) {
	serialized, err := json.Marshal(req)
	if err != nil {
		return nil, nil
	}
	body := bytes.NewBuffer(serialized)
	httpReq := httptest.NewRequest(http.MethodPost, "/register", body)
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func registerHandler(
	tokenGen domain.TokenGenerator,
	sessionGen domain.SessionGenerator,
	clock clockwork.Clock,
	idGen domain.IDGenerator,
) http.HandlerFunc {
	storage := inmem.NewStorage(tokenGen, sessionGen, clock, idGen)
	register := RegisterHandler(storage)
	return NewHttpHandler(register)
}
