package main

import (
	"context"
	"fmt"
	"log"
	"task/config"
	"task/internal/server"
	"task/pkg/db/redis"
	"task/pkg/logger"
)

const configPath = "./config/config-local"

func main() {
	log.Println("Starting api server")
	fmt.Println()

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s", cfg.Server.AppVersion, cfg.Logger.Level)

	db := redis.NewRedisClient(cfg)
	_, err = db.Ping(context.Background()).Result()

	if err != nil {
		appLogger.Fatalf("Redis init: %s", err)
		return
	}

	defer db.Close()

	appLogger.Infof("Postgres connected, Status: %#v", db.PoolStats())

	s := server.NewServer(cfg, db, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatalf("Server error: %s", err.Error())
	}
}
