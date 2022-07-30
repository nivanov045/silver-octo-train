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
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func main() {
	var cfg config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "store interval")
	flag.BoolVar(&cfg.Restore, "r", true, "restore")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "store file")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln("server::main: error in env parsing:", err)
	}
	log.Println("server::main: cfg:", cfg)

	storage := storage.New(cfg.StoreInterval, cfg.StoreFile, cfg.Restore)
	serv := service.New(storage)
	myapi := api.New(serv)
	log.Fatalln(myapi.Run(cfg.Address))
}
