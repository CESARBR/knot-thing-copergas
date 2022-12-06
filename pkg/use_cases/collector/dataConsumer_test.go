package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type conversionTest struct {
	valueToConvert  interface{}
	expectedInteger int
	expectedFloat32 interface{}
	expectedString  string
	expectedBool    bool
}

var integerConversionTests = []conversionTest{
	{valueToConvert: 528, expectedInteger: 528},
	{valueToConvert: 1821, expectedInteger: 1821},
	{valueToConvert: 890342, expectedInteger: 890342},
}

func TestGivenStringThenConvertToInteger(t *testing.T) {
	for _, test := range integerConversionTests {
		convertedInteger, err := convertValueToCorrectDataType(test.valueToConvert)
		expectedInteger := test.expectedInteger
		assert.Equal(t, expectedInteger, convertedInteger)
		assert.Nil(t, err)
	}

}

var float32ConversionTests = []conversionTest{
	{valueToConvert: 528.67, expectedFloat32: 528.67},
	{valueToConvert: 1821.8, expectedFloat32: 1821.8},
	{valueToConvert: 890342.45, expectedFloat32: 890342.45},
}

func TestGivenStringThenConvertToFloat32(t *testing.T) {
	for _, test := range float32ConversionTests {
		convertedFloat32, err := convertValueToCorrectDataType(test.valueToConvert)
		expectedFloat32 := test.expectedFloat32
		assert.Equal(t, expectedFloat32, convertedFloat32)
		assert.Nil(t, err)
	}

}

var stringConversionTests = []conversionTest{
	{valueToConvert: "data", expectedString: "data"},
	{valueToConvert: "paving stone", expectedString: "paving stone"},
}

func TestGivenStringThenReturnSameString(t *testing.T) {
	for _, test := range stringConversionTests {
		convertedValue, err := convertValueToCorrectDataType(test.valueToConvert)
		expectedValue := test.expectedString
		assert.Equal(t, expectedValue, convertedValue)
		assert.Nil(t, err)
	}
}

var boolConversionTests = []conversionTest{
	{valueToConvert: true, expectedBool: true},
	{valueToConvert: false, expectedBool: false},
}

func TestGivenBoolThenReturnSameBool(t *testing.T) {
	for _, test := range boolConversionTests {
		convertedValue, err := convertValueToCorrectDataType(test.valueToConvert)
		expectedValue := test.expectedBool
		assert.Equal(t, expectedValue, convertedValue)
		assert.Nil(t, err)
	}
}

func TestGivenEmptyStringThenReturnError(t *testing.T) {
	convertedValue, err := convertValueToCorrectDataType("")
	expectedValue := 0
	assert.Equal(t, expectedValue, convertedValue)
	assert.NotNil(t, err)
}
