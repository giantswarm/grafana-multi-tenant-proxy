package auth

import (
	"crypto/subtle"
	"net/http"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

const (
	realm = "Loki multi-tenant proxy"
)

type BasicAuthenticator struct {
	user       string
	pwd        string
	authConfig *pkg.Authn
	logger     *zap.Logger
}

func (a BasicAuthenticator) Authenticate(r *http.Request) (bool, string) {
	for _, v := range a.authConfig.Users {
		// Check user and password passed in the request and get OrgID
		if subtle.ConstantTimeCompare([]byte(a.user), []byte(v.Username)) == 1 && subtle.ConstantTimeCompare([]byte(a.pwd), []byte(v.Password)) == 1 {
			if !a.authConfig.KeepOrgID {
				return true, v.OrgID
			} else {
				return true, ""
			}
		}
	}
	return false, ""
}

func (a BasicAuthenticator) OnAuthenticationError(w http.ResponseWriter) {
	a.logger.Error("Basic authentication failed")
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}
