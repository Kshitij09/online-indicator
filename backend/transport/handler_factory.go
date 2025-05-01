package transport

import (
	"errors"
	"github.com/Kshitij09/online-indicator/transport/apierror"
	"github.com/Kshitij09/online-indicator/transport/handlers"
	"github.com/Kshitij09/online-indicator/transport/middlewares"
	"github.com/Kshitij09/online-indicator/transport/writer"
	"log"
	"net/http"
)

func NewHttpHandler(handler handlers.Handler, middlewares ...middlewares.Middleware) http.HandlerFunc {
	for _, mdl := range middlewares {
		handler = mdl(handler)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			var apiErr *apierror.APIError
			var writeErr error
			if errors.As(err, &apiErr) {
				writeErr = writer.ErrorJson(w, apiErr)
			} else {
				log.Println(err)
				writeErr = writer.ErrorJson(w, apierror.InternalServerError)
			}
			if writeErr != nil {
				log.Fatal(writeErr)
			}
		}
	}
}
