package forge

import (
	"errors"
	"fmt"
)

type Primative struct {
	valueType ValueType
	value     interface{}
}

func NewPrimative(valueType ValueType, value interface{}) *Primative {
	return &Primative{
		valueType: valueType,
		value:     value,
	}
}

func NewBoolean(value bool) *Primative {
	return NewPrimative(BOOLEAN, value)
}

func NewFloat(value float64) *Primative {
	return NewPrimative(FLOAT, value)
}

func NewInteger(value int64) *Primative {
	return NewPrimative(INTEGER, value)
}

func NewNull() *Primative {
	return NewPrimative(NULL, nil)
}

func NewString(value string) *Primative {
	return NewPrimative(STRING, value)
}

func (primative *Primative) GetType() ValueType {
	return primative.valueType
}

func (primative *Primative) GetValue() interface{} {
	return primative.value
}

func (primative *Primative) UpdateValue(value interface{}) error {
	// Valid types
	switch value.(type) {
	case bool:
		primative.valueType = BOOLEAN
	case float64:
		primative.valueType = FLOAT
	case int64:
		primative.valueType = INTEGER
	case nil:
		primative.valueType = NULL
	case string:
		primative.valueType = STRING
	default:
		msg := fmt.Sprintf("Unsupported type, %s must be of (bool, float64, int64, nil, string)", value)
		return errors.New(msg)

	}
	primative.value = value
	return nil
}

func (primative *Primative) AsBoolean() (bool, error) {
	return asBoolean(primative.value)
}

func (primative *Primative) AsFloat() (float64, error) {
	return asFloat(primative.value)
}

func (primative *Primative) AsInteger() (int64, error) {
	return asInteger(primative.value)
}

func (primative *Primative) AsString() (string, error) {
	return asString(primative.value)
}

func (primative *Primative) String() string {
	str, _ := primative.AsString()
	return str
}
