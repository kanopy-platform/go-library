package client

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/trinodb/trino-go-client/trino"
)

var (
	boolTrue  bool = true
	boolFalse bool = false
)

// ConnectionConfig constructs a base ServerURI for a Trino connection.
// if enforces HTTPS encryption by default, but can be configured to use HTTP.
// user and passerword are optional but required encyption is enabled.
// ConnectionConfigs are used to provide a connection string to New to create a new Client.
// It is a convenience struct to make conditional connection string configuration easier.
/*

Example usage:

conectionConfig := NewConnectionConfig(WithServerHost("trino.example.com:443"),WithUser("read-only-user"), WithPassword("password"))

serverURI, err := connectionConfig.Parse()
if err != nil {
	log.Fatal(err)
}

client, err := New(serverURI)
if err != nil {
	log.Fatal(err)
}

cursor, err := client.Query(context.Background(), "SELECT * FROM my_table WHERE id = ?", 1)
*/

type ConnectionConfig struct {
	// Encrypted indicates whether the connection should use HTTPS, optional, default: true
	Encrypted *bool `json:"encrypted,omitempty"`
	// Password is the password for authentication, optional
	Password string `json:"password,omitempty"`
	// ServerHost the Trino server hostname including port, e.g. "localhost:8080", required
	ServerHost string `json:"serverHost"`
	/// User is the username for authentication, optional
	User string `json:"user,omitempty"`
}

type connectionOption func(*ConnectionConfig)

// NewConnectionConfig creates a new ConnectionConfig with encryption enabled by default.
func NewConnectionConfig(opts ...connectionOption) *ConnectionConfig {
	out := &ConnectionConfig{
		Encrypted: &boolTrue,
	}

	for _, opt := range opts {
		opt(out)
	}

	return out
}

// WithUnencrypted enables or disables https encryption.
func WithEncrypted(encrypted bool) connectionOption {
	return func(c *ConnectionConfig) {
		c.Encrypted = &encrypted
	}
}

// WithPassword sets the password for authentication.
func WithPassword(password string) connectionOption {
	return func(c *ConnectionConfig) {
		c.Password = password
	}
}

// WithServerHost sets the server host for the Trino connection.
// The server host should include the port if not standard http(s) ports, e.g. "localhost:8080".
func WithServerHost(serverHost string) connectionOption {
	return func(c *ConnectionConfig) {
		c.ServerHost = serverHost
	}
}

// WithUser sets the user for authentication.
func WithUser(user string) connectionOption {
	return func(c *ConnectionConfig) {
		c.User = user
	}
}

// Parse validates the ConnejctionConfig values and formats the connection configuration into a ServerURI string for a trino.Config.
func (c *ConnectionConfig) Parse() (string, error) {
	if c.ServerHost == "" {
		return "", fmt.Errorf("ServerHost must be provided")
	}

	rawDSN := "http://"

	encrypted := c.Encrypted != nil && *c.Encrypted
	if encrypted {
		rawDSN = "https://"
	}

	if c.User == "" && c.Password != "" {
		return "", fmt.Errorf("user must be provided if password is set")
	}

	if c.User != "" && c.Password == "" {
		return "", fmt.Errorf("password must be provided if user is set")
	}

	if c.User != "" && c.Password != "" {
		if !encrypted {
			return "", fmt.Errorf("encryption must be abled for user authentication")
		}
		rawDSN = fmt.Sprintf("https://%s:%s@", c.User, c.Password)
	}

	rawDSN += c.ServerHost

	return rawDSN, nil
}

type Client struct {
	config     trino.Config
	dsn        string
	conn       *sql.DB
	retryCount int
}

func New(uri string, opts ...option) (*Client, error) {
	c := &Client{
		config: trino.Config{
			ServerURI: uri,
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	dsn, err := c.config.FormatDSN()
	if err != nil {
		// cannot log error directly because DSN contains username and password
		return nil, fmt.Errorf("malformed server URI, please check your connection string and try again")
	}
	c.dsn = dsn

	return c, nil
}

func (c *Client) Connect() error {
	db, err := sql.Open("trino", c.dsn)
	if err != nil {
		if db != nil {
			db.Close()
		}
		// cannot log error directly because DSN contains username and password
		return fmt.Errorf("malformed server URI, please check your connection string and try again")
	}
	c.conn = db
	log.Info("Connected to Trino")

	return nil
}

func (c *Client) Disconnect() {
	c.conn.Close()
	log.Info("Connection to Trino closed")
}

func (c *Client) Query(ctx context.Context, statement string, args ...any) (*sql.Rows, error) {
	if c.conn == nil {
		if e := c.Connect(); e != nil {
			return nil, e
		}
	}
	stmt, err := c.conn.PrepareContext(ctx, statement)
	if err != nil {
		return nil, fmt.Errorf("prepare Error: %s", err)
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

type option func(*Client) error

// WithCustomClient registers a custom HTTP client for the Trino connection.
// all coonnections issued from this client will use this custom client.
// This provides support for custom HTTP clients with the configuration like
// the TrinoTransport
/*


transport := NewTrinoTransport(WithUsername("demo-access-policy-name"))

customHTTP := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

trinoClient := New("https://trino.example.com:443", WithCustomClient("custom-client", customHTTP))



*/
func WithCustomClient(name string, customClient *http.Client) option {
	return func(c *Client) error {
		if err := trino.RegisterCustomClient(name, customClient); err != nil {
			//this is basically unreachable
			return err
		}
		c.config.CustomClientName = name

		return nil
	}
}
