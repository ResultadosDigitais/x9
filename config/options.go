package config

import (
	"io/ioutil"
	"os"
	"path"

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
	Threads                 *int     `yaml:"num_threads"`
	Debug                   *bool    `yaml:"debug_mode"`
	MaximumRepositorySize   *uint    `yaml:"maximum_repository_size"`
	MaximumFileSize         *uint    `yaml:"maximum_file_size"`
	CloneRepositoryTimeout  *int     `yaml:"clone_repository_timeout"`
	BlackListedRepositories []string `yaml:"blacklisted_repositories"`
}

func ParseConfig() error {
	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(path.Join(dir, "config.yaml"))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &Opts)
	return err
}

func GetBlacklistedRepositories() []string {
	return Opts.BlackListedRepositories
}
