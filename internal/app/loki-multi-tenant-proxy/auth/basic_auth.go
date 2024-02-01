package auth

import (
	"crypto/subtle"
	"net/http"

	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

const (
	realm = "Loki multi-tenant proxy"
)

type BasicAuthentication struct {
	mode string
	user string
	pwd  string
}

func (a BasicAuthentication) GetMode() string {
	return a.mode
}

func (a BasicAuthentication) IsAuthorized(r *http.Request, authConfig *pkg.Authn, logger *zap.Logger) (bool, string) {
	for _, v := range authConfig.Users {
		// Check user and password passed in the request and get OrgID
		if subtle.ConstantTimeCompare([]byte(a.user), []byte(v.Username)) == 1 && subtle.ConstantTimeCompare([]byte(a.pwd), []byte(v.Password)) == 1 {
			if !authConfig.KeepOrgID {
				return true, v.OrgID
			} else {
				return true, ""
			}
		}
	}
	return false, ""
}

func (a BasicAuthentication) WriteUnauthorisedResponse(w http.ResponseWriter, logger *zap.Logger) {
	logger.Error("Basic authentication failed")
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}
