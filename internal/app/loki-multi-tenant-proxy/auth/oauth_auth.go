package auth

import (
	"context"
	"crypto/subtle"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"go.uber.org/zap"
)

const (
	readUser = "read"
)

// Struct to represent the interesting part of the OAuth token payload section
type Payload struct {
	Iss string `json:"iss"`
	Aud string `json:"aud"`
}

type OAuthAuthenticator struct {
	token      string
	authConfig *pkg.Authn
	logger     *zap.Logger
}

func (a OAuthAuthenticator) Authenticate(r *http.Request) (bool, string) {
	// Decode OAuth token payload section
	payload, err := a.extractPayload()
	if err != nil {
		a.logger.Error(fmt.Sprintf("Error decoding token payload %s", a.token), zap.Error(err))
		return false, ""
	}
	// Token validation against identity provider
	err = a.validate(payload, r.Context())
	if err != nil {
		a.logger.Error(fmt.Sprintf("Error while validating OAuth token against identity provider %s", a.token), zap.Error(err))
		return false, ""
	}
	// Retrieve OrgId for user 'read'
	for _, v := range a.authConfig.Users {
		// Retrieve user 'read' and get OrgID
		if subtle.ConstantTimeCompare([]byte(readUser), []byte(v.Username)) == 1 {
			if !a.authConfig.KeepOrgID {
				return true, v.OrgID
			} else {
				return true, ""
			}
		}
	}
	return false, ""
}

func (a OAuthAuthenticator) OnAuthenticationError(w http.ResponseWriter) {
	a.logger.Error("OAuth authentication failed")
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}

// extractPayload decodes the payload section of the OAuth token
func (a OAuthAuthenticator) extractPayload() (Payload, error) {
	// Get payload section from the token
	sections := strings.Split(a.token, ".")
	if len(sections) <= 1 {
		return Payload{}, errors.New("Invalid token")
	}
	payload := sections[1]
	payloadDecoded, err := b64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return Payload{}, err
	}

	var p Payload
	err = json.Unmarshal(payloadDecoded, &p)
	return p, err
}

// validate validates the OAuth token against Dex
func (a OAuthAuthenticator) validate(payload Payload, ctx context.Context) error {
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
