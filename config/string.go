package configTypes

type StringValue struct {
	Name  string
	Value string
}

func (this StringValue) GetType() ConfigType   { return STRING }
func (this StringValue) GetValue() interface{} { return this.Value }
