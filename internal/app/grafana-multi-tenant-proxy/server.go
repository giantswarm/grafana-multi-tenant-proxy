package proxy

import (
	"fmt"
	"net/http"

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

	// Read the configuration
	proxyConfigLocation := c.String("proxy-config")
	authConfigLocation := c.String("auth-config")
	cfg, err := config.ReadConfigFiles(proxyConfigLocation, authConfigLocation)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not parse config %v", err), -1)
	}

	proxy := NewProxy(cfg, errorLogger)
	authenticationMiddleware := auth.NewAuthenticationMiddleware(cfg, logger, proxy.Handler())
	handlers := Logger(authenticationMiddleware.Authenticate(), logger)

	http.HandleFunc("/", handlers)

	// Reload config endpoint
	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}

		cfg, err = config.ReadConfigFiles(proxyConfigLocation, authConfigLocation)
		if err != nil {
			logger.Error("Could not reload config", zap.Error(err))
			w.WriteHeader(500)
		} else {
			authenticationMiddleware.ApplyConfig(cfg)
			proxy.ApplyConfig(cfg)
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
