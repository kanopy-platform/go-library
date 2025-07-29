package client

import (
	"fmt"
	"net/http"
)

type bearerTokenFunc func() (string, error)

func defaultBearerTokenFunc() (string, error) {
	// Implement your logic to retrieve the bearer token
	return "", nil
}

type transportOption func(*TrinoTransport)

// WithUsername sets the Trino user to be used in the X-Trino-User header.
func WithUsername(user string) transportOption {
	return func(tt *TrinoTransport) {
		tt.user = user
	}
}

// WithBaseTransport sets the base HTTP transport to be used by the TrinoTransport.
// The RoundTrip method of the base transport will be called to perform the actual HTTP request.
// the default is http.DefaultTransport. This can be used to support local testing configurations
// such as overriding the TLS verification settings.
func WithBaseTransport(baseTransport http.RoundTripper) transportOption {
	return func(tt *TrinoTransport) {
		tt.base = baseTransport
	}
}

// WithBearerTokenFunc sets a custom function to retrieve the bearer token for the Authorization header.
// it shoudl return a base64 encoded beaer token string.
// an empty string means no Authorization header will be set.
func WithBearerTokenFunc(fn func() (string, error)) transportOption {
	return func(tt *TrinoTransport) {
		tt.bearerTokenFn = fn
	}
}

// NewTrinoTransport creates a new TrinoTransport with the provided options.
// It initializes the transport with default values and applies the provided options.
// The default base transport is http.DefaultTransport, and the default bearer token function disabling
// the Authorization.
func NewTrinoTransport(opts ...transportOption) *TrinoTransport {
	tt := &TrinoTransport{
		base:          http.DefaultTransport,
		bearerTokenFn: defaultBearerTokenFunc,
	}

	for _, opt := range opts {
		opt(tt)
	}

	return tt
}

// TrinoTransport is a custom HTTP transport for Trino that adds the necessary headers
// such as X-Trino-User and Authorization (Bearer token) to the requests.
type TrinoTransport struct {
	// base is the underlying HTTP transport to use for making requests
	// defaults to http.DefaultTransport
	base http.RoundTripper
	// bearerTokenFn is a function that retrieves the bearer token to be used in the Authorization header
	// if it returns an empty string, no Authorization header will be set
	bearerTokenFn bearerTokenFunc
	// user is the Trino user to set in the X-Trino-User header
	user string
}

// RoundTrip injects the X-Trino-User header and the Authorization header with the bearer token
// when necessary, and then calls the base RoundTrip method to perform the actual HTTP request.
// If the bearerTokenFn returns an error, it will return an error instead of making the request.
func (tt *TrinoTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if tt.user != "" {
		req.Header.Set("X-Trino-User", tt.user)
	}

	token, err := tt.bearerTokenFn()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bearer token: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return tt.base.RoundTrip(req)
}
