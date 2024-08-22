package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	Proxy          ProxyConfig
	Authentication AuthenticationConfig
}

type ProxyConfig struct {
	TargetServerURL string `yaml:"targetServerURL"`
	KeepOrgID       bool   `yaml:"keepOrgId"`
}

// AuthenticationConfig contains a list of users
type AuthenticationConfig struct {
	Users []User `yaml:"users"`
}

// User defines a user credentials and its tenant
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	OrgID    string `yaml:"orgid"`
}

// ReadProxyConfigFile read a configuration file in the path `location` and returns a ProxyConfig object
func ReadProxyConfigFile(location string) (*ProxyConfig, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	config := ProxyConfig{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// ReadAuthConfigFile read a configuration file in the path `location` and returns a Config object
func ReadAuthConfigFile(location string) (*AuthenticationConfig, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	config := AuthenticationConfig{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
