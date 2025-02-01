package transport

import (
	"encoding/json"
	"github.com/Kshitij09/online-indicator/cmd/http-server/test"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStatusHandler_Success(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	authService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())
	acc := domain.Account{Name: "john"}
	acc, err := authService.CreateAccount(acc)
	if err != nil {
		t.Error(err)
	}
	session, err := authService.Login(acc.Name, acc.Token)
	if err != nil {
		t.Error(err)
	}
	statusService := service.NewStatusService(storage.Status(), storage.Session(), test.Config.OnlineThreshold, storage.Profile())
	err = statusService.Ping(session.Id)
	if err != nil {
		t.Error(err)
	}

	handler := NewHttpHandler(StatusHandler(storage, test.Config))

	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", acc.Id)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusOK
	if result.StatusCode != expectedStatusCode {
		t.Errorf("status code should be %d, got %d", expectedStatusCode, result.StatusCode)
	}
	var body StatusResponse
	err = json.NewDecoder(result.Body).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	expectedOnlineMillis := clock.Now().UnixMilli()
	expectedBody := StatusResponse{
		Id:         acc.Id,
		Username:   acc.Name,
		IsOnline:   true,
		LastOnline: &expectedOnlineMillis,
	}
	if !reflect.DeepEqual(body, expectedBody) {
		expectedBodyJson, _ := json.Marshal(expectedBody)
		bodyJson, _ := json.Marshal(body)
		t.Errorf("response body does not match\nexpected:\n%s\n\nactual:\n%s", string(expectedBodyJson), string(bodyJson))
	}
}

func TestStatusHandler_AccountNotFound(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	handler := NewHttpHandler(StatusHandler(storage, test.Config))

	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", "john")
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusNotFound
	if result.StatusCode != expectedStatusCode {
		t.Errorf("status code should be %d, got %d", expectedStatusCode, result.StatusCode)
	}
}
