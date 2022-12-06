package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/pkg/logging"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Server represents the server configuration properties
type Server struct {
	Port int
}

// Logger represents the logger configuration properties
type Logger struct {
	Level string
}

// Config represents the service configuration
type Config struct {
	Server
	Logger
}

func readFile(name string) {
	logger := logging.NewLogrus("error").Get("Config")
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("Error reading config file, %s", err)
	}
}

// IntegrationKNoTConfig holds the KNoT integration configuration.
type IntegrationKNoTConfig struct {
	UserToken               string `yaml:"user_token"`
	URL                     string `yaml:"url"`
	EventRoutingKeyTemplate string `yaml:"event_routing_key_template"`
	QueueName               string `yaml:"AMQPqueue"`
}

// Load returns the service configuration
func Load() Config {
	var configuration Config
	logger := logging.NewLogrus("info").Get("Config")
	viper.AddConfigPath("internal/config")
	viper.SetConfigType("yaml")

	if os.Getenv("ENV") == "development" {
		readFile("development")
		if err := viper.MergeInConfig(); err != nil {
			logger.Fatalf("Error reading config file, %s", err)
		}
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&configuration); err != nil {
		logger.Fatalf("Error unmarshalling configuration, %s", err)
	}

	return configuration
}

func LoadCopergasSetup() (entities.CopergasConfig, error) {
	var config entities.CopergasConfig
	var configErr error = nil

	yamlBytes, err := ioutil.ReadFile("internal/config/copergas_setup.yaml")
	if err != nil {
		configErr = err
	}

	unmarshalErr := yaml.Unmarshal(yamlBytes, &config)
	if unmarshalErr != nil {
		configErr = unmarshalErr
	}
	return config, configErr

}

func LoadKnotSetup() (IntegrationKNoTConfig, error) {
	var config IntegrationKNoTConfig
	var configErr error = nil

	yamlBytes, err := ioutil.ReadFile("internal/config/knot_setup.yaml")
	if err != nil {
		configErr = err
	}

	unmarshalErr := yaml.Unmarshal(yamlBytes, &config)
	if unmarshalErr != nil {
		configErr = unmarshalErr
	}
	return config, configErr
}

func LoadDeviceConfig() (map[string]entities.Device, error) {
	var config map[string]entities.Device
	var configErr error = nil

	yamlBytes, err := ioutil.ReadFile("internal/config/device_config.yaml")
	if err != nil {
		configErr = err
	}

	unmarshalErr := yaml.Unmarshal(yamlBytes, &config)
	if unmarshalErr != nil {
		configErr = unmarshalErr
	}
	return config, configErr
}

func LoadCodVarSensorIDMapping() (entities.CodVarSensorIDMapping, error) {
	var mapping entities.CodVarSensorIDMapping
	var mappingErr error
	yamlBytes, err := ioutil.ReadFile("internal/config/copergas_identifier_knot_sensor_mapping.yaml")
	if err != nil {
		mappingErr = err
	}
	unmarshalErr := yaml.Unmarshal(yamlBytes, &mapping)
	if unmarshalErr != nil {
		mappingErr = unmarshalErr
	}
	return mapping, mappingErr

}
