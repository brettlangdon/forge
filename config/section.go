package config

import "encoding/json"

type SectionValue struct {
	Name  string
	Value map[string]ConfigValue
}

func (this SectionValue) GetType() ConfigType   { return SECTION }
func (this SectionValue) GetValue() interface{} { return this.Value }

func (this SectionValue) Set(name string, value ConfigValue) {
	this.Value[name] = value
}

func (this SectionValue) Get(name string) ConfigValue {
	return this.Value[name]
}

func (this SectionValue) GetSection(name string) SectionValue {
	value := this.Value[name]
	return value.(SectionValue)
}

func (this SectionValue) GetString(name string) StringValue {
	value := this.Value[name]
	return value.(StringValue)
}

func (this SectionValue) GetInteger(name string) IntegerValue {
	value := this.Value[name]
	return value.(IntegerValue)
}

func (this SectionValue) GetFloat(name string) FloatValue {
	value := this.Value[name]
	return value.(FloatValue)
}

func (this SectionValue) Contains(name string) bool {
	_, ok := this.Value[name]
	return ok
}

func (this SectionValue) ToJSON() ([]byte, error) {
	data, err := this.ToMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func (this SectionValue) ToMap() (map[string]interface{}, error) {
	settings := make(map[string]interface{})
	for name, value := range this.Value {
		if value.GetType() == SECTION {
			data, err := value.(SectionValue).ToMap()
			if err != nil {
				return nil, err
			}
			settings[name] = data
		} else {
			settings[name] = value.GetValue()
		}
	}

	return settings, nil
}
