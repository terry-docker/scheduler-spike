package httpserver

import (
	"context"
	"encoding/json"
	"example.com/scheduler/scheduler"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"time"
)

var server *http.Server

func StartServer(port string, manager *scheduler.SchedulerManager) {
	mux := http.NewServeMux()

	mux.HandleFunc("/addTask", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Spec       string `json:"Spec"`
			VolumeName string `json:"VolumeName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if data.VolumeName == "" {
			http.Error(w, "Volume name not valid", http.StatusBadRequest)
			return
		}

		if err := manager.AddTask(data.Spec, data.VolumeName); err != nil {
			http.Error(w, "Failed to add task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Task added successfully"})

	})

	mux.HandleFunc("/removeTask", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			TaskID cron.EntryID `json:"TaskID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := manager.RemoveTask(data.TaskID); err != nil {
			http.Error(w, "Failed to remove task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Task removed successfully"})
	})

	mux.HandleFunc("/debug/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		manager.Scheduler.List()
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
