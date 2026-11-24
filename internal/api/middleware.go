package api

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	apiErrors "github.com/polygonid/sh-id-platform/internal/errors"
	"github.com/polygonid/sh-id-platform/internal/log"
)

// LogMiddleware returns a middleware that adds general log configuration to each context request
func LogMiddleware(ctx context.Context) StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, operationID string) StrictHandlerFunc {
		return func(ctxReq context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (interface{}, error) {
			if reqID := middleware.GetReqID(ctxReq); reqID != "" {
				log.With("req-id", reqID)
			}
			return f(ctx, w, r, args)
		}
	}
}

// BasicAuthMiddleware returns a middleware that performs an http basic authorization for endpoints configured with
// basic auth in the api spec.
// In uses the BasicAuthScopes value in context to figure if and endpoint needs authorization or not, because this
// value is injected automatically by openapi when basic auth is selected.
// Basic auth is disabled only when both user and pass are empty. If just one of them is empty the configuration is
// incomplete, so requests to protected endpoints are rejected instead of being let through.
func BasicAuthMiddleware(ctx context.Context, user, pass string) StrictMiddlewareFunc {
	authDisabled := user == "" && pass == ""
	return func(f StrictHandlerFunc, operationID string) StrictHandlerFunc {
		return func(ctxReq context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (interface{}, error) {
			if ctxReq.Value(BasicAuthScopes) != nil && !authDisabled {
				if user == "" || pass == "" {
					return nil, apiErrors.AuthError{Err: errors.New("unauthorized")}
				}
				userReq, passReq, ok := r.BasicAuth()
				if !ok {
					return nil, apiErrors.AuthError{Err: errors.New("unauthorized")}
				}
				if subtle.ConstantTimeCompare([]byte(user), []byte(userReq)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(passReq)) != 1 {
					return nil, apiErrors.AuthError{Err: errors.New("unauthorized")}
				}
			}
			return f(ctx, w, r, args)
		}
	}
}
