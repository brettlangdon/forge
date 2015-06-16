package config

type FloatValue struct {
	Name  string
	Value float64
}

func (this FloatValue) GetType() ConfigType   { return INTEGER }
func (this FloatValue) GetValue() interface{} { return this.Value }
