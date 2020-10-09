package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/ResultadosDigitais/x9/config"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config
var oidcVerifier *oidc.IDTokenVerifier
var ctx context.Context

func InitOIDC() error {
	ctx = context.Background()
	provider, err := oidc.NewProvider(ctx, config.AppOpts.OIDC.ServiceURL)
	if err != nil {
		return err
	}
	conf := &oidc.Config{
		ClientID: config.AppOpts.OIDC.ClientID,
	}
	oidcVerifier = provider.Verifier(conf)
	oauth2Config = oauth2.Config{
		ClientID:     config.AppOpts.OIDC.ClientID,
		ClientSecret: config.AppOpts.OIDC.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.AppOpts.OIDC.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return nil
}

func GetAuthCodeURL(state string, extraParams map[string]string) (string, error) {
	return oauth2Config.AuthCodeURL(state, getExtraAuthOptions(extraParams)...), nil
}

func GetRawIDToken(code string) (string, error) {
	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("Missing id_token")
	}
	return rawIDToken, nil
}

func VerifyToken(rawIDToken string) (*oidc.IDToken, error) {
	idToken, err := oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	return idToken, nil
}

func getExtraAuthOptions(extraParams map[string]string) []oauth2.AuthCodeOption {
	aco := []oauth2.AuthCodeOption{}
	for k, v := range extraParams {
		aco = append(aco, oauth2.SetAuthURLParam(k, v))
	}
	return aco
}

func GetState() string {
	b := make([]byte, 16)
	rand.Read(b)
	encoded := base64.StdEncoding.EncodeToString(b)

	return encoded

}
