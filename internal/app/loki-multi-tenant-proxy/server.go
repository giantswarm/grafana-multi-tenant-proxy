package proxy

import (
	"errors"
	"log"

	"github.com/urfave/cli/v2"
)

// Serve serves
func Serve(c *cli.Context) error {
	log.Print("Error is here")
	return cli.Exit(errors.New("test fake error"), -1)
}
