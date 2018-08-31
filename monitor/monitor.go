package monitor

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/service"
)

// ProcessMonitoring - start and process BB Archive mdb monitoring service
func ProcessMonitoring() {
	var err error
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "BB Archive MDB Monitoring Service",
		DisplayName: "BB Archive MDB Monitoring Service",
		Description: "This is a Go service that collects, monitors and analyses BB archive mdb data, checks integrity and alerts about abnormal behaviour.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	log.Info("Starting to monitor mdb...")

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		log.Error(err)
	}

	log.Info("Monitoring mdb successfully finished...")
}
