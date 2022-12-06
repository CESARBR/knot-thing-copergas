package entities

// Schema represents the thing's schema
type Schema struct {
	ValueType int    `yaml:"valueType"`
	Unit      int    `yaml:"unit"`
	TypeID    int    `yaml:"typeId"`
	Name      string `yaml:"name"`
}
