package main

import (
	"cms/controller"
	"cms/utils"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	log := utils.NewLogger()
	config, err := utils.NewConfig()
	if err != nil {
		log.WithError(err).Error("unable to load env")
		return
	}
	db := utils.ConnectDB(log, config)

	grpcServer := controller.NewServer(db, log)

	lis, err := net.Listen("tcp", config.GetPort())
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	// Start the gRPC server
	go func() {
		log.Info(fmt.Sprintf("Starting gRPC server on %s", config.GetPort()))
		if err := grpcServer.Serve(lis); err != nil {
			log.WithError(err).Fatal("Failed to serve")
		}
	}()

	// Graceful shutdown of the server on interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info("Server stopped gracefully.")
}
