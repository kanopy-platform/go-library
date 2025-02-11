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
	oktaClient *okta.APIClient
}

func NewClient(orgURL string, clientID string, jwkBytes []byte, scopes ...string) (*Client, error) {
	jwk := &jose.JSONWebKey{}
	if err := json.Unmarshal(jwkBytes, jwk); err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := jwk.Key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key data must be of type *rsa.PrivateKey")
	}

	pemKey := &strings.Builder{}
	if err := pem.Encode(pemKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey)}); err != nil {
		return nil, err
	}

	config, err := okta.NewConfiguration(
		okta.WithOrgUrl(orgURL),
		okta.WithAuthorizationMode("PrivateKey"),
		okta.WithClientId(clientID),
		okta.WithScopes(scopes),
		okta.WithPrivateKey(pemKey.String()),
	)
	if err != nil {
		return nil, err
	}

	client := okta.NewAPIClient(config)
	return &Client{oktaClient: client}, nil
}

type ListGroupUsersOpt func(okta.ApiListGroupUsersRequest) okta.ApiListGroupUsersRequest

func WithLimit(limit int32) ListGroupUsersOpt {
	return func(r okta.ApiListGroupUsersRequest) okta.ApiListGroupUsersRequest {
		return r.Limit(limit)
	}
}

func (c *Client) ListGroupUsers(ctx context.Context, groupId string, opts ...ListGroupUsersOpt) (*[]okta.GroupMember, error) {
	query := c.oktaClient.GroupAPI.ListGroupUsers(ctx, groupId)
	for _, opt := range opts {
		query = opt(query)
	}

	users, resp, err := c.oktaClient.GroupAPI.ListGroupUsersExecute(query)
	if err != nil {
		return nil, err
	}

	for resp.HasNextPage() {
		var nextSet []okta.GroupMember
		resp, err = resp.Next(&nextSet)
		users = append(users, nextSet...)
	}

	return &users, nil
}

func (c *Client) GroupByName(ctx context.Context, groupName string) (*okta.Group, error) {
	query := c.oktaClient.GroupAPI.ListGroups(ctx).Q(groupName)
	oktaGroups, _, err := query.Execute()
	if err != nil {
		return nil, err
	}

	// okta group names are unique
	if len(oktaGroups) != 1 {
		return nil, fmt.Errorf("unable to find okta group %q", groupName)
	}
	return &oktaGroups[0], nil
}
