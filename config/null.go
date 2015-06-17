package config

type NullValue struct {
	Name  string
	Value interface{}
}

func (this NullValue) GetType() ConfigType   { return NULL }
func (this NullValue) GetValue() interface{} { return nil }
