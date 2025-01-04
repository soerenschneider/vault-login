package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type VaultTokenSource struct {
	client     *api.Client
	authMethod api.AuthMethod
}

func NewVaultTokenSource(_ Config, auth api.AuthMethod) (*VaultTokenSource, error) {
	if auth == nil {
		return nil, errors.New("empty authmethod provided")
	}

	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	return &VaultTokenSource{
		client:     client,
		authMethod: auth,
	}, nil
}

func (s *VaultTokenSource) Receive(ctx context.Context) (string, error) {
	secret, err := s.client.Auth().Login(ctx, s.authMethod)
	if err != nil {
		return "", fmt.Errorf("failed to login to Vault: %w", err)
	}

	return secret.Auth.ClientToken, nil
}

func (s *VaultTokenSource) Cleanup(ctx context.Context) error {
	return s.client.Auth().Token().RevokeSelf("xxx")
}
