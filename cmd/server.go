package cmd

import (
	"app/internal/infras/app"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Load environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// serverCommand.Flags().IntVarP(&serverCfg.Port, "port", "p", serverCfg.Port, "port number to run on")
	// serverCommand.Flags().StringVar(&serverCfg.ApiPrefix, "prefix", serverCfg.ApiPrefix, "API prefix")
	rootCmd.AddCommand(serverCommand)
}

var (
	err       error

	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "Run App API Server",
		Run: func(cmd *cobra.Command, args []string) {
			server := app.NewApp()

			// Create a signal channel and start server in a goroutine
			signChan := make(chan os.Signal, 1)
			go func() {
				if err := server.Start(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("server.Start: %v", err)
				}
			}()

			// Listen to terminate signal from system or os
			signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
			<-signChan
			log.Println("shutting down")

			waitTime := 2 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), waitTime)
			defer func() {
				log.Println("Closing Database connection!")
				cancel()
			}()

			// Stop server
			log.Println("Stopping http server....")
			if err := server.Stop(ctx); err != nil {
				log.Println("Halted active connections")
			}
			select {
			case <-ctx.Done():
				log.Printf("timeout for %v seconds", waitTime)
			}
			log.Printf("Shudown completed!")
		},
	}
)

func LoadConfig() (app.Opts, error) {
	c := app.Opts{}

	if err := env.Parse(&c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	log.Printf("%+v\n", c)
	return c, nil
}
