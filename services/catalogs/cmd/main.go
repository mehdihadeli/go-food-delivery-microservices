package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/store-golang-microservice-sample/pkg/must"
	"github.com/store-golang-microservice-sample/pkg/shutdown"
	"io"
	nethttp "net/http"
	"os"
	"time"

	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Order Application
// @description catalogs service
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1
func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	cleanup, err := run(os.Stdout)
	defer cleanup()

	if err != nil {
		fmt.Printf("%v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func run(w io.Writer) (func(), error) {
	server, err := buildServer(w)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := server.Start(); err != nil && err != nethttp.ErrServerClosed {
			server.Fatal(errors.New("server could not be started"))
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(server.Config().Context.Timeout)*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Fatal(err)
		}
	}, nil
}

func buildServer(w io.Writer) (*http.Server, error) {
	var cfg http.Config
	readConfig(&cfg)


}

func readConfig(cfg *http.Config) {
	viper.SetConfigFile(`./config.json`)

	must.NotFailF(viper.ReadInConfig)
	must.NotFail(viper.Unmarshal(cfg))
}
