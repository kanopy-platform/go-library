package client

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocktrino "github.com/kanopy-platform/go-library/trino/testing"
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

// mockRetryDB creates a mock sql.DB with sqlmock for testing retry logic
// the DB will expect 3 retries and return errors for the first two queries
//func mockRetryDB(ctx context.Context, t *testing.T) (*sql.DB, error) {
//	ctrl := gomock.NewController(t)
//
//	mockStmt := mocktrino.NewMockStmtQueryContext(ctrl)
//	mockStmt.EXPECT().QueryContext(ctx, "SELECT 1").Return(nil, fmt.Errorf("mock error")).Times(2)
//	mockStmt.EXPECT().QueryContext(gomock.Any(), gomock.Any()).Return(&mocktrino.StubRows{}, nil).Times(1)
//
//	mockPrepare := mocktrino.NewMockConnPrepareContext(ctrl)
//	mockPrepare.EXPECT().PrepareContext(gomock.Any(), gomock.Any()).Return(mockStmt, nil).AnyTimes()
//
//	//	mockConn := mocktrino.NewMockConn(ctrl)
//	//	mockConn.EXPECT().PrepareConext(gomock.Any(), gomock.Any()).Return(mockStmt, nil).AnyTimes()
//	mockDriver := mocktrino.NewMockDriver(ctrl)
//	mockDriver.EXPECT().Open("mock://").Return(mockPrepare, nil).AnyTimes()
//
//	db, err := sql.Open("mock_trino", "mock://")
//
//	return db, err
//}

func TestClientRetry(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		retryCount int
		err        bool
	}{
		//		"retry 0 times": {
		//			retryCount: 0,
		//			err:        true,
		//		},
		//		"retry 1 time": {
		//			retryCount: 1,
		//			err:        true,
		//		},
		"retry 2 times": {
			retryCount: 2,
			err:        false,
		},
	}

	for name, tc := range testcases {
		//		ctx := context.TODO()
		//		db, err := mockRetryDB(ctx, t)

		ctx := context.Background()
		stmt := mocktrino.MockStmt{
			Rows: []*mocktrino.StubRows{
				nil,
				nil,
				{},
			},
			Err: []error{
				fmt.Errorf("mock error"),
				fmt.Errorf("mock error"),
				nil,
			},
		}
		conn := mocktrino.MockConn{
			Stmt: &stmt,
		}
		driver := mocktrino.MockDriver{
			Conn: &conn,
		}

		sql.Register("mock_trino", &driver)

		db, err := sql.Open("mock_trino", "mock://")
		require.NoError(t, err, name)

		client := &Client{
			conn:       db,
			retryCount: tc.retryCount,
		}

		rows, err := client.Query(ctx, "SELECT 1")
		if tc.err {
			assert.Error(t, err, name)
			assert.Nil(t, rows, name)
			continue
		}
		require.NoError(t, err, name)
		assert.NotNil(t, rows, name)
	}
}
