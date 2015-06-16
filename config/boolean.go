package config

type BooleanValue struct {
	Name  string
	Value bool
}

func (this BooleanValue) GetType() ConfigType   { return BOOLEAN }
func (this BooleanValue) GetValue() interface{} { return this.Value }
