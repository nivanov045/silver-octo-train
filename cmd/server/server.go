package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"

	"github.com/nivanov045/silver-octo-train/cmd/server/api"
	"github.com/nivanov045/silver-octo-train/cmd/server/service"
	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"0s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln("server::main: error in env parsing:", err)
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "address")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "store interval")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "store file")
	flag.Parse()
	log.Println("server::main: cfg:", cfg)

	storage := storage.New(cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	serv := service.New(storage)
	myapi := api.New(serv)
	log.Fatalln(myapi.Run(cfg.Address))
}
