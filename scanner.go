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
	curLine int
	curCol  int
	curTok  token.Token
	curCh   rune
	newline bool
	reader  *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	scanner := &Scanner{
		reader:  bufio.NewReader(reader),
		curLine: 0,
		curCol:  0,
		newline: false,
	}
	scanner.readRune()
	return scanner
}

func (this *Scanner) readRune() {
	if this.newline {
		this.curLine += 1
		this.curCol = 0
		this.newline = false
	} else {
		this.curCol += 1
	}

	nextCh, _, err := this.reader.ReadRune()
	if err != nil {
		this.curCh = eof
		return
	}

	this.curCh = nextCh

	if this.curCh == '\n' {
		this.newline = true
	}
}

func (this *Scanner) parseIdentifier() {
	this.curTok.ID = token.IDENTIFIER
	this.curTok.Literal = string(this.curCh)
	for {
		this.readRune()
		if !isLetter(this.curCh) && this.curCh != '_' {
			break
		}
		this.curTok.Literal += string(this.curCh)
	}

	if isBoolean(this.curTok.Literal) {
		this.curTok.ID = token.BOOLEAN
	} else if isNull(this.curTok.Literal) {
		this.curTok.ID = token.NULL
	} else if isInclude(this.curTok.Literal) {
		this.curTok.ID = token.INCLUDE
	}
}

func (this *Scanner) parseNumber() {
	this.curTok.ID = token.INTEGER
	this.curTok.Literal = string(this.curCh)
	digit := false
	for {
		this.readRune()
		if this.curCh == '.' && digit == false {
			this.curTok.ID = token.FLOAT
			digit = true
		} else if !isDigit(this.curCh) {
			break
		}
		this.curTok.Literal += string(this.curCh)
	}
}

func (this *Scanner) parseString() {
	this.curTok.ID = token.STRING
	this.curTok.Literal = string(this.curCh)
	for {
		this.readRune()
		if this.curCh == '"' {
			break
		}
		this.curTok.Literal += string(this.curCh)
	}
	this.readRune()
}

func (this *Scanner) parseComment() {
	this.curTok.ID = token.COMMENT
	this.curTok.Literal = ""
	for {
		this.readRune()
		if this.curCh == '\n' {
			break
		}
		this.curTok.Literal += string(this.curCh)
	}
	this.readRune()
}

func (this *Scanner) skipWhitespace() {
	for {
		this.readRune()
		if !isWhitespace(this.curCh) {
			break
		}
	}
}

func (this *Scanner) NextToken() token.Token {
	if isWhitespace(this.curCh) {
		this.skipWhitespace()
	}

	this.curTok = token.Token{
		ID:      token.ILLEGAL,
		Literal: string(this.curCh),
		Line:    this.curLine,
		Column:  this.curCol,
	}

	switch ch := this.curCh; {
	case isLetter(ch) || ch == '_':
		this.parseIdentifier()
	case isDigit(ch):
		this.parseNumber()
	case ch == '#':
		this.parseComment()
	case ch == eof:
		this.curTok.ID = token.EOF
		this.curTok.Literal = "EOF"
	default:
		this.readRune()
		this.curTok.Literal = string(ch)
		switch ch {
		case '=':
			this.curTok.ID = token.EQUAL
		case '"':
			this.parseString()
		case '{':
			this.curTok.ID = token.LBRACKET
		case '}':
			this.curTok.ID = token.RBRACKET
		case ';':
			this.curTok.ID = token.SEMICOLON
		case '.':
			this.curTok.ID = token.PERIOD
		}
	}

	return this.curTok
}
