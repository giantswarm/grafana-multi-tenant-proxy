package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/giantswarm/loki-multi-tenant-proxy/internal/app/loki-multi-tenant-proxy/auth"
	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var keepOrgID bool
var authConfigLocation string
var authConfig *pkg.Authn

func loadConfig() (*pkg.Authn, error) {
	config, err := pkg.ParseConfig(&authConfigLocation)
	config.KeepOrgID = keepOrgID
	return config, err
}

// Serve serves
func Serve(c *cli.Context) error {
	lokiServerURL, _ := url.Parse(c.String("loki-server"))
	addr := fmt.Sprintf(":%d", c.Int("port"))
	authConfigLocation = c.String("auth-config")
	keepOrgID = c.Bool("keep-orgid")
	logLevel := c.String("log-level")
	if logLevel == "" {
		logLevel = "INFO"
	}

	var logger *zap.Logger
	{
		zapConfig := zap.NewProductionConfig()
		level, err := zap.ParseAtomicLevel(logLevel)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Could not parse log level %v", err), -1)
		}
		zapConfig.Level = level

		logger = zap.Must(zapConfig.Build())
		defer logger.Sync()
	}

	errorLogger, err := zap.NewStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not create standard logger %v", err), -1)
	}

	authConfig, err = loadConfig()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not parse config %v", err), -1)
	}

	var reverseProxy *httputil.ReverseProxy
	{
		reverseProxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(lokiServerURL)
				r.Out.Host = lokiServerURL.Host
				r.Out.Header.Set("X-Forwarded-Host", lokiServerURL.Host)
				orgID := r.In.Context().Value(auth.OrgIDKey)

				if orgID != "" {
					r.Out.Header.Set("X-Scope-OrgID", orgID.(string))
				}
			},
			ErrorLog: errorLogger,
		}
	}

	authenticationMiddleware := auth.NewAuthenticationMiddleware(
		logger,
		ReverseLoki(reverseProxy),
		*authConfig,
	)

	handlers := Logger(
		authenticationMiddleware.Authenticate(),
		logger,
	)

	http.HandleFunc("/", handlers)
	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}
		authConfig, err := loadConfig()

		if err != nil {
			logger.Error("Could not reload config", zap.Error(err))
			w.WriteHeader(500)
		} else {
			authenticationMiddleware.ApplyConfig(*authConfig)
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		}
	})

	server := &http.Server{Addr: addr, ErrorLog: errorLogger}
	if err := server.ListenAndServe(); err != nil {
		return cli.Exit(fmt.Sprintf("Loki multi tenant proxy could not start %v", err), -1)
	}
	logger.Info("Starting HTTP server", zap.String("addr", addr))
	return nil
}
