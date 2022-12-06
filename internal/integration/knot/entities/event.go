package entities

// Event represents the thing's event
type Event struct {
	Change         bool        `yaml:"change"`
	TimeSec        int         `yaml:"timeSec"`
	LowerThreshold interface{} `yaml:"lowerThreshold"`
	UpperThreshold interface{} `yaml:"upperThreshold"`
}
