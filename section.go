package forge

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Section struct {
	parent *Section
	values map[string]Value
}

func NewSection() *Section {
	return &Section{
		values: make(map[string]Value),
	}
}

func NewChildSection(parent *Section) *Section {
	return &Section{
		parent: parent,
		values: make(map[string]Value),
	}
}

func (section *Section) GetType() ValueType {
	return SECTION
}

func (section *Section) GetValue() interface{} {
	return section.values
}

func (section *Section) UpdateValue(value interface{}) error {
	switch value.(type) {
	case map[string]Value:
		section.values = value.(map[string]Value)
		return nil
	}

	msg := fmt.Sprintf("Unsupported type, %s must be of type `map[string]Value`", value)
	return errors.New(msg)
}

func (section *Section) AddSection(name string) *Section {
	childSection := NewChildSection(section)
	section.values[name] = childSection
	return childSection
}

func (section *Section) Exists(name string) bool {
	_, err := section.Get(name)
	return err == nil
}

func (section *Section) Get(name string) (Value, error) {
	value, ok := section.values[name]
	var err error
	if ok == false {
		err = errors.New("Value does not exist")
	}
	return value, err
}

func (section *Section) GetBoolean(name string) (bool, error) {
	value, err := section.Get(name)
	if err != nil {
		return false, err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsBoolean()
	case *Section:
		return true, nil
	}

	return false, errors.New("Could not convert unknown value to boolean")
}

func (section *Section) GetFloat(name string) (float64, error) {
	value, err := section.Get(name)
	if err != nil {
		return float64(0), err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsFloat()
	}

	return float64(0), errors.New("Could not convert non-primative value to float")
}

func (section *Section) GetInteger(name string) (int64, error) {
	value, err := section.Get(name)
	if err != nil {
		return int64(0), err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsInteger()
	}

	return int64(0), errors.New("Could not convert non-primative value to integer")
}

func (section *Section) GetSection(name string) (*Section, error) {
	value, err := section.Get(name)
	if err != nil {
		return nil, err
	}

	if value.GetType() == SECTION {
		return value.(*Section), nil
	}
	return nil, errors.New("Could not fetch value as section")
}

func (section *Section) GetString(name string) (string, error) {
	value, err := section.Get(name)
	if err != nil {
		return "", err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsString()
	}

	return "", errors.New("Could not convert non-primative value to string")
}

func (section *Section) GetParent() *Section {
	return section.parent
}

func (section *Section) HasParent() bool {
	return section.parent != nil
}

func (section *Section) Set(name string, value Value) {
	section.values[name] = value
}

func (section *Section) SetBoolean(name string, value bool) {
	current, err := section.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		section.values[name] = NewBoolean(value)
	}
}

func (section *Section) SetFloat(name string, value float64) {
	current, err := section.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		section.values[name] = NewFloat(value)
	}
}

func (section *Section) SetInteger(name string, value int64) {
	current, err := section.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		section.values[name] = NewInteger(value)
	}
}

func (section *Section) SetNull(name string) {
	current, err := section.Get(name)

	// Already is a Null, nothing to do
	if err == nil && current.GetType() == NULL {
		return
	}
	section.Set(name, NewNull())
}

func (section *Section) SetString(name string, value string) {
	current, err := section.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		section.Set(name, NewString(value))
	}
}

func (section *Section) Resolve(name string) (Value, error) {
	// Used only in error state return value
	var value Value

	parts := strings.Split(name, ".")
	if len(parts) == 0 {
		return value, errors.New("No name provided")
	}

	var current Value
	current = section
	for _, part := range parts {
		if current.GetType() != SECTION {
			return value, errors.New("Trying to resolve value from non-section")
		}

		nextCurrent, err := current.(*Section).Get(part)
		if err != nil {
			return value, errors.New("Could not find value in section")
		}
		current = nextCurrent
	}
	return current, nil
}

func (section *Section) ToJSON() ([]byte, error) {
	data := section.ToMap()
	return json.Marshal(data)
}

func (section *Section) ToMap() map[string]interface{} {
	output := make(map[string]interface{})

	for key, value := range section.values {
		if value.GetType() == SECTION {
			output[key] = value.(*Section).ToMap()
		} else {
			output[key] = value.GetValue()
		}
	}
	return output
}
