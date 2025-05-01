package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/Kshitij09/online-indicator/transport/handlers"
	"net/http"
)

type RegisterRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func RegisterHandler(storage domain.Storage) handlers.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Body == http.NoBody {
			return apierror.SimpleAPIError(http.StatusBadRequest, "Request Body is missing")
		}
		var req RegisterRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&req)
		if decodeErr != nil {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", decodeErr))
		}

		acc := domain.Account{Name: req.Name}
		created, err := storage.Auth().Create(acc)
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			return apierror.SimpleAPIError(http.StatusConflict, "account already exists")
		}
		if errors.Is(err, domain.ErrEmptyName) {
			return apierror.SimpleAPIError(http.StatusBadRequest, "name is required")
		}
		response := RegisterResponse{
			Token: created.Token,
		}
		err = json.NewEncoder(w).Encode(response)
		return err
	}
}
