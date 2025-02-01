package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/apierror"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/handlers"
	"github.com/Kshitij09/online-indicator/domain"
	service2 "github.com/Kshitij09/online-indicator/domain/service"
	"net/http"
)

var PathAccountId = "id"

type StatusResponse struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	IsOnline   bool   `json:"is_online"`
	LastOnline *int64 `json:"last_online,omitempty"`
}

func StatusHandler(storage domain.Storage, config domain.Config) handlers.Handler {
	service := service2.NewStatusService(
		storage.Status(),
		storage.Session(),
		config.OnlineThreshold,
		storage.Profile(),
	)
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.PathValue(PathAccountId)
		if accountId == "" {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("path parameter '%s' missing", PathAccountId))
		}
		profileStatus, err := service.Status(accountId)
		if errors.Is(err, domain.ErrAccountNotFound) {
			return apierror.SimpleAPIError(http.StatusNotFound, fmt.Sprintf("account with id '%s' not found", accountId))
		}
		resp := toTransport(profileStatus)
		return json.NewEncoder(w).Encode(resp)
	}
}

func toTransport(status domain.ProfileStatus) StatusResponse {
	var lastOnlineEpochMillis *int64
	if !status.LastOnline.IsZero() {
		*lastOnlineEpochMillis = status.LastOnline.UnixMilli()
	}
	return StatusResponse{
		Id:         status.UserId,
		Username:   status.Username,
		IsOnline:   status.IsOnline,
		LastOnline: lastOnlineEpochMillis,
	}
}
