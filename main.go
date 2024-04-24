package main

import (
	"example.com/scheduler/cronjobs"
	"example.com/scheduler/persistence"
	"example.com/scheduler/scheduler"
	"example.com/scheduler/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}

	// Construct the file path
	filePath := filepath.Join(cwd, "test.json")

	persistenceManager := persistence.NewPersistenceManager(filePath)
	cronScheduler := cronjobs.New()
	schedulerManager := scheduler.NewSchedulerManager(cronScheduler, persistenceManager)

	// Start the HTTP server
	go httpserver.StartServer("8080", schedulerManager)

	// Run the scheduler
	if err := schedulerManager.RunScheduler(); err != nil {
		log.Fatalf("Failed to run scheduler: %v", err)
	}

}
