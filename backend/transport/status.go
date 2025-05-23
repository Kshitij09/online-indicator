package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/Kshitij09/online-indicator/domain/service"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/Kshitij09/online-indicator/transport/handlers"
	"net/http"
	"sort"
	"strings"
)

type StatusResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsOnline   bool   `json:"is_online"`
	LastOnline *int64 `json:"last_online,omitempty"`
}

type BatchStatusRequest struct {
	Ids []string `json:"ids"`
}
type BatchStatusResponse struct {
	Items []StatusResponse `json:"items"`
}

func StatusHandler(svc service.StatusService) handlers.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountId := r.PathValue(PathId)
		if accountId == "" {
			return apierror.SimpleAPIError(http.StatusBadRequest, fmt.Sprintf("path parameter '%s' missing", PathId))
		}
		profileStatus, err := svc.Status(accountId)
		if errors.Is(err, domain.ErrAccountNotFound) {
			return apierror.SimpleAPIError(http.StatusNotFound, fmt.Sprintf("account with id '%s' not found", accountId))
		}
		resp := toTransport(profileStatus)
		return json.NewEncoder(w).Encode(resp)
	}
}

func BatchStatusHandler(svc service.StatusService) handlers.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		request := BatchStatusRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return apierror.SimpleAPIError(http.StatusBadRequest, "invalid request body")
		}
		response := BatchStatusResponse{}
		statuses := svc.BatchStatus(request.Ids)
		items := toTransportItems(statuses)
		response.Items = items
		return json.NewEncoder(w).Encode(response)
	}
}

func toTransport(status domain.ProfileStatus) StatusResponse {
	var lastOnlineEpochMillis *int64
	if !status.LastOnline.IsZero() {
		var onlineMillis = status.LastOnline.UnixMilli()
		lastOnlineEpochMillis = &onlineMillis
	}
	return StatusResponse{
		Id:         status.UserId,
		Name:       status.Name,
		IsOnline:   status.IsOnline,
		LastOnline: lastOnlineEpochMillis,
	}
}

func toTransportItems(statuses map[string]domain.ProfileStatus) []StatusResponse {
	items := make([]StatusResponse, 0, len(statuses))
	for _, status := range statuses {
		items = append(items, toTransport(status))
	}
	sort.Slice(items, func(i, j int) bool {
		return strings.Compare(items[i].Name, items[j].Name) < 0
	})
	return items
}
