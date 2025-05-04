package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	service2 "github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/Kshitij09/online-indicator/transport/handlers"
	"github.com/jonboulle/clockwork"
	"net/http"
)

type PingRequest struct {
	SessionToken string `json:"sessionToken"`
}

func PingHandler(storage domain.Storage, config domain.Config, clock clockwork.Clock) handlers.Handler {
	service := service2.NewStatusService(
		storage.Session(),
		config.OnlineThreshold,
		storage.Profile(),
		clock,
	)
	return func(w http.ResponseWriter, r *http.Request) error {
		var req PingRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&req)
		if decodeErr != nil {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", decodeErr))
		}
		err := service.Ping(req.SessionToken)
		if errors.Is(err, domain.ErrSessionNotFound) {
			return apierror.SimpleAPIError(http.StatusUnauthorized, fmt.Sprintf("session not found: %s", req.SessionToken))
		}
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		return nil
	}
}
