package knot

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CESARBR/knot-thing-copergas/internal/config"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/network"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/values"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Protocol interface provides methods to handle KNoT Protocol
type Protocol interface {
	Close() error
	createDevice(device entities.Device) error
	deleteDevice(id string) error
	updateDevice(device entities.Device) error
	checkData(device entities.Device) error
	checkDeviceConfiguration(device entities.Device) error
	deviceExists(device entities.Device) bool
	generateID(device entities.Device, log *logrus.Entry) (string, error)
	checkTimeout(device entities.Device, log *logrus.Entry) entities.Device
	requestsKnot(deviceChan chan entities.Device, device entities.Device, oldState string, curState string, message string, log *logrus.Entry)
}
type networkWrapper struct {
	amqp       *network.AMQP
	publisher  network.Publisher
	subscriber network.Subscriber
}

type protocol struct {
	userToken           string
	network             *networkWrapper
	devices             map[string]entities.Device
	replyToAuthMessages string
}

func newProtocol(pipeDevices chan map[string]entities.Device, conf config.IntegrationKNoTConfig, deviceChan chan entities.Device, msgChan chan network.InMsg, log *logrus.Entry, devices map[string]entities.Device) (Protocol, error) {
	p := &protocol{}

	p.userToken = conf.UserToken
	p.replyToAuthMessages = conf.QueueName + "-auth-rpc"
	p.network = new(networkWrapper)
	p.network.amqp = network.NewAMQP(conf.URL)
	err := p.network.amqp.Start(log)
	if err != nil {
		log.Println("Knot connection error")
		return p, err
	} else {
		log.Println("Knot connected")
	}
	p.network.publisher = network.NewMsgPublisher(p.network.amqp)
	p.network.subscriber = network.NewMsgSubscriber(p.network.amqp)

	if err = p.network.subscriber.SubscribeToKNoTMessages(p.replyToAuthMessages, conf.QueueName, msgChan); err != nil {
		log.Errorln("Error to subscribe")
		return p, err
	}
	p.devices = make(map[string]entities.Device)
	p.devices = devices

	go handlerKnotAMQP(p.replyToAuthMessages, msgChan, deviceChan, log)
	go dataControl(pipeDevices, deviceChan, p, log)

	return p, nil
}

// Check for data to be updated
func (p *protocol) checkData(device entities.Device) error {

	if device.Data == nil {
		return nil
	}

	sliceSize := len(device.Data)
	nextData := 0
	for dataIndex, data := range device.Data {
		if data.Value == "" {
			return fmt.Errorf("invalid sensor value")
		} else if data.TimeStamp == "" {
			return fmt.Errorf("invalid sensor timestamp")
		}

		nextData = dataIndex + 1
		for nextData < sliceSize {
			if data.SensorID == device.Data[nextData].SensorID {
				return fmt.Errorf("repeated sensor id")
			}
			nextData++
		}
	}

	return nil
}

func isInvalidValueType(valueType int) bool {
	minValueTypeAllow := 1
	maxValueTypeAllow := 5
	return valueType < minValueTypeAllow || valueType > maxValueTypeAllow
}

// Check for device configuration
func (p *protocol) checkDeviceConfiguration(device entities.Device) error {
	sliceSize := len(device.Config)
	nextData := 0

	if device.Config == nil {
		return fmt.Errorf("sensor has no configuration")
	}

	// Check if the ids are correct, no repetition
	for dataIndex, config := range device.Config {
		if isInvalidValueType(config.Schema.ValueType) {
			return fmt.Errorf("invalid sensor id")
		}

		nextData = dataIndex + 1
		for nextData < sliceSize {
			if config.SensorID == device.Config[nextData].SensorID {
				return fmt.Errorf("repeated sensor id")
			}
			nextData++
		}
	}

	return nil

}

// Update the knot device information on map
func (p *protocol) updateDevice(device entities.Device) error {
	if _, checkDevice := p.devices[device.ID]; !checkDevice {

		return fmt.Errorf("device do not exist")
	}

	receiver := p.devices[device.ID]

	if p.checkDeviceConfiguration(device) == nil {
		receiver.Config = device.Config
	}
	if device.Name != "" {
		receiver.Name = device.Name
	}
	if device.Token != "" {
		receiver.Token = device.Token
	}
	if device.Error != "" {
		receiver.Error = device.Error
	}

	receiver.Data = nil
	oldState := receiver.State
	receiver.State = entities.KnotNew
	p.devices[device.ID] = receiver

	data, err := yaml.Marshal(&p.devices)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("internal/config/device_config.yaml", data, 0600)
	if err != nil {
		log.Fatal(err)
	}
	receiver.State = oldState
	if device.State != "" {
		receiver.State = device.State
	}
	if p.checkData(device) == nil {
		receiver.Data = device.Data
	}
	p.devices[device.ID] = receiver

	return nil
}

// Close closes the protocol.
func (p *protocol) Close() error {
	p.network.amqp.Stop()
	return nil
}

