package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiolouro/hello-go-docker-v2/server"
	"github.com/caiolouro/hello-go-docker-v2/storage"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	serverAddrFlagName          string = "addr"
	apiServerStorageDatabaseURL string = "database-url"
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
			&cli.StringFlag{Name: apiServerStorageDatabaseURL, EnvVars: []string{"DATABASE_URL"}},
		},
		Action: func(c *cli.Context) error {
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			stopper := make(chan struct{})
			go func() {
				<-done
				close(stopper)
			}()

			databaseURL := c.String(apiServerStorageDatabaseURL)
			s, err := storage.NewStorage(databaseURL)
			if err != nil {
				return fmt.Errorf("could not initialize storage: %w", err)
			}

			addr := c.String(serverAddrFlagName)
			server, err := server.NewServer(addr, s)
			if err != nil {
				return err
			}

			return server.Start(stopper)
		},
	}
}
