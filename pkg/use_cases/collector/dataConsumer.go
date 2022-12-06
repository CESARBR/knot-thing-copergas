package collector

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot"
	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	IntType                int = 1
	FloatType              int = 2
	StringType             int = 4
	conversionErrorMessage     = "type conversion error"
)

//
// Consumes the data received from copergás.
func measurementConsumer(
	filePath string,
	pipeDevices chan map[string]entities.Device,
	measurements chan entities.ReceivedData,
	pertinentVariables []int,
	integration *knot.Integration,
	variablesIdentifier entities.CodVarSensorIDMapping,
	log *logrus.Entry) {

	var receivedMeasurements entities.Measurements = entities.Measurements{Variables: make(map[int]entities.ReceivedData)}
	var number_pertinent_variables int = len(pertinentVariables)
	var device entities.Device
	devices := <-pipeDevices
	valueType := make(map[int]entities.VariableLastData)
	err := loadDataContextFile(filePath, &valueType)
	variablesIdentifierKeys := convertMapKeysToSlice(variablesIdentifier.Mapping)
	var valueToConvert interface{}
	if err == nil {
		device = getSensorsConfiguration(devices)
		updateTimestampMapping(&valueType, variablesIdentifierKeys, device)
		log.Println("Load context of the last sent data")

	} else {
		device = getSensorsConfiguration(devices)
		updateTimestampMapping(&valueType, variablesIdentifierKeys, device)
		err = writeLatestTimestampProceed(filePath, &valueType)
		verifyErrors(err)
	}

	integration.HandleDevice(device)

	for measurement := range measurements {
		if measurement.Error == nil {
			receivedMeasurements.Variables[measurement.CodVar] = measurement
			timestamp := formatTimestamp(measurement.Data.DataLeitura)
			_, ok := valueType[measurement.CodVar]
			if ok {
				if device.State == entities.KnotPublishing && valueType[measurement.CodVar].Timestamp != measurement.Data.DataLeitura {
					//Update timestamp
					_value := entities.VariableLastData{ValueType: valueType[measurement.CodVar].ValueType, Timestamp: measurement.Data.DataLeitura}
					valueType[measurement.CodVar] = _value
					//Update data information of the data context on the file
					err = writeLatestTimestampProceed(filePath, &valueType)
					verifyErrors(err)
					knotSensorID := variablesIdentifier.Mapping[measurement.Data.CodVar]
					valueToConvert = measurement.Data.ValorConv
					//When the value is a string, the 'ValorConv' field contains a null. So we get the 'ValorString' field.
					if valueType[measurement.CodVar].ValueType == StringType {
						valueToConvert = measurement.Data.ValorString
					}
					convertedValue, err := convertValueToCorrectDataType(valueToConvert)
					if err == nil {
						data := entities.Data{SensorID: knotSensorID, Value: convertedValue, TimeStamp: timestamp}
						device.Data = append(device.Data, data)
						integration.HandleDevice(device)
						device.Data = nil
					}
				}
			}
		}

		select {
		case devices = <-pipeDevices:
			log.Println("Update context of the last sent data")
			device = getSensorsConfiguration(devices)
			updateTimestampMapping(&valueType, variablesIdentifierKeys, device)
			err = writeLatestTimestampProceed(filePath, &valueType)
			verifyErrors(err)

		default:
			// It waits to receive the data of all the requested variables before processing the received information.
			if len(receivedMeasurements.Variables) == number_pertinent_variables {
				receivedMeasurements.Variables = make(map[int]entities.ReceivedData)
			}
		}
	}
}

// Verify erros
func verifyErrors(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Formate timeStamp
func formatTimestamp(datetime string) string {
	s := strings.Split(datetime, ":")
	yyyy_mm_ddThh := s[0]
	mm := s[1]
	slice := s[2]
	sec := slice[:2]
	zone := slice[2:] + s[3]

	return yyyy_mm_ddThh + ":" + mm + ":" + sec + ".0" + zone
}

// Get the sensors configuration
// 		Where is the receiver a new configuration, this function maps the sensors types and the previous timestamp
func getSensorsConfiguration(devices map[string]entities.Device) entities.Device {
	if len(devices) < 1 {
		return entities.Device{}
	}
	keys := make([]string, 0, len(devices))
	// get all the keys of new devices
	for key := range devices {
		keys = append(keys, key)
	}
	//we just use a device in this app, always will be the [0]
	return devices[keys[0]]
}

func updateTimestampMapping(valueType *map[int]entities.VariableLastData, keys []int, device entities.Device) {
	copValueType := *valueType
	//For each device, there are many configurations. So, get this configuration and set what is each sensor type, and the last timestamp send.
	for i := 0; i < len(device.Config); i++ {
		oldTimestamp := copValueType[keys[i]].Timestamp
		valueType := device.Config[i].Schema.ValueType
		variableLastData := entities.VariableLastData{ValueType: valueType, Timestamp: oldTimestamp}
		copValueType[keys[i]] = variableLastData
	}
	*valueType = copValueType

}

func convertMapKeysToSlice(mapping map[int]int) []int {
	const length = 0
	keys := make([]int, length)
	for key := range mapping {
		keys = append(keys, key)
	}
	return keys
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		// file exists
		return true

	}
	log.Println(err)
	return false
}

//Load the last timestamps of the data send and saved
func loadDataContextFile(filePath string, context *map[int]entities.VariableLastData) error {

	if !fileExists(filePath) {
		return fmt.Errorf("file does not exists")
	}

	yamlBytes, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlBytes, context)
}

func writeLatestTimestampProceed(filePath string, context *map[int]entities.VariableLastData) error {

	data, err := yaml.Marshal(context)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(filePath), data, 0600)
}

func isEmptyString(value string) bool {
	return value == ""
}

func convertValueToCorrectDataType(value interface{}) (interface{}, error) {
	switch i := value.(type) {
	case int:
		return i, nil
	case float32:
		return i, nil
	case float64:
		return i, nil
	case string:
		if isEmptyString(i) {
			return 0, fmt.Errorf(conversionErrorMessage)
		} else {
			return i, nil
		}
	case bool:
		return i, nil
	default:
		return nil, fmt.Errorf(conversionErrorMessage)
	}
}
