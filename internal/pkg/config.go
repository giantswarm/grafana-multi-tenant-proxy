package pkg

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Authn Contains a list of users
type Authn struct {
	Users     []User `yaml:"users"`
	KeepOrgID bool
}

// User Identifies a user including the tenant
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	OrgID    string `yaml:"orgid"`
}

// ParseConfig read a configuration file in the path `location` and returns an Authn object
func ParseConfig(location *string) (*Authn, error) {
	data, err := os.ReadFile(*location)
	if err != nil {
		return nil, err
	}
	authn := Authn{}
	err = yaml.Unmarshal([]byte(data), &authn)
	if err != nil {
		return nil, err
	}
	return &authn, nil
}
