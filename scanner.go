package forge

import (
	"bufio"
	"io"
	"strings"

	"github.com/brettlangdon/forge/token"
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

func isInclude(str string) bool {
	return strings.ToLower(str) == "include"
}

type Scanner struct {
	cur_line int
	cur_col  int
	cur_tok  token.Token
	cur_ch   rune
	newline  bool
	reader   *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	scanner := &Scanner{
		reader:   bufio.NewReader(reader),
		cur_line: 0,
		cur_col:  0,
		newline:  false,
	}
	scanner.readRune()
	return scanner
}

func (this *Scanner) readRune() {
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

func (this *Scanner) parseIdentifier() {
	this.cur_tok.ID = token.IDENTIFIER
	this.cur_tok.Literal = string(this.cur_ch)
	for {
		this.readRune()
		if !isLetter(this.cur_ch) && this.cur_ch != '_' {
			break
		}
		this.cur_tok.Literal += string(this.cur_ch)
	}

	if isBoolean(this.cur_tok.Literal) {
		this.cur_tok.ID = token.BOOLEAN
	} else if isNull(this.cur_tok.Literal) {
		this.cur_tok.ID = token.NULL
	} else if isInclude(this.cur_tok.Literal) {
		this.cur_tok.ID = token.INCLUDE
	}
}

func (this *Scanner) parseNumber() {
	this.cur_tok.ID = token.INTEGER
	this.cur_tok.Literal = string(this.cur_ch)
	digit := false
	for {
		this.readRune()
		if this.cur_ch == '.' && digit == false {
			this.cur_tok.ID = token.FLOAT
			digit = true
		} else if !isDigit(this.cur_ch) {
			break
		}
		this.cur_tok.Literal += string(this.cur_ch)
	}
}

func (this *Scanner) parseString() {
	this.cur_tok.ID = token.STRING
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

func (this *Scanner) parseComment() {
	this.cur_tok.ID = token.COMMENT
	this.cur_tok.Literal = ""
	for {
		this.readRune()
		if this.cur_ch == '\n' {
			break
		}
		this.cur_tok.Literal += string(this.cur_ch)
	}
	this.readRune()
}

func (this *Scanner) skipWhitespace() {
	for {
		this.readRune()
		if !isWhitespace(this.cur_ch) {
			break
		}
	}
}

func (this *Scanner) NextToken() token.Token {
	if isWhitespace(this.cur_ch) {
		this.skipWhitespace()
	}

	this.cur_tok = token.Token{
		ID:      token.ILLEGAL,
		Literal: string(this.cur_ch),
		Line:    this.cur_line,
		Column:  this.cur_col,
	}

	switch ch := this.cur_ch; {
	case isLetter(ch) || ch == '_':
		this.parseIdentifier()
	case isDigit(ch):
		this.parseNumber()
	case ch == '#':
		this.parseComment()
	case ch == eof:
		this.cur_tok.ID = token.EOF
		this.cur_tok.Literal = "EOF"
	default:
		this.readRune()
		this.cur_tok.Literal = string(ch)
		switch ch {
		case '=':
			this.cur_tok.ID = token.EQUAL
		case '"':
			this.parseString()
		case '{':
			this.cur_tok.ID = token.LBRACKET
		case '}':
			this.cur_tok.ID = token.RBRACKET
		case ';':
			this.cur_tok.ID = token.SEMICOLON
		case '.':
			this.cur_tok.ID = token.PERIOD
		}
	}

	return this.cur_tok
}
