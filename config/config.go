package config

type ConfigType int

const (
	SECTION ConfigType = iota
	INTEGER
	BOOLEAN
	FLOAT
	STRING
)

var configTypes = [...]string{
	SECTION: "SECTION",
	BOOLEAN: "BOOLEAN",
	INTEGER: "INTEGER",
	FLOAT:   "FLOAT",
	STRING:  "STRING",
}

func (this ConfigType) String() string {
	s := ""
	if 0 <= this && this < ConfigType(len(configTypes)) {
		s = configTypes[this]
	}

	if s == "" {
		s = "UNKNOWN"
	}

	return s
}

type ConfigValue interface {
	GetType() ConfigType
	GetValue() interface{}
}

type BooleanValue struct {
	Name  string
	Value bool
}

func (this BooleanValue) GetType() ConfigType   { return BOOLEAN }
func (this BooleanValue) GetValue() interface{} { return this.Value }

type IntegerValue struct {
	Name  string
	Value int64
}

func (this IntegerValue) GetType() ConfigType   { return INTEGER }
func (this IntegerValue) GetValue() interface{} { return this.Value }

type FloatValue struct {
	Name  string
	Value float64
}

func (this FloatValue) GetType() ConfigType   { return INTEGER }
func (this FloatValue) GetValue() interface{} { return this.Value }

type StringValue struct {
	Name  string
	Value string
}

func (this StringValue) GetType() ConfigType   { return STRING }
func (this StringValue) GetValue() interface{} { return this.Value }
