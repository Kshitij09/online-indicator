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

type PingRequest struct {
	SessionId string `json:"sessionId"`
}

func PingHandler(storage domain.Storage, config domain.Config) handlers.Handler {
	service := service2.NewStatusService(storage.Status(), storage.Session(), config.OnlineThreshold)
	return func(w http.ResponseWriter, r *http.Request) error {
		var req PingRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&req)
		if decodeErr != nil {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", decodeErr))
		}
		err := service.Ping(req.SessionId)
		if errors.Is(err, domain.ErrSessionNotFound) {
			return apierror.SimpleAPIError(http.StatusUnauthorized, fmt.Sprintf("session not found: %s", req.SessionId))
		}
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		return nil
	}
}
