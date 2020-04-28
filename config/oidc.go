package config

type OIDC struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ServiceURL   string
}

func ParseOIDCConfig() (OIDC, error) {
	oidc := OIDC{}
	var err error
	if oidc.ClientID, err = getEnv("OIDC_CLIENT_ID"); err != nil {
		return oidc, err
	}
	if oidc.ServiceURL, err = getEnv("OIDC_SERVICE_URL"); err != nil {
		return oidc, err
	}
	if oidc.ClientSecret, err = getEnv("OIDC_CLIENT_SECRET"); err != nil {
		return oidc, err
	}
	if oidc.RedirectURL, err = getEnv("OIDC_REDIRECT_URL"); err != nil {
		return oidc, err
	}
	return oidc, nil
}
