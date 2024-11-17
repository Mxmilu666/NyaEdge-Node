package main

import (
	"fmt"
	"nyaedge-node/source"
	"nyaedge-node/source/caddy"
	"nyaedge-node/source/server"
	"nyaedge-node/source/zaplogger"
	"os"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("NyaEdge-Node v0.0.1")
	logger, err := zaplogger.Setup()
	if err != nil {
		return
	}

	configFile := "config.yml"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Info("No config file")
		return
	}

	Config, err := source.ReadConfig(configFile)
	if err != nil {
		logger.Error("Error reading config file: %v", zap.Error(err))
		return
	}

	logger.Debug("Connecting to Caddy")
	connected, err := caddy.CheckCaddyAPI()
	if err != nil {
		logger.Error("Error checking Caddy API", zap.Error(err))
		return
	}
	if !connected {
		logger.Error("Caddy not running")
		return
	}

	logger.Info("Setup Caddy Server")

	if err := server.StartServer(Config, logger); err != nil {
		logger.Error("Server startup failed: %v", zap.Error(err))
		return
	}

}
