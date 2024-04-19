package httpserver

import (
	"context"
	"encoding/json"
	"example.com/scheduler/cronjobs"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"time"
)

var server *http.Server

func StartServer(port string, scheduler *cronjobs.CronScheduler, isRunning func() bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if isRunning() {
			fmt.Fprintln(w, "Service is running")
		} else {
			fmt.Fprintln(w, "Service is not running")
		}
	})

	mux.HandleFunc("/addTask", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec        string `json:"Spec"`
			Description string `json:"Description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		taskFunc := func() {
			log.Printf("Executing task: %s", data.Description)
		}

		if id, err := scheduler.AddTask(data.Spec, taskFunc); err != nil {
			http.Error(w, "Failed to add task", http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"taskID": id})
		}
	})

	mux.HandleFunc("/removeTask", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			TaskID cron.EntryID `json:"TaskID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		scheduler.RemoveTask(data.TaskID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Task removed successfully")
	})

	server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on port %s\n", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func StopServer() {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		}
		log.Println("HTTP server stopped")
	}
}
