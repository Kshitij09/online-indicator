package transport

import (
	"encoding/json"
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/domain/stubs"
	"github.com/Kshitij09/online-indicator/inmem"
	"github.com/Kshitij09/online-indicator/testfixtures"
	"github.com/jonboulle/clockwork"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strconv"
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
	session, err := authService.Login(acc.Id, acc.Token)
	if err != nil {
		t.Error(err)
	}
	statusService := service.NewStatusService(
		storage.Session(),
		testfixtures.Config.OnlineThreshold,
		storage.Profile(),
		clock,
	)
	err = statusService.Ping(session.Id)
	if err != nil {
		t.Error(err)
	}

	handler := NewHttpHandler(StatusHandler(storage, testfixtures.Config, clock))

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
	handler := NewHttpHandler(StatusHandler(storage, testfixtures.Config, clock))

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

func TestStatusHandler_NoLoginAsOffline(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	authService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())
	acc := domain.Account{Name: "john"}
	acc, err := authService.CreateAccount(acc)
	if err != nil {
		t.Error(err)
	}
	handler := NewHttpHandler(StatusHandler(storage, testfixtures.Config, clock))

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
	if body.IsOnline {
		t.Errorf("account should not be online")
	}
	if body.Username == "" || body.Id == "" {
		t.Errorf("account details should not be empty")
	}
}

func TestStatusHandler_MissingAccountId(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	handler := NewHttpHandler(StatusHandler(storage, testfixtures.Config, clock))

	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusBadRequest
	if result.StatusCode != expectedStatusCode {
		t.Errorf("status code should be %d, got %d", expectedStatusCode, result.StatusCode)
	}
}

func TestBatchStatusHandler_Success(t *testing.T) {
	seqGen := domain.NewSeqIdGenerator()
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, seqGen)
	authService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())
	accIds := make([]string, 0, 100)
	expected := make([]StatusResponse, 0, 100)
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("test%d", i)
		acc := domain.Account{Name: name}
		acc, err := authService.CreateAccount(acc)
		if err != nil {
			t.Error(err)
		}
		accIds = append(accIds, acc.Id)
		session, err := authService.Login(acc.Id, acc.Token)
		if err != nil {
			t.Error(err)
		}
		statusService := service.NewStatusService(
			storage.Session(),
			testfixtures.Config.OnlineThreshold,
			storage.Profile(),
			clock,
		)
		err = statusService.Ping(session.Id)
		if err != nil {
			t.Error(err)
		}
		lastOnlineMillis := clock.Now().UnixMilli()
		response := StatusResponse{
			Id:         acc.Id,
			Username:   acc.Name,
			IsOnline:   true,
			LastOnline: &lastOnlineMillis,
		}
		expected = append(expected, response)
	}

	handler := NewHttpHandler(BatchStatusHandler(storage, testfixtures.Config, clock))

	reqBody := BatchStatusRequest{Ids: accIds}
	req, err := testfixtures.CreateRequest(http.MethodPost, "/batch-status", reqBody)
	if err != nil {
		t.Error(err)
	}
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	result := recorder.Result()
	expectedStatusCode := http.StatusOK
	if result.StatusCode != expectedStatusCode {
		t.Errorf("status code should be %d, got %d", expectedStatusCode, result.StatusCode)
	}

	var body BatchStatusResponse
	err = json.NewDecoder(result.Body).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	slices.SortFunc(expected, statusResponseLessById)
	actual := body.Items
	slices.SortFunc(actual, statusResponseLessById)
	if !reflect.DeepEqual(actual, expected) {
		t.Error("Response body does not match")
	}
}

func TestStatusHandler_OnlineToOfflineAfterThreshold(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	authService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())

	// Register
	acc := domain.Account{Name: "john"}
	acc, err := authService.CreateAccount(acc)
	if err != nil {
		t.Error(err)
	}

	// Login
	session, err := authService.Login(acc.Id, acc.Token)
	if err != nil {
		t.Error(err)
	}

	// Ping to set online
	statusService := service.NewStatusService(
		storage.Session(),
		testfixtures.Config.OnlineThreshold,
		storage.Profile(),
		clock,
	)
	err = statusService.Ping(session.Id)
	if err != nil {
		t.Error(err)
	}

	handler := NewHttpHandler(StatusHandler(storage, testfixtures.Config, clock))

	// Verify online status
	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", acc.Id)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	result := recorder.Result()
	var body StatusResponse
	err = json.NewDecoder(result.Body).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	if !body.IsOnline {
		t.Error("user should be online")
	}

	// Advance time beyond threshold
	clock.Advance(testfixtures.Config.OnlineThreshold + 1)

	// Verify offline status
	req, err = http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", acc.Id)
	recorder = httptest.NewRecorder()
	handler(recorder, req)

	result = recorder.Result()
	body = StatusResponse{}
	err = json.NewDecoder(result.Body).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	if body.IsOnline {
		t.Error("user should be offline")
	}
}

func TestStatusHandler_OfflineToOnline(t *testing.T) {
	staticGen := stubs.StaticGenerator{StubValue: "123"}
	clock := clockwork.NewFakeClock()
	storage := inmem.NewStorage(staticGen, staticGen, clock, staticGen)
	authService := service.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())

	// Register
	acc := domain.Account{Name: "john"}
	acc, err := authService.CreateAccount(acc)
	if err != nil {
		t.Error(err)
	}

	// Login
	session, err := authService.Login(acc.Id, acc.Token)
	if err != nil {
		t.Error(err)
	}

	// Advance time beyond initial ping by logic
	clock.Advance(testfixtures.Config.OnlineThreshold + 1)

	handler := NewHttpHandler(StatusHandler(storage, testfixtures.Config, clock))

	// Verify offline status
	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", acc.Id)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	result := recorder.Result()
	var body StatusResponse
	err = json.NewDecoder(result.Body).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	if body.IsOnline {
		t.Error("user should be offline")
	}

	// Ping to set online
	statusService := service.NewStatusService(
		storage.Session(),
		testfixtures.Config.OnlineThreshold,
		storage.Profile(),
		clock,
	)
	err = statusService.Ping(session.Id)
	if err != nil {
		t.Error(err)
	}

	// Verify online status
	req, err = http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("id", acc.Id)
	recorder = httptest.NewRecorder()
	handler(recorder, req)

	result = recorder.Result()
	body = StatusResponse{}
	err = json.NewDecoder(result.Body).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	if !body.IsOnline {
		t.Error("user should be online")
	}
}

func statusResponseLessById(first StatusResponse, second StatusResponse) int {
	aId, _ := strconv.Atoi(first.Id)
	bId, _ := strconv.Atoi(second.Id)
	if aId < bId {
		return -1
	} else if aId > bId {
		return 1
	}
	return 0
}
