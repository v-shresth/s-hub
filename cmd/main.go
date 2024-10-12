package main

import (
	"cms/clients"
	"cms/server"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	log := clients.NewLogger()
	config, err := clients.NewConfig()
	if err != nil {
		log.WithError(err).Error("unable to load env")
		return
	}
	systemDB := clients.GetSystemDB(log, config)

	grpcServer := server.NewServer(systemDB, log, config)

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
