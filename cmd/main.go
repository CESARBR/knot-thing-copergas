package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/pkg/logging"
	"github.com/CESARBR/knot-thing-copergas/pkg/use_cases/collector"
)

func monitorSignals(sigs chan os.Signal, quit chan bool, logger logging.Logger) {
	signal := <-sigs
	logger.Infof("signal %s received", signal)
	quit <- true
}

func main() {

	logrus := logging.NewLogrus("info")
	logger := logrus.Get("Main")
	logger.Info("Starting Copergás Service")

	// Signal Handler
	sigs := make(chan os.Signal, 1)
	quit := make(chan bool, 1)
	pipeDevices := make(chan map[string]entities.Device)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go monitorSignals(sigs, quit, logger)

	collectorChan := make(chan bool, 1)

	devices, knotConfig, copergasCollector, variables_identifier_mapping, err := collector.LoadConfigs(logrus.Get("Copergas"))
	if err != nil {
		logger.Panic("Error: ", err)
		os.Exit(1)
	} else {
		logger.Info("files config ok")
	}

	knotIntegration, err := knot.NewKNoTIntegration(pipeDevices, knotConfig, logger, devices)
	if err != nil {
		logger.Panic("Error: ", err)
		os.Exit(1)
	}

	go copergasCollector.Start(pipeDevices, collectorChan, knotIntegration, variables_identifier_mapping, logger)

	for {
		select {
		case started := <-collectorChan:
			if started {
				logger.Info("Collector Started!")
			}
		case <-quit:
			os.Exit(0)
		}
	}
}
