package config

type ConfigType int

const (
	SECTION ConfigType = iota
	INTEGER
	BOOLEAN
	FLOAT
	STRING
	NULL
)

var configTypes = [...]string{
	SECTION: "SECTION",
	BOOLEAN: "BOOLEAN",
	INTEGER: "INTEGER",
	FLOAT:   "FLOAT",
	STRING:  "STRING",
	NULL:    "NULL",
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
