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

func (this *Section) GetType() ValueType {
	return SECTION
}

func (this *Section) GetValue() interface{} {
	return this.values
}

func (this *Section) UpdateValue(value interface{}) error {
	switch value.(type) {
	case map[string]Value:
		this.values = value.(map[string]Value)
		return nil
	}

	msg := fmt.Sprintf("Unsupported type, %s must be of type `map[string]Value`", value)
	return errors.New(msg)
}

func (this *Section) AddSection(name string) *Section {
	section := NewChildSection(this)
	this.values[name] = section
	return section
}

func (this *Section) Exists(name string) bool {
	_, err := this.Get(name)
	return err == nil
}

func (this *Section) Get(name string) (Value, error) {
	value, ok := this.values[name]
	var err error
	if ok == false {
		err = errors.New("Value does not exist")
	}
	return value, err
}

func (this *Section) GetBoolean(name string) (bool, error) {
	value, err := this.Get(name)
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

func (this *Section) GetFloat(name string) (float64, error) {
	value, err := this.Get(name)
	if err != nil {
		return float64(0), err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsFloat()
	}

	return float64(0), errors.New("Could not convert non-primative value to float")
}

func (this *Section) GetInteger(name string) (int64, error) {
	value, err := this.Get(name)
	if err != nil {
		return int64(0), err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsInteger()
	}

	return int64(0), errors.New("Could not convert non-primative value to integer")
}

func (this *Section) GetSection(name string) (*Section, error) {
	value, err := this.Get(name)
	if err != nil {
		return nil, err
	}

	if value.GetType() == SECTION {
		return value.(*Section), nil
	}
	return nil, errors.New("Could not fetch value as section")
}

func (this *Section) GetString(name string) (string, error) {
	value, err := this.Get(name)
	if err != nil {
		return "", err
	}

	switch value.(type) {
	case *Primative:
		return value.(*Primative).AsString()
	}

	return "", errors.New("Could not convert non-primative value to string")
}

func (this *Section) GetParent() *Section {
	return this.parent
}

func (this *Section) HasParent() bool {
	return this.parent != nil
}

func (this *Section) Set(name string, value Value) {
	this.values[name] = value
}

func (this *Section) SetBoolean(name string, value bool) {
	current, err := this.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		this.values[name] = NewBoolean(value)
	}
}

func (this *Section) SetFloat(name string, value float64) {
	current, err := this.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		this.values[name] = NewFloat(value)
	}
}

func (this *Section) SetInteger(name string, value int64) {
	current, err := this.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		this.values[name] = NewInteger(value)
	}
}

func (this *Section) SetNull(name string) {
	current, err := this.Get(name)

	// Already is a Null, nothing to do
	if err == nil && current.GetType() == NULL {
		return
	}
	this.Set(name, NewNull())
}

func (this *Section) SetString(name string, value string) {
	current, err := this.Get(name)

	// Exists just update the value/type
	if err == nil {
		current.UpdateValue(value)
	} else {
		this.Set(name, NewString(value))
	}
}

func (this *Section) Resolve(name string) (Value, error) {
	// Used only in error state return value
	var value Value

	parts := strings.Split(name, ".")
	if len(parts) == 0 {
		return value, errors.New("No name provided")
	}

	var current Value
	current = this
	for _, part := range parts {
		if current.GetType() != SECTION {
			return value, errors.New("Trying to resolve value from non-section")
		}

		next_current, err := current.(*Section).Get(part)
		if err != nil {
			return value, errors.New("Could not find value in section")
		}
		current = next_current
	}
	return current, nil
}

func (this *Section) ToJSON() ([]byte, error) {
	data := this.ToMap()
	return json.Marshal(data)
}

func (this *Section) ToMap() map[string]interface{} {
	output := make(map[string]interface{})

	for key, value := range this.values {
		if value.GetType() == SECTION {
			output[key] = value.(*Section).ToMap()
		} else {
			output[key] = value.GetValue()
		}
	}
	return output
}
