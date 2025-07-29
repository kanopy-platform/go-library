package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testURI(username, password, server string) string {
	return fmt.Sprintf("https://%s:%s@%s", username, password, server)
}

func defaultTestURI() string {
	return testURI("my-account@example.com", "password", "trino.example.com:443?catalog=test")
}

func TestNew(t *testing.T) {
	t.Parallel()

	uri := defaultTestURI()

	c, err := New(uri)
	assert.NoError(t, err)
	assert.Equal(t, uri, c.config.ServerURI)
}

func TestWithCustomClient(t *testing.T) {
	t.Parallel()

	customClientName := "custom"

	c, err := New(defaultTestURI(), WithCustomClient(customClientName, http.DefaultClient))
	require.NoError(t, err)
	assert.Equal(t, c.config.CustomClientName, customClientName)
}

func TestConnectionConfig(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		username  string
		password  string
		server    string
		encrypted bool
		expected  string
		err       bool
	}{
		"valid connection with username and password": {
			username:  "foo",
			password:  "bar",
			server:    "trino.example.com:443?catalog=test",
			encrypted: true,
			expected:  "https://foo:bar@trino.example.com:443?catalog=test",
		},
		"valid no user encrypted": {
			username:  "",
			password:  "",
			server:    "trino.example.com:443?catalog=test",
			encrypted: true,
			expected:  "https://trino.example.com:443?catalog=test",
		},
		"invalid no servername": {
			username:  "foo",
			password:  "bar",
			server:    "",
			encrypted: true,
			err:       true,
		},
		"invalid no user with password": {
			username:  "",
			password:  "bar",
			server:    "trino.example.com:443?catalog=test",
			encrypted: true,
			err:       true,
		},
		"invalid no password with user": {
			username:  "foo",
			password:  "",
			server:    "trino.example.com:443?catalog=test",
			encrypted: true,
			err:       true,
		},
		"invalid user and unencrypted": {
			username:  "foo",
			password:  "bar",
			server:    "trino.example.com:443?catalog=test",
			encrypted: false,
			err:       true,
		},
	}

	for name, tc := range testcases {
		conConf := NewConnectionConfig(WithUser(tc.username), WithPassword(tc.password), WithServerHost(tc.server), WithEncrypted(tc.encrypted))

		serverURI, err := conConf.Parse()
		if tc.err {
			assert.Error(t, err, name)
			assert.Empty(t, serverURI, name)
			continue
		}
		require.NoError(t, err, name)
		assert.Equal(t, tc.expected, serverURI, name)
	}
}
