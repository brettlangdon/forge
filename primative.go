package forge

import (
	"errors"
	"fmt"
	"math"
	"strconv"
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
	switch val := primative.value.(type) {
	case bool:
		return val, nil
	case float64:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case nil:
		return false, nil
	case string:
		return val != "", nil
	}

	msg := fmt.Sprintf("Could not convert value %s to type BOOLEAN", primative.value)
	return false, errors.New(msg)
}

func (primative *Primative) AsFloat() (float64, error) {
	switch val := primative.value.(type) {
	case bool:
		floatVal := float64(0)
		if val {
			floatVal = float64(1)
		}

		return floatVal, nil
	case float64:
		return val, nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	}

	msg := fmt.Sprintf("Could not convert value %s to type FLOAT", primative.value)
	return 0, errors.New(msg)
}

func (primative *Primative) AsInteger() (int64, error) {
	switch val := primative.value.(type) {
	case bool:
		intVal := int64(0)
		if val {
			intVal = int64(1)
		}
		return intVal, nil
	case float64:
		return int64(math.Trunc(val)), nil
	case int64:
		return val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	}

	msg := fmt.Sprintf("Could not convert value %s to type INTEGER", primative.value)
	return 0, errors.New(msg)
}

func (primative *Primative) AsNull() (interface{}, error) {
	switch val := primative.value.(type) {
	case nil:
		return val, nil
	}

	msg := fmt.Sprintf("Could not convert value %s to nil", primative.value)
	return 0, errors.New(msg)
}

func (primative *Primative) AsString() (string, error) {
	switch val := primative.value.(type) {
	case bool:
		strVal := "False"
		if val {
			strVal = "True"
		}
		return strVal, nil
	case float64:
		return strconv.FormatFloat(val, 10, -1, 64), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case nil:
		return "Null", nil
	case string:
		return val, nil
	}

	msg := fmt.Sprintf("Could not convert value %s to type STRING", primative.value)
	return "", errors.New(msg)
}

func (primative *Primative) String() string {
	str, _ := primative.AsString()
	return str
}
