package auth

import (
	"crypto/subtle"
	"net/http"

	"go.uber.org/zap"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/config"
)

const (
	realm = "Grafana multi-tenant proxy"
)

type BasicAuthenticator struct {
	user   string
	pwd    string
	config *config.Config
	logger *zap.Logger
}

<<<<<<< HEAD:internal/app/grafana-multi-tenant-proxy/handler/auth/basic_auth.go
func (a BasicAuthenticator) Authenticate(r *http.Request, targetServer *config.TargetServer) (bool, string) {
	for _, v := range a.config.Authentication.Users {
		// Check user and password passed in the request and get OrgID
		if subtle.ConstantTimeCompare([]byte(a.user), []byte(v.Username)) == 1 && subtle.ConstantTimeCompare([]byte(a.pwd), []byte(v.Password)) == 1 {
			if !targetServer.KeepOrgID {
=======
func (a BasicAuthenticator) Authenticate(r *http.Request) (bool, string) {
	for _, v := range a.config.Authentication.Users {
		// Check user and password passed in the request and get OrgID
		if subtle.ConstantTimeCompare([]byte(a.user), []byte(v.Username)) == 1 && subtle.ConstantTimeCompare([]byte(a.pwd), []byte(v.Password)) == 1 {
			if !a.config.Proxy.KeepOrgID {
>>>>>>> 2eb33b2 (Improve config management):internal/app/grafana-multi-tenant-proxy/auth/basic_auth.go
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
