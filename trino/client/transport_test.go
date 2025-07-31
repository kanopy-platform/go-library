package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrinoTransport(t *testing.T) {
	var outHeader http.Header
	handler := func(w http.ResponseWriter, r *http.Request) {
		outHeader = r.Header
		w.WriteHeader(http.StatusOK)
	}
	testServer := httptest.NewServer(http.HandlerFunc(handler))

	testCases := map[string]struct {
		user          string
		bearerTokenFn func() (string, error)
		baseTransport http.RoundTripper
	}{
		"defaults": {},
		"with user": {
			user: "testuser",
		},
		"with bearer token": {
			bearerTokenFn: func() (string, error) {
				return "testtoken", nil
			},
		},
		"with user and bearer token": {
			user: "testuser",
			bearerTokenFn: func() (string, error) {
				return "testtoken", nil
			},
		},
	}

	for name, tc := range testCases {

		opts := []transportOption{}
		if tc.user != "" {
			opts = append(opts, WithUsername(tc.user))
		}
		if tc.bearerTokenFn != nil {
			opts = append(opts, WithBearerTokenFunc(tc.bearerTokenFn))
		}

		tt := NewTrinoTransport(opts...)
		assert.NotNil(t, tt, name)

		httpClient := &http.Client{
			Transport: tt,
		}

		req, err := http.NewRequest("GET", testServer.URL, nil)
		assert.NoError(t, err, name)

		_, err = httpClient.Do(req)
		assert.NoError(t, err, name)

		if tc.user != "" {
			assert.Equal(t, tc.user, outHeader.Get("X-Trino-User"), name)
		} else {
			assert.Empty(t, outHeader.Get("X-Trino-User"), name)
		}

		if tc.bearerTokenFn != nil {
			token, err := tc.bearerTokenFn()
			assert.NoError(t, err, name)
			assert.Equal(t, "Bearer "+token, outHeader.Get("Authorization"), name)
		} else {
			assert.Empty(t, outHeader.Get("Authorization"), name)
		}
	}
}
