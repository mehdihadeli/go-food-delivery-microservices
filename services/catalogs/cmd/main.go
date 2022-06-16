package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/prometheus/common/config"
	"log"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("WriterService")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
