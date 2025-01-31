package transport

import (
	"github.com/Kshitij09/online-indicator/cmd/http-server/test"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler_Unauthorized(t *testing.T) {
	clock := clockwork.NewFakeClock()
	staticGen := stubs.StaticGenerator{}
	storage := inmem.NewStorage(staticGen, staticGen, clock)
	handler := NewHttpHandler(PingHandler(storage))

	recorder := httptest.NewRecorder()

	body := PingRequest{
		SessionId: "123",
	}
	req, err := test.CreateRequest(http.MethodPost, "/ping", body)
	if err != nil {
		t.Error(err)
	}

	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusUnauthorized
	if result.StatusCode != expectedStatusCode {
		t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
	}
}

func TestPingHandler_OK(t *testing.T) {
	staticGen := stubs.StaticGenerator{
		StubValue: "123",
	}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock)
	account := domain.Account{
		Name: "John Doe",
	}
	account, err := storage.Auth().Create(account)
	if err != nil {
		t.Error(err)
	}
	loginService := domain.NewLoginService(storage.Auth(), storage.Session())
	session, err := loginService.Login(account.Name, account.Token)
	if err != nil {
		t.Error(err)
	}

	handler := NewHttpHandler(PingHandler(storage))
	recorder := httptest.NewRecorder()

	body := PingRequest{SessionId: session.Id}
	req, err := test.CreateRequest(http.MethodPost, "/ping", body)
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusOK
	if result.StatusCode != expectedStatusCode {
		t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
	}
}
