package config

type AppConfig struct {
	ApplicationSecretKey string
	DatabaseConfig       DatabaseConfig
	AppURL               string
	OIDC                 OIDC
}

var AppOpts AppConfig

func ParseAppConfig() error {
	var err error

	if AppOpts.DatabaseConfig, err = ParseDatabaseConfig(); err != nil {
		return err
	}
	if AppOpts.ApplicationSecretKey, err = getEnv("APP_SECRET_KEY"); err != nil {
		return err
	}
	if AppOpts.AppURL, err = getEnv("APP_URL"); err != nil {
		return err
	}
	if AppOpts.OIDC, err = ParseOIDCConfig(); err != nil {
		return err
	}

	return nil
}
