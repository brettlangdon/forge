package token

import "fmt"

type Token struct {
	ID      TokenID
	Literal string
	Line    int
	Column  int
}

func (this Token) String() string {
	return fmt.Sprintf(
		"ID<%s> Literal<%s> Line<%s> Column<%s>",
		this.ID, this.Literal, this.Line, this.Column,
	)
}
