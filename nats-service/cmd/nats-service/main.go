package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/internal/inmemory"
	"main/internal/log"
	"main/internal/subscriber"
)

func main() {
	cfg := config.GetConfig()
	log.ConfigureLogger(cfg.LogLevel)
	memoryCache := inmemory.NewCache()

	subscriberService := subscriber.NewSubscriber(&cfg, memoryCache)

	subscriberService.RestoreCache()
	go subscriberService.Run()

	application := app.NewApp(&cfg, memoryCache)
	application.Run()
}
