package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/auth"
	"github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy/config"
)

type Proxy struct {
	proxies     map[string]*httputil.ReverseProxy
	proxyConfig *config.ProxyConfig
	errorLogger *log.Logger
}

func NewProxy(config config.Config, errorLogger *log.Logger) Proxy {
	return Proxy{
		proxies:     make(map[string]*httputil.ReverseProxy),
		proxyConfig: &config.Proxy,
		errorLogger: errorLogger,
	}
}

func (p Proxy) ApplyConfig(config config.Config) {
	*p.proxyConfig = config.Proxy
}

// Create the reverse proxy handler to target server
func (p Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.Host

		if proxy, ok := p.proxies[host]; ok {
			proxy.ServeHTTP(w, r)
			return
		}

		for _, v := range p.proxyConfig.TargetServers {
			if v.Host == host {
				proxy := p.newProxy(v.Target)
				p.proxies[host] = proxy
				proxy.ServeHTTP(w, r)
				return
			}
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
			r.Out.Header.Set("X-Forwarded-Host", targetServerURL.Host)

			orgID := r.In.Context().Value(auth.OrgIDKey)
			if orgID != "" {
				r.Out.Header.Set("X-Scope-OrgID", orgID.(string))
			}
		},
		ErrorLog: p.errorLogger,
	}
}
