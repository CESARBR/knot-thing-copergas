package knot

import (
	"github.com/CESARBR/knot-thing-copergas/internal/config"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/network"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Integration implements an KNoT integration.
type Integration struct {
	protocol Protocol
}

var deviceChan = make(chan entities.Device)
var msgChan = make(chan network.InMsg)

// New creates a new KNoT integration.
func NewKNoTIntegration(pipeDevices chan map[string]entities.Device, conf config.IntegrationKNoTConfig, log *logrus.Entry, devices map[string]entities.Device) (*Integration, error) {
	var err error
	KNoTInteration := Integration{}

	KNoTInteration.protocol, err = newProtocol(pipeDevices, conf, deviceChan, msgChan, log, devices)
	if err != nil {
		return nil, errors.Wrap(err, "new knot protocol")
	}

	return &KNoTInteration, nil
}

// HandleUplinkEvent sends an UplinkEvent.
func (i *Integration) HandleDevice(device entities.Device) {
	device.State = ""
	deviceChan <- device
}

// Close closes the integration.
func (integration *Integration) Close() error {
	return integration.protocol.Close()
}
