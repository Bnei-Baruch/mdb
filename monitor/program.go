package monitor

import (
	"github.com/Bnei-Baruch/mdb/monitor/agent"
	"github.com/Bnei-Baruch/mdb/monitor/config"
	log "github.com/Sirupsen/logrus"
	"github.com/kardianos/service"
	"github.com/spf13/viper"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		log.Info("Monitoring program running in terminal.")
	} else {
		log.Info("Monitoring program running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	log.Info("Monitoring program stopping...")
	close(p.exit)
	return nil
}

func (p *program) run() error {
	log.Infof("Monitoring program running %v.", service.Platform())

	/*
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case tm := <-ticker.C:
				log.Infof("Monitoring program still running at %v...", tm)
			case <-p.exit:
				ticker.Stop()
				return nil
			}
		}
	*/

	// If no other options are specified, load the config file and run.
	c := config.NewConfig()
	err := c.LoadConfig(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal("Error during loading configuration file: " + err.Error())
	}

	if len(c.Inputs) == 0 {
		log.Fatalf("Error: no inputs found, did you provide a valid config file?")
	}

	if len(c.Outputs) == 0 {
		log.Fatalf("Error: no outputs found, did you provide a valid config file?")
	}

	if int64(c.Agent.Interval) <= 0 {
		log.Fatalf("Error monitoring agent interval must be positive; found %s",
			c.Agent.Interval)
	}

	if int64(c.Agent.FlushInterval) <= 0 {
		log.Fatalf("Error monitoring agent flush_interval must be positive; found %s",
			c.Agent.FlushInterval)
	}

	ag, err := agent.NewAgent(c)
	if err != nil {
		log.Fatal("Error occurred during initializing new monitoring agent : " + err.Error())
	}

	err = ag.Connect()
	if err != nil {
		log.Fatal("Error occurred during connection in monitoring agent" + err.Error())
	}

	shutdown := make(chan struct{})

	// TODO : Impelement configuration watcher and monitor reloading logic here

	ag.Run(shutdown)

	return nil
}
