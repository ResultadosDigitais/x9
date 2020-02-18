package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Name    = "x9"
	Version = "0.1"
)

// Config has the main analyzer configuration
type Config struct {
	GitHubAccessTokens           []Account
	SlackWebhook                 string
	BlacklistedExtensions        []string          `yaml:"blacklisted_extensions"`
	BlacklistedPaths             []string          `yaml:"blacklisted_paths"`
	BlacklistedEntropyExtensions []string          `yaml:"blacklisted_entropy_extensions"`
	Signatures                   []ConfigSignature `yaml:"signatures"`
}

//Account stores the credentials from a git account
type Account struct {
	AccountName string
	Token       string
}

//ConfigSignature stores the match signature configuration
type ConfigSignature struct {
	Name     string `yaml:"name"`
	Part     string `yaml:"part"`
	Match    string `yaml:"match,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
	Verifier string `yaml:"verifier,omitempty"`
}

//ParseConfig reads the configuration from the config.yaml file and
func ParseConfig() (*Config, error) {
	config := &Config{}

	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(path.Join(dir, "config.yaml"))
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, err
	}
	if err := config.ParseAuthConfig(); err != nil {
		return config, err
	}

	config.SlackWebhook = os.Getenv("SLACK_WEBHOOK")

	return config, nil
}

//ParseAuthConfig get the github accounts from environment
func (c *Config) ParseAuthConfig() error {
	env, ok := os.LookupEnv("GITHUB_ACCOUNTS")
	if !ok {
		return errors.New("No accounts found")
	}

	accountList := strings.Split(env, ",")

	accounts := []Account{}
	for _, account := range accountList {
		credentials := strings.Split(account, ":")
		if len(credentials) != 2 {
			return errors.New("Invalid format")
		}
		accounts = append(accounts, Account{credentials[0], credentials[1]})
	}

	c.GitHubAccessTokens = accounts
	return nil
}

//UnmarshalYAML gets the fields from  a yaml file
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = Config{}
	type plain Config
	err := unmarshal((*plain)(c))

	if err != nil {
		return err
	}

	return nil
}
