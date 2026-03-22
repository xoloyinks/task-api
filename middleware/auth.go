package middleware

import (
	"context"
	"net/http"
	"strings"
	"task-tracker-api/utils"
)

func AuthMiddleware(next utils.AppHandler) utils.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// 1. get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return &utils.AppError{Code: http.StatusUnauthorized, Message: "missing authorization header"}
		}

		// 2. check format is "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return &utils.AppError{Code: http.StatusUnauthorized, Message: "invalid authorization format"}
		}

		// 3. verify the token
		claims, err := utils.VerifyJWT(parts[1])
		if err != nil {
			return &utils.AppError{Code: http.StatusUnauthorized, Message: "invalid or expired token"}
		}

		// 4. attach claims to request context
		ctx := context.WithValue(r.Context(), utils.ClaimsKey, claims)

		// 5. call the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}
