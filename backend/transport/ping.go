package transport

import (
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

func PingHandler(storage domain.Storage, config domain.Config, clock clockwork.Clock, lastSeen domain.LastSeenDao) handlers.Handler {
	service := service2.NewStatusService(
		storage.Session(),
		config.OnlineThreshold,
		storage.Profile(),
		lastSeen,
		clock,
	)
	return func(w http.ResponseWriter, r *http.Request) error {
		sessionToken := r.Header.Get(HeaderSessionToken)
		if sessionToken == "" {
			return apierror.SimpleAPIError(http.StatusUnauthorized, fmt.Sprintf("header '%s' is missing", HeaderSessionToken))
		}
		id := r.PathValue(PathId)
		if id == "" {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("path parameter '%s' missing", PathId))
		}
		err := service.Ping(id, sessionToken)
		if errors.Is(err, domain.ErrSessionNotFound) {
			return apierror.SimpleAPIError(http.StatusUnauthorized, "session not found")
		}
		if errors.Is(err, domain.ErrInvalidSession) {
			return apierror.SimpleAPIError(http.StatusUnauthorized, "invalid session")
		}
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		return nil
	}
}
