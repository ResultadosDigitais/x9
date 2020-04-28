package sast

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Config has the main analyzer configuration
type Config struct {
	BlacklistedExtensions        []string          `yaml:"blacklisted_extensions"`
	BlacklistedPaths             []string          `yaml:"blacklisted_paths"`
	BlacklistedEntropyExtensions []string          `yaml:"blacklisted_entropy_extensions"`
	Signatures                   []ConfigSignature `yaml:"signatures"`
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
func ParseConfig() (Config, error) {
	config := Config{}

	dir, _ := os.Getwd()
	data, err := ioutil.ReadFile(path.Join(dir, "sast/config.yaml"))
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
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
