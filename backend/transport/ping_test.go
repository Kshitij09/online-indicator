package transport

import (
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	test2 "github.com/Kshitij09/online-indicator/test"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler_Unauthorized(t *testing.T) {
	clock := clockwork.NewFakeClock()
	staticGen := stubs.StaticGenerator{}
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	handler := NewHttpHandler(PingHandler(storage, test2.Config))

	recorder := httptest.NewRecorder()

	body := PingRequest{
		SessionId: "123",
	}
	req, err := test2.CreateRequest(http.MethodPost, "/ping", body)
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
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	account := domain.Account{
		Name: "John Doe",
	}
	account, err := storage.Auth().Create(account)
	if err != nil {
		t.Error(err)
	}
	loginService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())
	session, err := loginService.Login(account.Name, account.Token)
	if err != nil {
		t.Error(err)
	}

	handler := NewHttpHandler(PingHandler(storage, test2.Config))
	recorder := httptest.NewRecorder()

	body := PingRequest{SessionId: session.Id}
	req, err := test2.CreateRequest(http.MethodPost, "/ping", body)
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusOK
	if result.StatusCode != expectedStatusCode {
		t.Errorf("expected statusCode %d, got %d", expectedStatusCode, result.StatusCode)
	}
}
