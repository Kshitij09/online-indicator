package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	service2 "github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/Kshitij09/online-indicator/transport/handlers"
	"net/http"
)

type LoginRequest struct {
	Id     string `json:"id"`
	ApiKey string `json:"apikey"`
}

type LoginResponse struct {
	SessionToken string `json:"sessionToken"`
}

func LoginHandler(storage domain.Storage) handlers.Handler {
	service := service2.NewAuthService(storage.Auth(), storage.Session(), storage.Profile())
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Body == http.NoBody {
			return apierror.SimpleAPIError(http.StatusBadRequest, "Request Body is missing")
		}
		var req LoginRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&req)
		if decodeErr != nil {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", decodeErr))
		}

		if req.Id == "" {
			return apierror.SimpleAPIError(http.StatusBadRequest, "id is required")
		}

		session, err := service.Login(req.Id, req.ApiKey)
		if errors.Is(err, domain.ErrAccountNotFound) {
			return apierror.SimpleAPIError(http.StatusNotFound, "account does not exist")
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return apierror.SimpleAPIError(http.StatusUnauthorized, "invalid credentials")
		}
		response := LoginResponse{
			SessionToken: session.Token,
		}
		err = json.NewEncoder(w).Encode(response)
		return err
	}
}
