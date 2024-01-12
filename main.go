package main

import (
	"auth-server/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	fmt.Println("Starting auth server...")
	if err := app().Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func app() *cli.App {
	return &cli.App{
		Name:  "auth-server",
		Usage: "Identity Server for authentication using JWT and OAuth2",
		Commands: []*cli.Command{
			authServerCmd(),
		},
	}
}

func authServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Starts the auth server",
		Action: func(c *cli.Context) error {
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			stopper := make(chan struct{})
			go func() {
				<-done
				close(stopper)
			}()
			server, err := server.NewServer()
			if err != nil {
				return err
			}
			return server.Start(stopper)
		},
	}
}
