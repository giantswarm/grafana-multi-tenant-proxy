package auth

import (
	"net/http"
	
	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

// /////////////////// Unknown Authentication //////////////////////////////////////////
type UnknowAuthentication struct {
	mode string
}

func (a UnknowAuthentication) GetMode() string {
	return a.mode
}

func (a UnknowAuthentication) IsAuthorized(r *http.Request, authConfig *pkg.Authn, logger *zap.Logger) (bool, string) {
	return false, ""
}

func (a UnknowAuthentication) WriteUnauthorisedResponse(w http.ResponseWriter, logger *zap.Logger) {
	logger.Error("Unknown authentication mode")
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}
