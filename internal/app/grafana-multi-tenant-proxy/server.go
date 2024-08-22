package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/auth"
	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/config"
)

func initLogger(logLevel string) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return nil, err
	}
	zapConfig.Level = level

	return zapConfig.Build()
}

// Serve serves requests to the proxy
func Serve(c *cli.Context) error {
	logLevel := c.String("log-level")
	if logLevel == "" {
		logLevel = "INFO"
	}

	logger, err := initLogger(logLevel)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not create logger %v", err), -1)
	}
	// Ensure that the logger is flushed before the program exits
	defer logger.Sync()

	errorLogger, err := zap.NewStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not create error logger %v", err), -1)
	}

	proxyConfigLocation := c.String("proxy-config")
	authConfigLocation := c.String("auth-config")
	config, err := parseConfig(proxyConfigLocation, authConfigLocation)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not parse config %v", err), -1)
	}

	var reverseProxy *httputil.ReverseProxy
	{
		targetServerURL, err := url.Parse(config.Proxy.TargetServerURL)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Could not parse target server url %v", err), -1)
		}
		reverseProxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(targetServerURL)
				r.Out.Host = targetServerURL.Host
				r.Out.Header.Set("X-Forwarded-Host", targetServerURL.Host)
				orgID := r.In.Context().Value(auth.OrgIDKey)

				if orgID != "" {
					r.Out.Header.Set("X-Scope-OrgID", orgID.(string))
				}
			},
			ErrorLog: errorLogger,
		}
	}

	authenticationMiddleware := auth.NewAuthenticationMiddleware(
		config,
		logger,
		ReverseTarget(reverseProxy),
	)

	handlers := Logger(
		authenticationMiddleware.Authenticate(),
		logger,
	)

	http.HandleFunc("/", handlers)

	// Reload config endpoint
	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}

		config, err = parseConfig(proxyConfigLocation, authConfigLocation)
		if err != nil {
			logger.Error("Could not reload config", zap.Error(err))
			w.WriteHeader(500)
		} else {
			authenticationMiddleware.ApplyConfig(config)
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		}
	})

	// Start the server
	addr := fmt.Sprintf(":%d", c.Int("port"))
	server := &http.Server{Addr: addr, ErrorLog: errorLogger}
	if err := server.ListenAndServe(); err != nil {
		return cli.Exit(fmt.Sprintf("Grafana multi tenant proxy could not start %v", err), -1)
	}
	logger.Info("Starting HTTP server", zap.String("addr", addr))
	return nil
}

func parseConfig(proxyConfigLocation string, authConfigLocation string) (config.Config, error) {
	proxyConfig, err := config.ParseProxyConfig(proxyConfigLocation)
	if err != nil {
		return config.Config{}, err
	}
	authConfig, err := config.ParseAuthConfig(authConfigLocation)
	if err != nil {
		return config.Config{}, err
	}
	return config.Config{
		Proxy:          proxyConfig,
		Authentication: authConfig,
	}, nil
}
