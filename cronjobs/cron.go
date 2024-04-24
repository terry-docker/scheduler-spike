package cronjobs

import (
	"example.com/scheduler/common"
	"example.com/scheduler/persistence"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

type ICronScheduler interface {
	Start()
	Stop()
	List()
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

var _ ICronScheduler = (*CronScheduler)(nil)

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

func (cs *CronScheduler) List() {
	entries := cs.Cron.Entries()
	for _, entry := range entries {
		fmt.Printf("ID: %v, Schedule: %v, Next Run: %v\n", entry.ID, entry.Schedule, entry.Next)
	}
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
	task()
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

func (cs *CronScheduler) LoadTasks(persistence *persistence.PersistenceManager) error {
	tasks, err := persistence.LoadTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		t := task
		if _, err := cs.AddTask(t.Spec, common.Cmd(t.VolumeName)); err != nil {
			log.Printf("Error adding task: %s", err)
		}
	}

	return nil
}
