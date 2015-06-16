package config

import "encoding/json"

type SectionValue struct {
	Name     string
	Value    map[string]ConfigValue
	Comments []string
	Includes []string
}

func NewNamedSection(name string) *SectionValue {
	return &SectionValue{
		Name:     name,
		Value:    make(map[string]ConfigValue),
		Comments: make([]string, 0),
		Includes: make([]string, 0),
	}
}

func NewAnonymousSection() *SectionValue {
	return &SectionValue{
		Value:    make(map[string]ConfigValue),
		Comments: make([]string, 0),
		Includes: make([]string, 0),
	}
}

func (this SectionValue) GetType() ConfigType   { return SECTION }
func (this SectionValue) GetValue() interface{} { return this.Value }

func (this *SectionValue) AddComment(comment string) {
	this.Comments = append(this.Comments, comment)
}

func (this *SectionValue) AddInclude(include string) {
	this.Includes = append(this.Includes, include)
}

func (this *SectionValue) Set(name string, value ConfigValue) {
	this.Value[name] = value
}

func (this *SectionValue) Get(name string) ConfigValue {
	return this.Value[name]
}

func (this *SectionValue) GetSection(name string) SectionValue {
	value := this.Value[name]
	return value.(SectionValue)
}

func (this *SectionValue) GetString(name string) StringValue {
	value := this.Value[name]
	return value.(StringValue)
}

func (this *SectionValue) GetInteger(name string) IntegerValue {
	value := this.Value[name]
	return value.(IntegerValue)
}

func (this *SectionValue) GetFloat(name string) FloatValue {
	value := this.Value[name]
	return value.(FloatValue)
}

func (this *SectionValue) Contains(name string) bool {
	_, ok := this.Value[name]
	return ok
}

func (this *SectionValue) ToJSON() ([]byte, error) {
	data, err := this.ToMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func (this *SectionValue) ToMap() (map[string]interface{}, error) {
	settings := make(map[string]interface{})
	for name, value := range this.Value {
		if value.GetType() == SECTION {
			data, err := value.(*SectionValue).ToMap()
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
