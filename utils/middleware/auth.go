package middleware

import (
	"api_shope/utils/helper"
	"context"
	"net/http"
	"strings"
)

type key int

const UserContextKey key = 0

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.WriteError(w, http.StatusUnauthorized, "tak ada token")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := helper.ValidateJWT(tokenString)
		if err != nil {
			helper.WriteError(w, http.StatusForbidden, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
