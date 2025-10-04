package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/MujiRahman/golang-simple-note/internal/helper"
	"github.com/MujiRahman/golang-simple-note/internal/service"
	"github.com/MujiRahman/golang-simple-note/pkg/contextkey"
)

// AuthMiddleware returns a middleware that validates JWT and injects userID into request context.
// It accepts a UserService to parse token (DI).
func AuthMiddleware(userSvc service.UserService) func(h httprouter.Handle) httprouter.Handle {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			token, err := extractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				helper.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
				return
			}
			uid, err := userSvc.ParseToken(token)
			if err != nil {
				helper.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}
			// inject userID into context
			ctx := context.WithValue(r.Context(), contextkey.UserIDKey, uid)
			// pass the httprouter.Params along, by creating new request with ctx and using httprouter.Handle
			h(w, r.WithContext(ctx), ps)
		}
	}
}

func extractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("no header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid header")
	}
	return parts[1], nil
}
