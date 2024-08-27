package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/config"
	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/handler"
	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/handler/auth"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grafana_multi_tenant_proxy_http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"handler", "code", "method"})

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:                            "grafana_multi_tenant_proxy_http_request_duration_seconds",
			Help:                            "Histogram of latencies for HTTP requests.",
			Buckets:                         []float64{.05, 0.1, .25, .5, .75, 1, 2, 5, 20, 60},
		},
		[]string{"handler", "method"},
	)
	responseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grafana_multi_tenant_proxy_http_response_size_bytes",
			Help:    "Histogram of response size for HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"handler", "method"},
	)
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
		return cli.Exit(fmt.Sprintf("Could not create standard error logger %v", err), -1)
	}

	// Read the configuration
	proxyConfigLocation := c.String("proxy-config")
	authConfigLocation := c.String("auth-config")
	cfg, err := config.ReadConfigFiles(proxyConfigLocation, authConfigLocation)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Could not parse config %v", err), -1)
	}

	proxy := handler.NewProxy(cfg, logger, errorLogger)
	authenticationMiddleware := auth.NewAuthenticationMiddleware(cfg, logger, proxy.Handler())
	handlers := handler.Logger(authenticationMiddleware.Authenticate(), logger)

	// Register Prometheus collectors
	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	// We handle metrics first to avoid calling the authentication middleware
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", instrumentHandler("default", handlers))

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

func instrumentHandler(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	handlerLabel := prometheus.Labels{"handler": handlerName}
	return promhttp.InstrumentHandlerDuration(
		requestDuration.MustCurryWith(handlerLabel),
		promhttp.InstrumentHandlerResponseSize(
			responseSize.MustCurryWith(handlerLabel),
			promhttp.InstrumentHandlerCounter(
				requestsTotal.MustCurryWith(handlerLabel),
				handler,
			),
		),
	)
}
