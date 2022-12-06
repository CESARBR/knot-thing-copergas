package config

import (
	"testing"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidMappingSuccess(t *testing.T) {
	const testCodVar = 79565
	codVarSensorIDMapping := map[int]int{
		79565: 1,
	}

	expected_mapping := entities.CodVarSensorIDMapping{Mapping: codVarSensorIDMapping}
	sensorIDMapping, err := LoadCodVarSensorIDMapping()
	assert.Nil(t, err)
	assert.Equal(t, expected_mapping.Mapping[testCodVar], sensorIDMapping.Mapping[testCodVar])

}
