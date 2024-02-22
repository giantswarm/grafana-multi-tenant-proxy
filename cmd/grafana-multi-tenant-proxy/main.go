package main

import (
	"os"

	proxy "github.com/giantswarm/grafana-multi-tenant-proxy/internal/app/grafana-multi-tenant-proxy"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "Grafana Multitenant Proxy"
	app.Usage = "Makes Grafana Labs applications multi tenant"
	app.Version = version
	app.Authors = []*cli.Author{
		{Name: "Angel Barrera", Email: "angel@k8spin.cloud"},
		{Name: "Pau Rosello", Email: "pau@k8spin.cloud"},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Runs the Grafana multi tenant proxy",
			Action: proxy.Serve,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "port",
					Usage: "Port to expose this proxy",
					Value: 3501,
				}, &cli.StringFlag{
					Name:  "target-server",
					Usage: "Target server endpoint",
					Value: "http://localhost:3500",
				}, &cli.StringFlag{
					Name:  "auth-config",
					Usage: "AuthN yaml configuration file path",
					Value: "authn.yaml",
				}, &cli.StringFlag{
					Name:  "log-level",
					Usage: "Log level (DEBUG, INFO, WARN, ERROR, PANIC, FATAL)",
					Value: "INFO",
				}, &cli.BoolFlag{
					Name:  "keep-orgid",
					Usage: "Don't change OrgID header (proxy is only used for authent)",
				},
			},
		},
	}

	app.Run(os.Args)
}
