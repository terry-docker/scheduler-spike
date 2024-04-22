package cronjobs

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

type Scheduler interface {
	Start()
	Stop()
	AddTask(spec string, task func()) (cron.EntryID, error)
	RemoveTask(id cron.EntryID) error
}

type CronTask struct {
	ID   cron.EntryID
	Spec string
	Task func()
}

type CronScheduler struct {
	Cron  *cron.Cron
	Tasks map[cron.EntryID]CronTask
}

var _ Scheduler = (*CronScheduler)(nil)

func New() *CronScheduler {
	return &CronScheduler{
		Cron:  cron.New(),
		Tasks: make(map[cron.EntryID]CronTask),
	}
}

func (cs *CronScheduler) Start() {
	cs.Cron.Start()
}

func (cs *CronScheduler) Stop() {
	cs.Cron.Stop()
}

func (cs *CronScheduler) AddTask(spec string, task func()) (cron.EntryID, error) {
	id, err := cs.Cron.AddFunc(spec, task)
	if err != nil {
		return 0, err
	}
	cs.Tasks[id] = CronTask{
		ID:   id,
		Spec: spec,
		Task: task,
	}
	log.Printf("Added new task: %v with spec: %s", id, spec)
	return id, nil
}

func (cs *CronScheduler) RemoveTask(id cron.EntryID) error {
	if _, exists := cs.Tasks[id]; !exists {
		return fmt.Errorf("task with ID %d does not exist", id)
	}

	cs.Cron.Remove(id)
	delete(cs.Tasks, id)
	log.Printf("Removed task: %v", id)
	return nil
}
