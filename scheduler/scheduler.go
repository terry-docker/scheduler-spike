package scheduler

import (
	"example.com/scheduler/cronjobs"
	httpserver "example.com/scheduler/server"
	"github.com/kardianos/service"
	"sync"
)

type Scheduler struct {
	service       service.Service
	mutex         sync.Mutex
	Running       bool
	Name          string
	CronScheduler *cronjobs.CronScheduler
	logger        service.Logger
}

func NewScheduler(name string) *Scheduler {
	return &Scheduler{
		Name:          name,
		CronScheduler: cronjobs.New(),
	}
}

func (s *Scheduler) SetService(srv service.Service) {
	s.service = srv
}

func (s *Scheduler) SetLogger(logger service.Logger) {
	s.logger = logger
}

func (s *Scheduler) Start(_ service.Service) error {
	s.SetRunning(true)
	s.logger.Infof("[%s] Scheduler started with cron jobs.\n", s.Name)
	s.startCronJobs() // Initialize and start cron jobs
	return nil
}

func (s *Scheduler) Stop(_ service.Service) error {
	s.logger.Infof("[%s] Stopping scheduler\n", s.Name)
	s.CronScheduler.Stop()
	httpserver.StopServer() // Stop the HTTP server gracefully
	s.SetRunning(false)
	return nil
}

func (s *Scheduler) startCronJobs() {
	// Setup and start cron jobs
	s.CronScheduler.Start()
}

func (s *Scheduler) SetRunning(state bool) {
	s.mutex.Lock()
	s.Running = state
	s.mutex.Unlock()
}

func (s *Scheduler) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Running
}
