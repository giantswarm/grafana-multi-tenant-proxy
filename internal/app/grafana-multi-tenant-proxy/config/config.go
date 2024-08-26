package config

type Config struct {
	Proxy          ProxyConfig
	Authentication AuthenticationConfig
}

func ReadConfigFiles(proxyConfigLocation string, authConfigLocation string) (*Config, error) {
	proxyConfig, err := readProxyConfigFile(proxyConfigLocation)
	if err != nil {
		return nil, err
	}
	authConfig, err := readAuthConfigFile(authConfigLocation)
	if err != nil {
		return nil, err
	}
	return &Config{
		Proxy:          *proxyConfig,
		Authentication: *authConfig,
	}, nil
}
