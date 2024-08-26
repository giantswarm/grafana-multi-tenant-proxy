package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/config"
	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/handler/auth"
)

type Proxy struct {
	proxies     map[string]*httputil.ReverseProxy
	proxyConfig *config.ProxyConfig
	logger      *zap.Logger
	errorLogger *log.Logger
}

func NewProxy(config *config.Config, logger *zap.Logger, errorLogger *log.Logger) Proxy {
	return Proxy{
		proxies:     make(map[string]*httputil.ReverseProxy),
		proxyConfig: &config.Proxy,
		logger:      logger,
		errorLogger: errorLogger,
	}
}

func (p Proxy) ApplyConfig(config *config.Config) {
	*p.proxyConfig = config.Proxy
	// Clear the existing proxies
	for k := range p.proxies {
		delete(p.proxies, k)
	}
}

// Create the reverse proxy handler to target server
func (p Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if proxy, ok := p.proxies[r.Host]; ok {
			proxy.ServeHTTP(w, r)
			return
		}

		// Find the target server for the host
		server := p.proxyConfig.FindTargetServer(r.Host)
		if server != nil {
			proxy := p.newProxy(server.Target)
			p.proxies[r.Host] = proxy
			proxy.ServeHTTP(w, r)
			return
		}

		p.errorLogger.Print("Target server not configured")
		w.WriteHeader(404)
		_, err := w.Write([]byte("Not Found\n"))
		if err != nil {
			p.logger.Error("Could not write response", zap.Error(err))
		}
	}
}

func (p Proxy) newProxy(targetServer string) *httputil.ReverseProxy {
	targetServerURL, err := url.Parse(targetServer)
	if err != nil {
		log.Println("target parse fail:", err)
	}
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetServerURL)
			r.Out.Host = targetServerURL.Host
			r.Out.Header.Add("X-Forwarded-Host", targetServerURL.Host)

			orgID := r.In.Context().Value(auth.OrgIDKey)
			if orgID != "" {
				r.Out.Header.Set("X-Scope-OrgID", orgID.(string))
			}
		},
		ErrorLog: p.errorLogger,
	}
}
