package configTypes

type IntegerValue struct {
	Name  string
	Value int64
}

func (this IntegerValue) GetType() ConfigType   { return INTEGER }
func (this IntegerValue) GetValue() interface{} { return this.Value }
