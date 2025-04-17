package config

import (
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

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

// readAuthConfigFile read a configuration file in the path `location` and returns a Config object
func readAuthConfigFile(location string) (*AuthenticationConfig, error) {
	data, err := os.ReadFile(filepath.Clean(location))
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
