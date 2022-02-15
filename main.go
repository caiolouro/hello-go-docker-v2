package main

import (
	"os"
	"os/signal"
	"syscall"

	// Make sure you change this line to match your module
	"github.com/caiolouro/hello-go-docker-v2/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	serverAddrFlagName string = "addr"
)

func main() {
	if err := app().Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("could not run application")
	}
}

func app() *cli.App {
	return &cli.App{
		Name:  "server",
		Usage: "The HTTP server",
		Commands: []*cli.Command{
			serverCmd(),
		},
	}
}

func serverCmd() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "starts the server",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: serverAddrFlagName, EnvVars: []string{"SERVER_ADDR"}},
		},
		Action: func(c *cli.Context) error {
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			stopper := make(chan struct{})
			go func() {
				<-done
				close(stopper)
			}()

			addr := c.String(serverAddrFlagName)
			server, err := server.NewServer(addr)
			if err != nil {
				return err
			}

			return server.Start(stopper)
		},
	}
}