// Create a new knot device
func (p *protocol) createDevice(device entities.Device) error {

	if device.State != "" {
		return fmt.Errorf("device cannot be created, unknown source")
	} else {

		device.State = entities.KnotNew

		p.devices[device.ID] = device

		return nil
	}
}

// Create a new device ID
func (p *protocol) generateID(device entities.Device, log *logrus.Entry) (string, error) {
	delete(p.devices, device.ID)
	var err error
	device.ID, err = tokenIDGenerator()
	device.Token = ""
	p.devices[device.ID] = device

	log.Println(" generated a new Device ID : ", device.ID)

	return device.ID, err
}

// Check if the device exists
// 	return true if the device is on map
func (p *protocol) deviceExists(device entities.Device) bool {

	if _, checkDevice := p.devices[device.ID]; checkDevice {

		return true
	}
	return false
}

// Generated a new Device ID
func tokenIDGenerator() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// Just formart the Error message
func errorFormat(device entities.Device, strError string) entities.Device {
	device.Error = strError
	device.State = entities.KnotError
	return device
}

// Delete the knot device from map
func (p *protocol) deleteDevice(id string) error {
	if _, d := p.devices[id]; !d {
		return fmt.Errorf("device do not exist")
	}

	delete(p.devices, id)
	return nil
}

//non-blocking channel to update devices on the other routin
func updateDeviceMap(pipeDevices chan map[string]entities.Device, devices map[string]entities.Device) {
	pipeDevices <- devices
}

func verifyErrors(err error, log *logrus.Entry) {
	if err != nil {
		log.Errorln(err)
	}
}

//init the timeout couter
func initTimeout(deviceChan chan entities.Device, device entities.Device) {
	go func(deviceChan chan entities.Device, device entities.Device) {
		time.Sleep(20 * time.Second)
		device.Error = values.ErrorTimeout
		deviceChan <- device
	}(deviceChan, device)
}

// check response time
func (p *protocol) requestsKnot(deviceChan chan entities.Device, device entities.Device, oldState string, curState string, message string, log *logrus.Entry) {
	device.State = oldState
	initTimeout(deviceChan, device)
	device.State = curState
	err := p.updateDevice(device)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Println(message)
		switch oldState {
		case entities.KnotNew:
			err = p.network.publisher.PublishDeviceRegister(p.userToken, &device)
		case entities.KnotRegistered:
			err = p.network.publisher.PublishDeviceAuth(p.replyToAuthMessages, p.userToken, &device)
		case entities.KnotAuth:
			err = p.network.publisher.PublishDeviceUpdateConfig(p.userToken, &device)
		}
		verifyErrors(err, log)
	}
}

// Control device paths
func dataControl(pipeDevices chan map[string]entities.Device, deviceChan chan entities.Device, p *protocol, log *logrus.Entry) {
	pipeDevices <- p.devices

	for device := range deviceChan {

		if !p.deviceExists(device) {
			if device.Error != values.ErrorTimeout {
				log.Error("device id received does not match the stored")
			}
		} else {

			device = p.checkTimeout(device, log)
			if device.State != entities.KnotOff && device.Error != values.ErrorTimeout {

				err := p.updateDevice(device)
				verifyErrors(err, log)
				device = p.devices[device.ID]

				if device.Name == "" {
					log.Fatalln("Device has no name")
				} else if device.State == entities.KnotNew {
					if device.Token != "" {
						device.State = entities.KnotRegistered
					} else {
						id, err := p.generateID(device, log)
						if err != nil {
							device.State = entities.KnotOff
							log.Error(err)
						} else {
							device.ID = id
							go updateDeviceMap(pipeDevices, p.devices)
						}
					}
				}
			} else if device.Error == values.ErrorTimeout {
				device.Error = values.NoError
				log.Error("Timeout received")
			}
			switch device.State {

			// If the device status is new, request a device registration
			case entities.KnotNew:

				p.requestsKnot(deviceChan, device, device.State, entities.KnotWaitReg, "send a register request", log)

			// If the device is already registered, ask for device authentication
			case entities.KnotRegistered:

				p.requestsKnot(deviceChan, device, device.State, entities.KnotWaitAuth, "send a auth request", log)

			// The device has a token and authentication was successful.
			case entities.KnotAuth:

				p.requestsKnot(deviceChan, device, device.State, entities.KnotWaitConfig, "send a updateconfig request", log)

			//everything is ok with knot device
			case entities.KnotReady:
				device.State = entities.KnotPublishing
				err := p.updateDevice(device)
				if err != nil {
					log.Errorln(err)
				} else {
					go updateDeviceMap(pipeDevices, p.devices)
				}
			// Send the new data that comes from the device to Knot Cloud
			case entities.KnotPublishing:
				if p.checkData(device) == nil {
					log.Println("send data of device ", device.Data[0].SensorID)

					err := p.network.publisher.PublishDeviceData(p.userToken, &device, device.Data)
					if err != nil {
						log.Errorln(err)
					} else {
						device.Data = nil
						err = p.updateDevice(device)
						verifyErrors(err, log)
					}
				} else {
					log.Println("invalid data, has no data to send")
				}

			// If the device is already registered, ask for device authentication
			case entities.KnotAlreadyReg:

				var err error
				if device.Token == "" {
					device.ID, err = p.generateID(device, log)
					if err != nil {
						log.Error(err)
					} else {
						go updateDeviceMap(pipeDevices, p.devices)
						p.requestsKnot(deviceChan, device, entities.KnotNew, entities.KnotWaitReg, "send a register request", log)
					}
				} else {

					p.requestsKnot(deviceChan, device, entities.KnotRegistered, entities.KnotWaitAuth, "send a Auth request", log)

				}

			// Just delete
			case entities.KnotForceDelete:
				var err error
				log.Println("delete a device")

				device.ID, err = p.generateID(device, log)
				if err != nil {
					log.Error(err)
				} else {
					go updateDeviceMap(pipeDevices, p.devices)
					p.requestsKnot(deviceChan, device, entities.KnotNew, entities.KnotWaitReg, "send a register request", log)
				}

			// Handle errors
			case entities.KnotError:
				log.Println("ERROR: ")
				switch device.Error {
				// If the device is new to the chirpstack platform, but already has a registration in Knot, first the device needs to ask to unregister and then ask for a registration.
				case "thing's config not provided":
					log.Println("thing's config not provided")

				default:
					log.Println("ERROR WITHOUT HANDLER" + device.Error)

				}
				device.State = entities.KnotNew
				device.Error = ""
				err := p.updateDevice(device)
				verifyErrors(err, log)

			// ignore the device
			case entities.KnotOff:

			}

		}
	}
}

