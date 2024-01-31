package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
)

func TestOAuthAuthentication_IsAuthorized(t *testing.T) {
	authConfig := &pkg.Authn{
		Users: []pkg.User{
			{
				Username: "user1",
				Password: "pass1",
				OrgID:    "org1",
			},
			{
				Username: "user2",
				Password: "pass2",
				OrgID:    "org2",
			},
		},
		KeepOrgID: false,
	}

	logger := zap.NewNop()

	tests := []struct {
		name        string
		token       string
		expected    bool
		expectedOrg string
	}{
		{
			name:        "Valid token",
			token:       "valid_token",
			expected:    true,
			expectedOrg: "org1",
		},
		{
			name:        "Invalid token",
			token:       "invalid_token",
			expected:    false,
			expectedOrg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := OAuthAuthentication{
				token: tt.token,
			}

			result, orgID := auth.IsAuthorized(&http.Request{}, authConfig, logger)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedOrg, orgID)
		})
	}
}

func TestOAuthAuthentication_decodeOAuthToken(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg5MjlhMzdkM2Y2OGM0Njg1OTJjOGIyODhhYjBhMTk0OGQ3MmQ4YzUifQ.eyJpc3MiOiJodHRwczovL2RleC5nb2xlbS5nYXdzLmdpZ2FudGljLmlvIiwic3ViIjoiQ2lRMU5UTTNaakk1WkMwek5UWTNMVFExTW1FdE9UQmxNUzAzTnpNeU5EUTVZalUzWldFU0RXZHBZVzUwYzNkaGNtMHRZV1EiLCJhdWQiOiIycVJnTmI1cVFCazNRcVd4STFXTGdCNXpUUTFNNGVCKyIsImV4cCI6MTcwNjcxNzAxNywiaWF0IjoxNzA2NzE1MjE3LCJhdF9oYXNoIjoiYndsX0tYSUJtbHJUbm9IRXN2RUNLdyIsImVtYWlsIjoibWFyaWVAZ2lhbnRzd2FybS5pbyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsiZ2lhbnRzd2FybS1hZDpnaWFudHN3YXJtLWFkbWlucyIsImdpYW50c3dhcm0tYWQ6R1MgU3VwcG9ydCAtIE1TIHRlYW1zIiwiZ2lhbnRzd2FybS1hZDpHaWFudCBTd2FybSBHbG9iYWwiLCJnaWFudHN3YXJtLWFkOkdpYW50U3dhcm0iLCJnaWFudHN3YXJtLWFkOkdpYW50IFN3YXJtIEVVIiwiZ2lhbnRzd2FybS1hZDpEZXZlbG9wZXJzIl0sIm5hbWUiOiJNYXJpZSBSb3F1ZSJ9.UyfIohHXBVocgv2nb-lgwVU09LJDwzHOHDb20HVZPTPMVBTLPWPzCgryg2KCXxAO1eyspdbcEQA-ZnQoqW_S6QajVyMCQyqLAECRa5h90dIvENvgj3jdcjDhCZl8q5k7Jl0WUMsBFMFMoaa3GKslM0tNcb5s-g1m0ylZocKu46qbJpiF7xWVg4ak_eWoyjb7lvBmCOSWavNHvl0Wc0Rq8HlwZHQl9Bmr5w1gZYKBcdYeMTL9_I0vnTF3UkQsvpQRsVUG9j9z86rCx3T8LsQcY_4jpOnvRVvFRbTWAWcbErvDdfOdte1TVWgBVttKq-WNBgS2HVVjk8jUAlU62k4MPA"
	expectedPayload := Payload{
		Iss: "https://dex.golem.gaws.gigantic.io",
		Aud: "2qRgNb5qQBk3QqWxI1WLgB5zTQ1M4eB+",
	}

	auth := OAuthAuthentication{
		token: token,
	}

	payload, err := auth.decodeOAuthToken()

	assert.NoError(t, err)
	assert.Equal(t, expectedPayload, payload)
}
