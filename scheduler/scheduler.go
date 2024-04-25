package scheduler

import (
	"fmt"
	"log"

	"example.com/scheduler/common"
	"example.com/scheduler/cronjobs"
	"example.com/scheduler/persistence"
	"github.com/google/uuid"
	"github.com/kardianos/service"
	"github.com/robfig/cron/v3"
)

type SchedulerManager struct {
	Scheduler   *cronjobs.CronScheduler
	Persistence *persistence.Manager
	service     service.Service
	logger      service.Logger
	tasks       []TaskMap
}

type TaskMap struct {
	volumeName string
	cronID     cron.EntryID
}

func NewSchedulerManager(scheduler *cronjobs.CronScheduler, persistence *persistence.Manager) *SchedulerManager {
	return &SchedulerManager{
		Scheduler:   scheduler,
		Persistence: persistence,
	}
}

func (sm *SchedulerManager) AddTask(spec string, volumeName string) error {
	// Add task to the scheduler
	_, err := sm.Scheduler.AddTask(spec, common.Cmd(volumeName))
	if err != nil {
		return err
	}

	// Record new task
	// TODO Check pinata for existing uuid solutions.
	task := common.TaskConfig{
		ID:         fmt.Sprint(uuid.New()),
		Spec:       spec,
		VolumeName: volumeName,
	}

	tasks, err := sm.Persistence.LoadTasks()
	if err != nil {
		return err
	}
	tasks = append(tasks, task)

	// Persist new task configuration
	return sm.Persistence.SaveTasks(tasks)
}

// TODO need to remove the task from json, was working when I was storing the cron id but I moved to a uuid and broke it.
// We don't need to store the cron id's as it changes between runs, leaning on storing a key/value map to associate the TaskConfig with the cron id.
func (sm *SchedulerManager) RemoveTask(taskID cron.EntryID) error {
	sm.Scheduler.RemoveTask(taskID)

	// Load, filter out removed task, save updated list
	tasks, err := sm.Persistence.LoadTasks()
	if err != nil {
		return err
	}
	filteredTasks := []common.TaskConfig{}
	for _, task := range tasks {
		if task.ID != fmt.Sprint(taskID) {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return sm.Persistence.SaveTasks(filteredTasks)
}

func (sm *SchedulerManager) RunScheduler() error {
	// TODO update this
	svcConfig := &service.Config{
		Name:        "GoCronScheduler",
		DisplayName: "Go Cron CronScheduler",
		Description: "This service runs a Go-based cron scheduler.",
	}

	// Load tasks from the persistence storage
	if err := sm.Scheduler.LoadTasks(sm.Persistence); err != nil {
		log.Fatalf("Failed to load tasks: %v", err)
		return err
	}

	// Create the service with the CronScheduler
	s, err := service.New(sm, svcConfig)
	if err != nil {
		log.Fatal("Failed to create service:", err)
		return err
	}

	// Attach the service to the CronScheduler for complete initialization
	sm.SetService(s)

	// Get the logger from the service
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal("Failed to create logger:", err)
		return err
	}
	sm.SetLogger(logger)

	// Run the service
	if err := s.Run(); err != nil {
		log.Println("Error:", err)
		return err
	}
	return nil
}

func (sm *SchedulerManager) SetService(srv service.Service) {
	sm.service = srv
}

func (sm *SchedulerManager) SetLogger(logger service.Logger) {
	sm.logger = logger
}

func (sm *SchedulerManager) Start(s service.Service) error {
	sm.logger.Infof("Scheduler started with cron jobs.\n")
	sm.Scheduler.Start()

	return nil
}

func (sm *SchedulerManager) Stop(_ service.Service) error {
	sm.logger.Infof("[%s] Stopping scheduler\n")
	sm.service.Stop()
	return nil
}
