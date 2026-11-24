package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	apiErrors "github.com/polygonid/sh-id-platform/internal/errors"
)

func TestBasicAuthMiddleware(t *testing.T) {
	type credentials struct{ user, password string }
	for _, tc := range []struct {
		name       string
		configured credentials
		sent       *credentials // nil sends no Authorization header
		protected  bool         // endpoint declares basic auth in the api spec
		wantCalled bool
	}{
		{name: "correct credentials", configured: credentials{"user", "password"}, sent: &credentials{"user", "password"}, protected: true, wantCalled: true},
		{name: "wrong password", configured: credentials{"user", "password"}, sent: &credentials{"user", "wrong"}, protected: true},
		{name: "no credentials", configured: credentials{"user", "password"}, protected: true},
		{name: "endpoint without basic auth", configured: credentials{"user", "password"}, protected: false, wantCalled: true},
		{name: "auth disabled with empty user and password", configured: credentials{}, protected: true, wantCalled: true},
		{name: "empty configured password, no credentials", configured: credentials{"user", ""}, protected: true},
		{name: "empty configured password, empty password sent", configured: credentials{"user", ""}, sent: &credentials{"user", ""}, protected: true},
		{name: "empty configured user, no credentials", configured: credentials{"", "password"}, protected: true},
		{name: "empty configured user, empty user sent", configured: credentials{"", "password"}, sent: &credentials{"", "password"}, protected: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			called := false
			next := func(context.Context, http.ResponseWriter, *http.Request, interface{}) (interface{}, error) {
				called = true
				return nil, nil
			}
			handler := BasicAuthMiddleware(ctx, tc.configured.user, tc.configured.password)(next, "GetIdentities")

			reqCtx := ctx
			if tc.protected {
				// The generated router sets this untyped string key, and the middleware reads it.
				reqCtx = context.WithValue(ctx, BasicAuthScopes, []string{}) //nolint:staticcheck // must match the generated key
			}
			req := httptest.NewRequest(http.MethodGet, "/v2/identities", nil)
			if tc.sent != nil {
				req.SetBasicAuth(tc.sent.user, tc.sent.password)
			}

			_, err := handler(reqCtx, httptest.NewRecorder(), req, nil)

			assert.Equal(t, tc.wantCalled, called)
			if tc.wantCalled {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &apiErrors.AuthError{})
			}
		})
	}
}
