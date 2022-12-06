package collector

import (
	"net/http"
	"sync"

	"github.com/CESARBR/knot-thing-copergas/internal/config"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/pkg/logging"
	"github.com/sirupsen/logrus"
)

// Represents the data collector.
type Collector struct {
	setup                        entities.CopergasConfig
	pertinentVariables           []int
	timeBetweenRequestsInSeconds float32
	dataChan                     chan entities.ReceivedData
	getTokenSynchronization      chan struct{}
	obtainedTokenSynchronization chan string
	httpClient                   *http.Client
	logger                       logging.Logger
}

// Creates a new collector for data from copergas.
func LoadConfigs(collectorLogger logging.Logger) (map[string]entities.Device, config.IntegrationKNoTConfig, Collector, entities.CodVarSensorIDMapping, error) {
	var copergasCollector Collector = Collector{}
	var collectorError error
	devices := make(map[string]entities.Device)
	var variables_identifier_mapping entities.CodVarSensorIDMapping

	knotSetup, err := config.LoadKnotSetup()
	if err != nil {
		collectorLogger.Errorf("Load configuration knot error: ", err)
		return devices, knotSetup, copergasCollector, variables_identifier_mapping, err
	}

	devices, err = config.LoadDeviceConfig()
	if err != nil {
		collectorLogger.Errorf("Load configuration device error: ", err)
		return devices, knotSetup, copergasCollector, variables_identifier_mapping, err
	}

	variables_identifier_mapping, err = config.LoadCodVarSensorIDMapping()
	if err != nil {
		collectorLogger.Errorf("Load identifiers mapping error: ", err)
		return devices, knotSetup, copergasCollector, variables_identifier_mapping, err
	}

	applicationSetup, err := config.LoadCopergasSetup()
	if err != nil {
		collectorLogger.Errorf("Load configuration error: ", err)
		collectorError = err
	} else {
		pertinentVariables := applicationSetup.PertinentVariables

		var timeBetweenRequestsInSeconds float32 = applicationSetup.TimeBetweenRequestsInSeconds

		data := make(chan entities.ReceivedData, len(pertinentVariables))
		getTokenSync := make(chan struct{})
		obtainedTokenSync := make(chan string)

		client := CreatesHTTPClient(timeBetweenRequestsInSeconds)

		copergasCollector.setup = applicationSetup
		copergasCollector.pertinentVariables = pertinentVariables
		copergasCollector.timeBetweenRequestsInSeconds = timeBetweenRequestsInSeconds
		copergasCollector.dataChan = data
		copergasCollector.getTokenSynchronization = getTokenSync
		copergasCollector.obtainedTokenSynchronization = obtainedTokenSync
		copergasCollector.httpClient = client
		copergasCollector.logger = collectorLogger

	}
	return devices, knotSetup, copergasCollector, variables_identifier_mapping, collectorError
}

// Starts data collection and authentication token management services.
func (c *Collector) Start(pipeDevices chan map[string]entities.Device, started chan bool, integration *knot.Integration, variables_identier entities.CodVarSensorIDMapping, logger *logrus.Entry) {

	var tokenMutex sync.Mutex
	c.logger.Info("Starting Copergás collector.")
	for _, codVar := range c.pertinentVariables {
		go GetInstantaneousMeasurement(codVar, c.timeBetweenRequestsInSeconds, c.httpClient, c.dataChan, c.setup, c.getTokenSynchronization, c.obtainedTokenSynchronization, c.logger)
	}

	go TokenHandler(c.httpClient, c.setup, c.getTokenSynchronization, c.obtainedTokenSynchronization, &tokenMutex, c.logger)

	go measurementConsumer(c.setup.DataContextCache, pipeDevices, c.dataChan, c.pertinentVariables, integration, variables_identier, logger)

	started <- true
}
