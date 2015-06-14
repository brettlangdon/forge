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
	INTEGER
	FLOAT
	STRING
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
	INTEGER:    "INTEGER",
	FLOAT:      "FLOAT",
	STRING:     "STRING",
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
