package okta

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type Client struct {
	*okta.APIClient
}

func NewClientFromJWKBytes(orgURL string, clientID string, jwkBytes []byte, scopes ...string) (*Client, error) {
	jwk, err := jwkFromBytes(jwkBytes)
	if err != nil {
		return nil, err
	}
	pem, err := jwkToRSA(jwk)
	if err != nil {
		return nil, err
	}
	return NewClient(orgURL, clientID, pem, scopes...)
}

func NewClient(orgURL string, clientID string, key string, scopes ...string) (*Client, error) {
	config, err := okta.NewConfiguration(
		okta.WithOrgUrl(orgURL),
		okta.WithClientId(clientID),
		okta.WithScopes(scopes),
		okta.WithAuthorizationMode("PrivateKey"),
		okta.WithPrivateKey(key),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build okta configuration: %w", err)
	}

	client := okta.NewAPIClient(config)
	return &Client{client}, nil
}

type ListGroupUsersOpt func(okta.ApiListGroupUsersRequest) okta.ApiListGroupUsersRequest

func WithLimit(limit int32) ListGroupUsersOpt {
	return func(r okta.ApiListGroupUsersRequest) okta.ApiListGroupUsersRequest {
		return r.Limit(limit)
	}
}

func (c *Client) ListGroupUsers(ctx context.Context, groupId string, opts ...ListGroupUsersOpt) ([]okta.GroupMember, error) {
	query := c.GroupAPI.ListGroupUsers(ctx, groupId)
	for _, opt := range opts {
		query = opt(query)
	}

	users, resp, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to query okta group users: %w", err)
	}

	for resp.HasNextPage() {
		var nextSet []okta.GroupMember
		resp, err = resp.Next(&nextSet)
		if err != nil {
			return nil, fmt.Errorf("failed to receive pagination results: %w", err)
		}
		users = append(users, nextSet...)
	}

	return users, nil
}

func (c *Client) GroupByName(ctx context.Context, groupName string) (*okta.Group, error) {
	query := c.GroupAPI.ListGroups(ctx).Q(groupName)
	oktaGroups, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to query okta group: %w", err)
	}

	for _, group := range oktaGroups {
		if group.Profile != nil && group.Profile.Name != nil && *group.Profile.Name == groupName {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("unable to find okta group %q", groupName)
}

func jwkFromBytes(bytes []byte) (*jose.JSONWebKey, error) {
	jwk := &jose.JSONWebKey{}
	if err := json.Unmarshal(bytes, jwk); err != nil {
		return nil, fmt.Errorf("failed to marhsal jwk bytes to json: %w", err)
	}
	return jwk, nil
}

func jwkToRSA(jwk *jose.JSONWebKey) (string, error) {
	rsaPrivateKey, ok := jwk.Key.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("key data must be of type *rsa.PrivateKey")
	}

	pemKey := &strings.Builder{}
	if err := pem.Encode(pemKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey)}); err != nil {
		return "", fmt.Errorf("failed to marshal pkcs1 private key: %w", err)
	}

	return pemKey.String(), nil
}
