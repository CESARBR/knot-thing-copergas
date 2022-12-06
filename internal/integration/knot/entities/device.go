package entities

// States that represent the current status of the device on the Knot network
const (
	KnotNew         string = "new"
	KnotAlreadyReg  string = "alreadyRegistered"
	KnotRegistered  string = "registered"
	KnotForceDelete string = "forceDelete"
	KnotReady       string = "readyToSendData"
	KnotPublishing  string = "SendData"
	KnotAuth        string = "authenticated"
	KnotError       string = "error"
	KnotWaitReg     string = "waitResponseRegister"
	KnotWaitAuth    string = "waitResponseAuth"
	KnotWaitConfig  string = "waitResponseConfig"
	KnotOff         string = "ignore"
)

// Device represents the device domain entity
type Device struct {
	// KNoT Protocol properties
	ID     string   `yaml:"id"`
	Token  string   `yaml:"token"`
	Name   string   `yaml:"name"`
	Config []Config `yaml:"config"`
	State  string   `yaml:"state"`
	Data   []Data   `yaml:"data"`
	Error  string

	// LoRaWAN properties

	// KNoT Protocol status
	// status string enum(ready ou online, register, auth, config)
}

type CopergasDevice struct {
	ID      int
	Name    string
	Sensors []Sensor
}

type Sensor struct {
	ID    int
	Name  string
	Value float32
}
