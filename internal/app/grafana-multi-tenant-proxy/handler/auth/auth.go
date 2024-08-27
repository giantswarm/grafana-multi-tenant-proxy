package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/giantswarm/grafana-multi-tenant-proxy/pkg/config"
)

type key int

const (
	// OrgIDKey Key used to pass tenant id though the middleware context
	OrgIDKey key = iota
)

// INTERFACE to handle different type of authentication
type Authenticator interface {
	Authenticate(r *http.Request, targetServer *config.TargetServer) (bool, string)
	OnAuthenticationError(w http.ResponseWriter)
}

type AuthenticationMiddleware struct {
	handler http.HandlerFunc
	config  *config.Config
	logger  *zap.Logger
}

func NewAuthenticationMiddleware(config *config.Config, logger *zap.Logger, handler http.HandlerFunc) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		handler: handler,
		config:  config,
		logger:  logger,
	}
}

// ////////////////////////////////////////////////////////////////////////////////////
// Authenticate can be used as a middleware chain to authenticate every request before proxying the request
func (am AuthenticationMiddleware) Authenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticator, err := newAuthenticator(r, am.config, am.logger)
		if err != nil {
			am.logger.Error("Error while authenticating request", zap.String("url", r.URL.String()), zap.Error(err))
			w.WriteHeader(401)
			_, err := w.Write([]byte("Unauthorised\n"))
			if err != nil {
				am.logger.Error("Could not write response", zap.Error(err))
			}
			return
		}

		targetServer := am.config.Proxy.FindTargetServer(r.Host)
		if targetServer == nil {
			am.logger.Error("Target server not configured",
				zap.String("host", r.Host),
				zap.String("url", r.URL.String()),
				zap.Error(err),
			)
			w.WriteHeader(404)
			_, err := w.Write([]byte("Not found\n"))
			if err != nil {
				am.logger.Error("Could not write response", zap.Error(err))
			}
			return
		}

		am.logger.Debug(fmt.Sprintf("Authentication mode: %T", authenticator))
		ok, orgID := authenticator.Authenticate(r, targetServer)
		if !ok {
			authenticator.OnAuthenticationError(w)
			return
		}
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		am.handler(w, r.WithContext(ctx))
	}
}

func (am AuthenticationMiddleware) ApplyConfig(config *config.Config) {
	*am.config = *config
}

// newAuthenticator returns the authentication mode used by the request and its credentials
func newAuthenticator(r *http.Request, config *config.Config, logger *zap.Logger) (Authenticator, error) {
	// OAuth token is favorite authentication mode
	token := r.Header.Get("X-Id-Token")
	if token != "" {
		return OAuthAuthenticator{
			token:  token,
			config: config,
			logger: logger,
		}, nil
	}
	// If no oauth token, we are looking for basicAuth
	user, pwd, ok := r.BasicAuth()
	if ok {
		return BasicAuthenticator{
			user:   user,
			pwd:    pwd,
			config: config,
			logger: logger,
		}, nil
	}
	return nil, errors.New("unsupported authentication")
}
