package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Proxy          *ProxyConfig
	Authentication *AuthenticationConfig
}

type ProxyConfig struct {
	TargetServerURL string `yaml:"targetServerURL"`
	KeepOrgID       bool   `yaml:"keepOrgId"`
}

// AuthenticationConfig contains a list of users
type AuthenticationConfig struct {
	Users []User `yaml:"users"`
}

// User Identifies a user including the tenant
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	OrgID    string `yaml:"orgid"`
}

// ParseConfig read a configuration file in the path `locatino` and returns a ProxyConfig object
func ParseProxyConfig(location string) (*ProxyConfig, error) {
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

// ParseAuthConfig read a configuration file in the path `location` and returns a Config object
func ParseAuthConfig(location string) (*AuthenticationConfig, error) {
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
