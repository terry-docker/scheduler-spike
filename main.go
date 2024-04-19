package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/kardianos/service"
)

type program struct {
	tag      string
	interval int32
	exit     chan struct{}
	logger   service.Logger // Use this logger instead of the global one
}

func (p *program) Start(s service.Service) error {
	// Use p.logger here to ensure you're using the correct logger instance
	if service.Interactive() {
		p.logger.Infof("[%s] Running in terminal.", p.tag)
	} else {
		p.logger.Infof("[%s] Running under service manager.", p.tag)
	}
	p.exit = make(chan struct{})

	go p.run()
	return nil
}

func (p *program) run() {
	// Use p.logger here
	p.logger.Infof("[%s] I'm running %v.", p.tag, service.Platform())
	interval := time.Duration(p.interval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case tm := <-ticker.C:
			p.logger.Infof("[%s] Still running at %v...", p.tag, tm)
		case <-p.exit:
			p.logger.Infof("[%s] Stopping", p.tag)
			return
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Use p.logger here
	p.logger.Infof("[%s] Stopping", p.tag)
	close(p.exit)
	return nil
}

func main() {
	wg := new(sync.WaitGroup)

	oranges, ologger := createService("Orange", 2)
	pairs, plogger := createService("Pair", 5)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := oranges.Run(); err != nil {
			ologger.Error(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := pairs.Run(); err != nil {
			plogger.Error(err)
		}
	}()

	wg.Wait()
}

func createService(name string, interval int32) (service.Service, service.Logger) {
	svcFlag := flag.String(name, "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        name,
		DisplayName: name + " Go Service",
		Description: "This is an example Go service that logs at given intervals.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	prg := &program{
		tag:      name,
		interval: interval,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	// Assign the service logger after it's been created
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	prg.logger = logger // Correct assignment of the logger to the program

	if len(*svcFlag) != 0 {
		if err := service.Control(s, *svcFlag); err != nil {
			log.Fatal(err)
		}
	}

	return s, logger
}
