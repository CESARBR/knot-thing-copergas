package entities

// Config represents the thing's config
type Config struct {
	SensorID int    `yaml:"sensorId"`
	Schema   Schema `yaml:"schema"`
	Event    Event  `yaml:"event"`
}

type CopergasConfig struct {
	Credentials struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	Endpoints struct {
		APIUrl    string `yaml:"APIUrl"`
		AuthToken string `yaml:"authToken"`
		Variable  string `yaml:"variable"`
	}
	PertinentVariables []int  `yaml:"pertinentVariables"`
	LogFilename        string `yaml:"logFilename"`

	TimeBetweenRequestsInSeconds float32 `yaml:"timeBetweenRequestsInSeconds"`
	DataContextCache             string  `yaml:"metadataCacheFilePath"`
}

type CodVarSensorIDMapping struct {
	Mapping map[int]int `yaml:"identifier_mapping"`
}
