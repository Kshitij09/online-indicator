package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/apierror"
	"github.com/Kshitij09/online-indicator/cmd/http-server/transport/handlers"
	"github.com/Kshitij09/online-indicator/domain"
	"net/http"
)

type LoginRequest struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(storage domain.Storage) handlers.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Body == http.NoBody {
			return apierror.SimpleAPIError(http.StatusBadRequest, "Request Body is missing")
		}
		var req LoginRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&req)
		if decodeErr != nil {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", decodeErr))
		}

		if req.Name == "" {
			return apierror.SimpleAPIError(http.StatusBadRequest, "name is required")
		}

		acc, err := storage.Auth().Login(req.Name, req.Token)
		if errors.Is(err, domain.ErrAccountNotFound) {
			return apierror.SimpleAPIError(http.StatusNotFound, "account does not exist")
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return apierror.SimpleAPIError(http.StatusUnauthorized, "invalid credentials")
		}
		response := LoginResponse{
			Token: acc.Token,
		}
		err = json.NewEncoder(w).Encode(response)
		return err
	}
}
