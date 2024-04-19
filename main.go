package main

import (
	"example.com/scheduler/scheduler"
	"example.com/scheduler/server"
	"github.com/kardianos/service"
	"log"
)

func main() {
	svcConfig := &service.Config{
		Name:        "GoCronScheduler",
		DisplayName: "Go Cron Scheduler",
		Description: "This service runs a Go-based cron scheduler.",
	}

	// Create the Scheduler instance without the service initially
	prg := scheduler.NewScheduler("GoCronScheduler")

	// Create the service with the Scheduler
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal("Failed to create service:", err)
	}

	// Attach the service to the Scheduler for complete initialization
	prg.SetService(s)

	// Get the logger from the service
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	prg.SetLogger(logger)

	// Start the HTTP server
	go httpserver.StartServer("8080", prg.CronScheduler, func() bool {
		return prg.IsRunning()
	})

	// Run the service
	if err := s.Run(); err != nil {
		log.Println("Error:", err)
	}
}
