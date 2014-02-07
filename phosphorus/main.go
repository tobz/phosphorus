package main

import "os"
import "os/signal"
import "github.com/tobz/phosphorus/log"
import "github.com/tobz/phosphorus/server"
import _ "github.com/tobz/phosphorus/network/all"
import _ "github.com/tobz/phosphorus/database/all"

func main() {
	log.Server.Info("main", "Phosphorus is starting...")

	// Pull in our configuration.
	config := &server.Config{}
	config, err := server.NewConfig("phosphorus.yaml")
	if err != nil {
		log.Server.Error("main", "Couldn't load configuration: %s", err)
		os.Exit(1)
	}

	// Create the server and start 'er up.
	server := server.NewServer(config)
	err = server.Start()
	if err != nil {
		log.Server.Error("main", "Caught an error while starting server: %s", err)
		os.Exit(1)
	}

	// Spin until we get a signal to stop.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	for {
		select {
		case sig := <-sigChan:
			log.Server.Warn("main", "Caught %s signal.  Starting shutdown.", sig)
			server.Stop()
			log.Server.Info("main", "Shutdown complete.  Exiting.")
			os.Exit(0)
		}
	}
}
