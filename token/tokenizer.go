package token

import (
	"bufio"
	"io"
	"strings"
)

var eof = rune(0)

func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9')
}

func isWhitespace(ch rune) bool {
	return (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r')
}

func isBoolean(str string) bool {
	lower := strings.ToLower(str)
	return lower == "true" || lower == "false"
}

func isNull(str string) bool {
	return strings.ToLower(str) == "null"
}

type Tokenizer struct {
	cur_line int
	cur_col  int
	cur_tok  Token
	cur_ch   rune
	newline  bool
	reader   *bufio.Reader
}

func NewTokenizer(reader io.Reader) *Tokenizer {
	tokenizer := &Tokenizer{
		reader:   bufio.NewReader(reader),
		cur_line: 0,
		cur_col:  0,
		newline:  false,
	}
	tokenizer.readRune()
	return tokenizer
}

func (this *Tokenizer) readRune() {
	if this.newline {
		this.cur_line += 1
		this.cur_col = 0
		this.newline = false
	} else {
		this.cur_col += 1
	}

	next_ch, _, err := this.reader.ReadRune()
	if err != nil {
		this.cur_ch = eof
		return
	}

	this.cur_ch = next_ch

	if this.cur_ch == '\n' {
		this.newline = true
	}
}

func (this *Tokenizer) parseIdentifier() {
	this.cur_tok.ID = IDENTIFIER
	this.cur_tok.Literal = string(this.cur_ch)
	for {
		this.readRune()
		if !isLetter(this.cur_ch) && this.cur_ch != '_' {
			break
		}
		this.cur_tok.Literal += string(this.cur_ch)
	}

	if isBoolean(this.cur_tok.Literal) {
		this.cur_tok.ID = BOOLEAN
	} else if isNull(this.cur_tok.Literal) {
		this.cur_tok.ID = NULL
	}
}

func (this *Tokenizer) parseNumber() {
	this.cur_tok.ID = INTEGER
	this.cur_tok.Literal = string(this.cur_ch)
	digit := false
	for {
		this.readRune()
		if this.cur_ch == '.' && digit == false {
			this.cur_tok.ID = FLOAT
			digit = true
		} else if !isDigit(this.cur_ch) {
			break
		}
		this.cur_tok.Literal += string(this.cur_ch)
	}
}

func (this *Tokenizer) parseString() {
	this.cur_tok.ID = STRING
	this.cur_tok.Literal = string(this.cur_ch)
	for {
		this.readRune()
		if this.cur_ch == '"' {
			break
		}
		this.cur_tok.Literal += string(this.cur_ch)
	}
	this.readRune()
}

func (this *Tokenizer) skipWhitespace() {
	for {
		this.readRune()
		if !isWhitespace(this.cur_ch) {
			break
		}
	}
}

func (this *Tokenizer) NextToken() Token {
	if isWhitespace(this.cur_ch) {
		this.skipWhitespace()
	}

	this.cur_tok = Token{
		ID:      ILLEGAL,
		Literal: string(this.cur_ch),
		Line:    this.cur_line,
		Column:  this.cur_col,
	}

	switch ch := this.cur_ch; {
	case isLetter(ch) || ch == '_':
		this.parseIdentifier()
	case isDigit(ch):
		this.parseNumber()
	case ch == eof:
		this.cur_tok.ID = EOF
		this.cur_tok.Literal = "EOF"
	default:
		this.readRune()
		this.cur_tok.Literal = string(ch)
		switch ch {
		case '=':
			this.cur_tok.ID = EQUAL
		case '"':
			this.parseString()
		case '{':
			this.cur_tok.ID = LBRACKET
		case '}':
			this.cur_tok.ID = RBRACKET
		case ';':
			this.cur_tok.ID = SEMICOLON
		case '.':
			this.cur_tok.ID = PERIOD
		}
	}

	return this.cur_tok
}