//check if response was received by comparing previous state with the new one
func (p *protocol) checkTimeout(device entities.Device, log *logrus.Entry) entities.Device {

	if device.Error == values.ErrorTimeout {
		curDevice := p.devices[device.ID]
		if device.State == entities.KnotNew && curDevice.State == entities.KnotWaitReg {
			return device

		} else if device.State == entities.KnotRegistered && curDevice.State == entities.KnotWaitAuth {
			return device

		} else if device.State == entities.KnotAuth && curDevice.State == entities.KnotWaitConfig {
			return device

		} else {
			device.State = entities.KnotOff
			return device
		}
	}
	return device
}

// Handle amqp messages
func handlerAMQPmessage(replyToAuthMessages string, deviceChan chan entities.Device, message network.InMsg, log *logrus.Entry, nextState string) error {
	receiver := network.DeviceGenericMessage{}
	device := entities.Device{}
	err := json.Unmarshal([]byte(string(message.Body)), &receiver)
	verifyErrors(err, log)
	device.ID = receiver.ID
	device.Name = receiver.Name
	device.Error = receiver.Error
	if network.BindingKeyRegistered == message.RoutingKey && receiver.Token != "" {
		device.Token = receiver.Token
	}

	if device.Error == values.ErrorAlreadyReg {
		device.State = entities.KnotAlreadyReg
		deviceChan <- device
		return fmt.Errorf(device.Error)

	} else if device.Error == values.ErrorFailValidation {
		device.State = entities.KnotAuth
		deviceChan <- device
		return fmt.Errorf(device.Error)

	} else if replyToAuthMessages == message.RoutingKey && device.Error != values.NoError {
		device.State = entities.KnotForceDelete
		deviceChan <- device
		return fmt.Errorf(device.Error)

	} else if device.Error != values.NoError {
		deviceChan <- errorFormat(device, device.Error)
		return fmt.Errorf(device.Error)

	} else {
		log.Println("going to ", nextState)
		device.State = nextState
		deviceChan <- device
		return nil
	}
}

// Handles messages coming from AMQP
func handlerKnotAMQP(replyToAuthMessages string, msgChan <-chan network.InMsg, deviceChan chan entities.Device, log *logrus.Entry) {

	for message := range msgChan {

		switch message.RoutingKey {

		// Registered msg from knot
		case network.BindingKeyRegistered:

			err := handlerAMQPmessage(replyToAuthMessages, deviceChan, message, log, entities.KnotRegistered)
			verifyErrors(err, log)

		// Unregistered
		case network.BindingKeyUnregistered:

			err := handlerAMQPmessage(replyToAuthMessages, deviceChan, message, log, entities.KnotForceDelete)
			verifyErrors(err, log)

		// Receive a auth msg
		case replyToAuthMessages:
			err := handlerAMQPmessage(replyToAuthMessages, deviceChan, message, log, entities.KnotAuth)
			verifyErrors(err, log)

		case network.BindingKeyUpdatedConfig:

			err := handlerAMQPmessage(replyToAuthMessages, deviceChan, message, log, entities.KnotReady)
			verifyErrors(err, log)

		}
	}
}
