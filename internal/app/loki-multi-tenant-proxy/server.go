package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/giantswarm/loki-multi-tenant-proxy/internal/pkg"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// Serve serves
func Serve(c *cli.Context) error {
	lokiServerURL, _ := url.Parse(c.String("loki-server"))
	addr := fmt.Sprintf(":%d", c.Int("port"))
	authConfigLocation := c.String("auth-config")
	authConfig, _ := pkg.ParseConfig(&authConfigLocation)
	authConfig.KeepOrgID = c.Bool("keep-orgid")

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

	var reverseProxy *httputil.ReverseProxy
	{
		reverseProxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(lokiServerURL)
				r.Out.Host = lokiServerURL.Host
				r.Out.Header.Set("X-Forwarded-Host", lokiServerURL.Host)
				orgID := r.In.Context().Value(OrgIDKey)

				if orgID != "" {
					logger.Info("url", zap.String("url", r.In.URL.String()))
					tenantIDsInUrl := extractTenantIDsInURL(r.In.URL)
					if len(tenantIDsInUrl) > 0 {
						logger.Info("Tenant ID found in URL", zap.String("tenant_ids", strings.Join(tenantIDsInUrl, ",")))
						r.Out.Header.Set("X-Scope-OrgID", tenantIDsInUrl[0])
					} else {
						logger.Info("Tenant ID from header", zap.String("tenant_ids", orgID.(string)))
						r.Out.Header.Set("X-Scope-OrgID", orgID.(string))
					}
				}
			},
			ErrorLog: errorLogger,
		}
	}

	handlers := Logger(
		BasicAuth(
			ReverseLoki(reverseProxy),
			authConfig,
		),
		logger,
	)

	http.HandleFunc("/", handlers)
	server := &http.Server{Addr: addr, ErrorLog: errorLogger}
	if err := server.ListenAndServe(); err != nil {
		return cli.Exit(fmt.Sprintf("Loki multi tenant proxy could not start %v", err), -1)
	}
	logger.Info("Starting HTTP server", zap.String("addr", addr))
	return nil
}

func extractTenantIDsInURL(url *url.URL) []string {
	tenantIDs := []string{}
	if strings.HasPrefix(url.Path, "/loki/api/v1/query_range") || strings.HasPrefix(url.Path, "/loki/api/v1/index/stats") {
		query := url.Query().Get("query")
		if strings.Contains(query, "x5f4k") {
			tenantIDs = append(tenantIDs, "x5f4k")
		}
	}
	return tenantIDs
}
