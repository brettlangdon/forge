package token

type TokenID int

const (
	ILLEGAL TokenID = iota
	EOF

	LBRACKET
	RBRACKET
	EQUAL
	SEMICOLON
	PERIOD

	IDENTIFIER
	BOOLEAN
	INTEGER
	FLOAT
	STRING
	NULL
	COMMENT
)

var tokenNames = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	LBRACKET:   "LBRACKET",
	RBRACKET:   "RBRACKET",
	EQUAL:      "EQUAL",
	SEMICOLON:  "SEMICOLON",
	PERIOD:     "PERIOD",
	IDENTIFIER: "IDENTIFIER",
	BOOLEAN:    "BOOLEAN",
	INTEGER:    "INTEGER",
	FLOAT:      "FLOAT",
	STRING:     "STRING",
	NULL:       "NULL",
	COMMENT:    "COMMENT",
}

func (this TokenID) String() string {
	s := ""
	if 0 <= this && this < TokenID(len(tokenNames)) {
		s = tokenNames[this]
	}

	if s == "" {
		s = "UNKNOWN"
	}

	return s
}
