package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	Name    = "x9"
	Version = "0.1"
)

var (
	Opts = Options{}
)

type Options struct {
	Threads                  *int     `yaml:"num_threads"`
	MaximumRepositorySize    *int     `yaml:"maximum_repository_size"`
	MaximumFileSize          *int     `yaml:"maximum_file_size"`
	CloneRepositoryTimeout   *int     `yaml:"clone_repository_timeout"`
	BlackListedRepositories  []string `yaml:"blacklisted_repositories"`
	SlackActionsUsersAllowed []string
	DatabaseConfig           DatabaseConfig
	SlackWebhook             string
	SlackSecretKey           string
	GithubSecretWebhook      string
	GithubAccountUser        string
	GithubAccountToken       string
}

func ParseConfig() error {
	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(path.Join(dir, "config.yaml"))
	if err != nil {
		return err
	}
	Opts.SlackActionsUsersAllowed = getSlackUsersAllowed()
	if Opts.DatabaseConfig, err = ParseDatabaseConfig(); err != nil {
		return err
	}
	Opts.SlackWebhook, _ = getEnv("SLACK_WEBHOOK")
	Opts.SlackSecretKey, _ = getEnv("SLACK_SECRET_KEY")

	if Opts.GithubSecretWebhook, err = getEnv("GITHUB_SECRET_WEBHOOK"); err != nil {
		return err
	}
	if Opts.GithubAccountUser, err = getEnv("GITHUB_ACCOUNT_USER"); err != nil {
		return err
	}
	if Opts.GithubAccountToken, err = getEnv("GITHUB_ACCOUNT_TOKEN"); err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &Opts)
	return err
}

func getEnv(envName string) (string, error) {
	env, ok := os.LookupEnv(envName)
	if !ok {
		return env, errors.New(fmt.Sprintf("Missing environment variable: %s", envName))
	}
	return env, nil
}
func getSlackUsersAllowed() []string {
	if env, ok := os.LookupEnv("SLACK_USERS_ALLOWED"); ok {
		return strings.Split(env, ",")
	}
	return nil
}

func GetBlacklistedRepositories() []string {
	return Opts.BlackListedRepositories
}
