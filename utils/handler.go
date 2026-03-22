package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

var validate = validator.New()

func Make(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// check if it's a known AppError
			if appErr, ok := err.(*AppError); ok {
				WriteJSON(w, appErr.Code, map[string]string{
					"error": appErr.Message,
				})
				return
			}

			// check if it's a validation error
			if _, ok := err.(validator.ValidationErrors); ok {
				WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
					"error":  "validation failed",
					"fields": FormatValidationErrors(err),
				})
				return
			}

			// unknown error — 500
			WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
