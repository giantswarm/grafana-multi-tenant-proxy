package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

type key int

const (
	// OrgIDKey Key used to pass loki tenant id though the middleware context
	OrgIDKey key = iota
)

// INTERFACE to handle different type of authentication
type Authenticator interface {
	Authenticate(r *http.Request) (bool, string)
	OnAuthenticationError(w http.ResponseWriter)
}

// ////////////////////////////////////////////////////////////////////////////////////
// Authenticate can be used as a middleware chain to authenticate every request before proxying the request
func Authenticate(handler http.HandlerFunc, authConfig *pkg.Authn, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			for _, value := range values {
				logger.Info(fmt.Sprintf("Header %s = %s", name, value))
			}
		}
		for _, cookie := range r.Cookies() {
			logger.Info(fmt.Sprintf("Cookie %s", cookie))
		}

		authent, err := newAuthenticator(r, authConfig, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while authenticating request %s", r.URL), zap.Error(err))
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised\n"))
			return
		}
		logger.Debug(fmt.Sprintf("Authentication mode: %T", authent))
		ok, orgID := authent.Authenticate(r)
		if !ok {
			authent.OnAuthenticationError(w)
			return
		}
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		handler(w, r.WithContext(ctx))
	}
}

// newAuthenticator returns the authentication mode used by the request and its credentials
func newAuthenticator(r *http.Request, authConfig *pkg.Authn, logger *zap.Logger) (Authenticator, error) {
	// OAuth token is favorite authentication mode
	token := r.Header.Get("X-Id-Token")
	if token != "" {
		logger.Info(fmt.Sprintf("Token = %s", token))
		return OAuthAuthenticator{
			token:      token,
			authConfig: authConfig,
			logger:     logger,
		}, nil
	}
	// If no oauth token, we are looking for basicAuth
	user, pwd, ok := r.BasicAuth()
	if ok {
		return BasicAuthenticator{
			user:       user,
			pwd:        pwd,
			authConfig: authConfig,
			logger:     logger,
		}, nil
	}
	return nil, errors.New("Unsupported authentication")
}
