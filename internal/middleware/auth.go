package middleware

import (
	"factorio/internal/config"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

func NewAuthMiddleware(cfgManager *config.ConfigManager) grove.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKeyHeader := r.Header.Get("authorization")
			if apiKeyHeader == "" {
				grove.WriteErrorToResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			apiKey := apiKeyHeader[len("Bearer "):]
			cfg := cfgManager.GetConfig().Auth
			if apiKey != cfg.ApiKey {
				grove.WriteErrorToResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
