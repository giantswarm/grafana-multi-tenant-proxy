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

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/auth"
	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/pkg"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var keepOrgID bool
var authConfigLocation string
var authConfig *pkg.Authn

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
			NativeHistogramBucketFactor:     1.1,
			NativeHistogramMaxBucketNumber:  100,
			NativeHistogramMinResetDuration: 1 * time.Hour,
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

func loadConfig() (*pkg.Authn, error) {
	config, err := pkg.ParseConfig(&authConfigLocation)
	config.KeepOrgID = keepOrgID
	return config, err
}

// Serve serves
func Serve(c *cli.Context) error {
	targetServerURL, _ := url.Parse(c.String("target-server"))
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
		logger,
		ReverseTarget(reverseProxy),
		*authConfig,
	)

	handlers := Logger(
		authenticationMiddleware.Authenticate(),
		logger,
	)

	// Register Prometheus collectors
	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	// We handle metrics first to avoid calling the authentication middleware
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", instrumentHandler("authentication", handlers))
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
