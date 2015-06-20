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

func (this *Primative) GetType() ValueType {
	return this.valueType
}

func (this *Primative) GetValue() interface{} {
	return this.value
}

func (this *Primative) UpdateValue(value interface{}) error {
	// Valid types
	switch value.(type) {
	case bool:
		this.valueType = BOOLEAN
	case float64:
		this.valueType = FLOAT
	case int64:
		this.valueType = INTEGER
	case nil:
		this.valueType = NULL
	case string:
		this.valueType = STRING
	default:
		msg := fmt.Sprintf("Unsupported type, %s must be of (bool, float64, int64, nil, string)", value)
		return errors.New(msg)

	}
	this.value = value
	return nil
}

func (this *Primative) AsBoolean() (bool, error) {
	return asBoolean(this.value)
}

func (this *Primative) AsFloat() (float64, error) {
	return asFloat(this.value)
}

func (this *Primative) AsInteger() (int64, error) {
	return asInteger(this.value)
}

func (this *Primative) AsString() (string, error) {
	return asString(this.value)
}

func (this *Primative) String() string {
	str, _ := this.AsString()
	return str
}
