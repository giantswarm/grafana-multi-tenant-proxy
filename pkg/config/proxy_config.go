package config

import (
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

type ProxyConfig struct {
	TargetServers []TargetServer `yaml:"targetServers"`
}

type TargetServer struct {
	Name      string `yaml:"name"`
	Host      string `yaml:"host"`
	Target    string `yaml:"target"`
	KeepOrgID bool   `yaml:"keepOrgId"`
}

// readProxyConfigFile read a configuration file in the path `location` and returns a ProxyConfig object
func readProxyConfigFile(location string) (*ProxyConfig, error) {
	data, err := os.ReadFile(filepath.Clean(location))
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

func (p ProxyConfig) FindTargetServer(host string) *TargetServer {
	for _, v := range p.TargetServers {
		if v.Host == host {
			return &v
		}
	}
	return nil
}
