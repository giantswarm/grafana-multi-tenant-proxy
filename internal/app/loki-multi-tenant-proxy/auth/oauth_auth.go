package auth

import (
	"context"
	"crypto/subtle"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

type OAuthAuthentication struct {
	mode  string
	token string
}

func (a OAuthAuthentication) GetMode() string {
	return a.mode
}

func (a OAuthAuthentication) IsAuthorized(r *http.Request, authConfig *pkg.Authn, logger *zap.Logger) (bool, string) {
	// Decode OAuth token payload section
	payload, err := a.decodeOAuthToken()
	if err != nil {
		logger.Error(fmt.Sprintf("Error decoding token payload %s", a.token), zap.Error(err))
		return false, ""
	}
	// Token validation against identity provider
	err = a.validateOAuthToken(payload, r.Context())
	if err != nil {
		logger.Error(fmt.Sprintf("Error while validating OAuth token against identity provider %s", a.token), zap.Error(err))
		return false, ""
	}
	// Retrieve OrgId for user 'read'
	for _, v := range authConfig.Users {
		// Retrieve user 'read' and get OrgID
		if subtle.ConstantTimeCompare([]byte(readUser), []byte(v.Username)) == 1 {
			if !authConfig.KeepOrgID {
				return true, v.OrgID
			} else {
				return true, ""
			}
		}
	}
	return false, ""
}

func (a OAuthAuthentication) WriteUnauthorisedResponse(w http.ResponseWriter, logger *zap.Logger) {
	logger.Error("OAuth authentication failed")
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}

// decodeOAuthToken decodes the payload section of the OAuth token
func (a OAuthAuthentication) decodeOAuthToken() (Payload, error) {
	// Get payload section from the token
	payload := strings.Split(a.token, ".")[1]
	payloadDecoded, err := b64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return Payload{}, err
	}

	var p Payload
	err = json.Unmarshal(payloadDecoded, &p)
	return p, err
}

// validateOAuthToken validates the OAuth token against Dex
func (a OAuthAuthentication) validateOAuthToken(payload Payload, ctx context.Context) error {
	// Initialize a provider by specifying dex's issuer URL.
	provider, err := oidc.NewProvider(ctx, payload.Iss)
	if err != nil {
		return err
	}
	// Create an ID token parser, but only trust ID tokens issued to 'clientId'
	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: payload.Aud})
	// Verify token validity
	_, err = idTokenVerifier.Verify(ctx, a.token)
	return err
}
