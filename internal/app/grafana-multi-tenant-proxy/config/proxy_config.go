package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

type ProxyConfig struct {
	TargetServers []TargetServer `yaml:"targetServers"`
	KeepOrgID     bool           `yaml:"keepOrgId"`
}

type TargetServer struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Target string `yaml:"target"`
}

// ReadProxyConfigFile read a configuration file in the path `location` and returns a ProxyConfig object
func readProxyConfigFile(location string) (*ProxyConfig, error) {
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
