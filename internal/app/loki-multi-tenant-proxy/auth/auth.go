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
type AuthenticationMiddleware struct {
	handler    http.HandlerFunc
	authConfig *pkg.Authn
	logger     *zap.Logger
}

func NewAuthenticationMiddleware(logger *zap.Logger, handler http.HandlerFunc, authConfig pkg.Authn) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		handler:    handler,
		authConfig: &authConfig,
		logger:     logger,
	}
}

// ////////////////////////////////////////////////////////////////////////////////////
// Authenticate can be used as a middleware chain to authenticate every request before proxying the request
func (am AuthenticationMiddleware) Authenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticator, err := newAuthenticator(r, am.authConfig, am.logger)
		if err != nil {
			am.logger.Error(fmt.Sprintf("Error while authenticating request %s", r.URL), zap.Error(err))
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised\n"))
			return
		}
		am.logger.Debug(fmt.Sprintf("Authentication mode: %T", authenticator))
		ok, orgID := authenticator.Authenticate(r)
		if !ok {
			authenticator.OnAuthenticationError(w)
			return
		}
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		am.handler(w, r.WithContext(ctx))
	}
}

func (am AuthenticationMiddleware) ApplyConfig(authConfig pkg.Authn) {
	*am.authConfig = authConfig
}

// newAuthenticator returns the authentication mode used by the request and its credentials
func newAuthenticator(r *http.Request, authConfig *pkg.Authn, logger *zap.Logger) (Authenticator, error) {
	// OAuth token is favorite authentication mode
	token := r.Header.Get("X-Id-Token")
	if token != "" {
		logger.Debug(fmt.Sprintf("OAuth Token = %s", token))
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
	return nil, errors.New("unsupported authentication")
}
