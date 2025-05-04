package transport

import (
	"encoding/json"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/Kshitij09/online-indicator/testfixtures"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler_Unauthorized(t *testing.T) {
	clock := clockwork.NewFakeClock()
	staticGen := stubs.StaticGenerator{StubValue: "expected_token"}
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	handler := NewHttpHandler(PingHandler(storage, testfixtures.Config, clock))
	t.Run("session token not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		body := PingRequest{
			SessionToken: "non_existent_token",
		}
		req, err := testfixtures.CreateRequest(http.MethodPost, "/ping", body)
		if err != nil {
			t.Error(err)
		}
		req.SetPathValue(PathId, "random") // invalid account id
		handler(recorder, req)
		result := recorder.Result()
		expectedStatusCode := http.StatusUnauthorized
		if result.StatusCode != expectedStatusCode {
			var errBody apierror.APIError
			json.NewDecoder(result.Body).Decode(&errBody)
			firstMsg := errBody.Errors[0]["msg"]
			t.Logf("error: %s", firstMsg)
			t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
		}
	})

	t.Run("invalid session token", func(t *testing.T) {
		// Create account and session
		account := domain.Account{
			Name: "John Doe",
		}
		account, err := storage.Auth().Create(account)
		if err != nil {
			t.Error(err)
		}

		handler := NewHttpHandler(PingHandler(storage, testfixtures.Config, clock))
		recorder := httptest.NewRecorder()

		body := PingRequest{
			SessionToken: "invalid_token",
		}
		req, err := testfixtures.CreateRequest(http.MethodPost, "/ping", body)
		if err != nil {
			t.Error(err)
		}
		req.SetPathValue(PathId, account.Id) // valid account id

		handler(recorder, req)

		result := recorder.Result()
		expectedStatusCode := http.StatusUnauthorized
		if result.StatusCode != expectedStatusCode {
			var errBody apierror.APIError
			json.NewDecoder(result.Body).Decode(&errBody)
			firstMsg := errBody.Errors[0]["msg"]
			t.Logf("error: %s", firstMsg)
			t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
		}
	})
}

func TestPingHandler_BadRequest(t *testing.T) {
	clock := clockwork.NewFakeClock()
	staticGen := stubs.StaticGenerator{}
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	handler := NewHttpHandler(PingHandler(storage, testfixtures.Config, clock))

	recorder := httptest.NewRecorder()

	body := PingRequest{
		SessionToken: "123",
	}
	req, err := testfixtures.CreateRequest(http.MethodPost, "/ping", body)
	if err != nil {
		t.Error(err)
	}
	handler(recorder, req)
	if err != nil {
		t.Error(err)
	}

	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusBadRequest
	if result.StatusCode != expectedStatusCode {
		var errBody apierror.APIError
		json.NewDecoder(result.Body).Decode(&errBody)
		firstMsg := errBody.Errors[0]["msg"]
		t.Logf("error: %s", firstMsg)
		t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
	}
}

func TestPingHandler_OK(t *testing.T) {
	staticGen := stubs.StaticGenerator{
		StubValue: "123",
	}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	account := domain.Account{
		Name: "John Doe",
	}
	account, err := storage.Auth().Create(account)
	if err != nil {
		t.Error(err)
	}
	loginService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())
	session, err := loginService.Login(staticGen.StubValue, account.ApiKey)
	if err != nil {
		t.Error(err)
	}

	handler := NewHttpHandler(PingHandler(storage, testfixtures.Config, clock))
	recorder := httptest.NewRecorder()

	body := PingRequest{SessionToken: session.Token}
	req, err := testfixtures.CreateRequest(http.MethodPost, "/ping", body)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue(PathId, session.AccountId)
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusOK
	if result.StatusCode != expectedStatusCode {
		var errBody apierror.APIError
		json.NewDecoder(result.Body).Decode(&errBody)
		firstMsg := errBody.Errors[0]["msg"]
		t.Logf("error: %s", firstMsg)
		t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
	}
}
