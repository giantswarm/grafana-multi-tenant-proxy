package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

type key int

// Struct to represent the interesting part of the OAuth token payload section
type Payload struct {
	Iss string `json:"iss"`
	Aud string `json:"aud"`
}

const (
	// OrgIDKey Key used to pass loki tenant id though the middleware context
	OrgIDKey key = iota
)

// INTERFACE to handle different type of authentication
type Authentication interface {
	GetMode() string
	IsAuthorized(r *http.Request, authConfig *pkg.Authn, logger *zap.Logger) (bool, string)
	WriteUnauthorisedResponse(w http.ResponseWriter, logger *zap.Logger)
}

// ////////////////////////////////////////////////////////////////////////////////////
// Authentication can be used as a middleware chain to authenticate every request before proxying the request
func Authenticate(handler http.HandlerFunc, authConfig *pkg.Authn, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			for _, value := range values {
				logger.Info(fmt.Sprintf("Header %s = %s", name, value))
			}
		}
		authent := getAuthentication(r, logger)
		logger.Info(fmt.Sprintf("Authentication mode: %s", authent.GetMode()))
		ok, orgID := authent.IsAuthorized(r, authConfig, logger)
		if !ok {
			authent.WriteUnauthorisedResponse(w, logger)
			return
		}
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		handler(w, r.WithContext(ctx))
	}
}

// getAuthentication returns the authentication mode used by the request and its credentials
func getAuthentication(r *http.Request, logger *zap.Logger) Authentication {
	// OAuth token is favorite authentication mode
	token := r.Header.Get("X-Id-Token")
	if token != "" {
		logger.Info(fmt.Sprintf("Token = %s", token))
		return OAuthAuthentication{
			mode:  "oauth",
			token: token,
		}
	}
	// If no oauth token, we are looking for basicAuth
	user, pwd, ok := r.BasicAuth()
	if ok {
		return BasicAuthentication{
			mode: "basic",
			user: user,
			pwd:  pwd,
		}
	}
	return UnknowAuthentication{
		mode: "unknown",
	}
}
