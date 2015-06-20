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

func (scanner *Scanner) readRune() {
	if scanner.newline {
		scanner.curLine += 1
		scanner.curCol = 0
		scanner.newline = false
	} else {
		scanner.curCol += 1
	}

	nextCh, _, err := scanner.reader.ReadRune()
	if err != nil {
		scanner.curCh = eof
		return
	}

	scanner.curCh = nextCh

	if scanner.curCh == '\n' {
		scanner.newline = true
	}
}

func (scanner *Scanner) parseIdentifier() {
	scanner.curTok.ID = token.IDENTIFIER
	scanner.curTok.Literal = string(scanner.curCh)
	for {
		scanner.readRune()
		if !isLetter(scanner.curCh) && scanner.curCh != '_' {
			break
		}
		scanner.curTok.Literal += string(scanner.curCh)
	}

	if isBoolean(scanner.curTok.Literal) {
		scanner.curTok.ID = token.BOOLEAN
	} else if isNull(scanner.curTok.Literal) {
		scanner.curTok.ID = token.NULL
	} else if isInclude(scanner.curTok.Literal) {
		scanner.curTok.ID = token.INCLUDE
	}
}

func (scanner *Scanner) parseNumber() {
	scanner.curTok.ID = token.INTEGER
	scanner.curTok.Literal = string(scanner.curCh)
	digit := false
	for {
		scanner.readRune()
		if scanner.curCh == '.' && digit == false {
			scanner.curTok.ID = token.FLOAT
			digit = true
		} else if !isDigit(scanner.curCh) {
			break
		}
		scanner.curTok.Literal += string(scanner.curCh)
	}
}

func (scanner *Scanner) parseString() {
	scanner.curTok.ID = token.STRING
	scanner.curTok.Literal = string(scanner.curCh)
	for {
		scanner.readRune()
		if scanner.curCh == '"' {
			break
		}
		scanner.curTok.Literal += string(scanner.curCh)
	}
	scanner.readRune()
}

func (scanner *Scanner) parseComment() {
	scanner.curTok.ID = token.COMMENT
	scanner.curTok.Literal = ""
	for {
		scanner.readRune()
		if scanner.curCh == '\n' {
			break
		}
		scanner.curTok.Literal += string(scanner.curCh)
	}
	scanner.readRune()
}

func (scanner *Scanner) skipWhitespace() {
	for {
		scanner.readRune()
		if !isWhitespace(scanner.curCh) {
			break
		}
	}
}

func (scanner *Scanner) NextToken() token.Token {
	if isWhitespace(scanner.curCh) {
		scanner.skipWhitespace()
	}

	scanner.curTok = token.Token{
		ID:      token.ILLEGAL,
		Literal: string(scanner.curCh),
		Line:    scanner.curLine,
		Column:  scanner.curCol,
	}

	switch ch := scanner.curCh; {
	case isLetter(ch) || ch == '_':
		scanner.parseIdentifier()
	case isDigit(ch):
		scanner.parseNumber()
	case ch == '#':
		scanner.parseComment()
	case ch == eof:
		scanner.curTok.ID = token.EOF
		scanner.curTok.Literal = "EOF"
	default:
		scanner.readRune()
		scanner.curTok.Literal = string(ch)
		switch ch {
		case '=':
			scanner.curTok.ID = token.EQUAL
		case '"':
			scanner.parseString()
		case '{':
			scanner.curTok.ID = token.LBRACKET
		case '}':
			scanner.curTok.ID = token.RBRACKET
		case ';':
			scanner.curTok.ID = token.SEMICOLON
		case '.':
			scanner.curTok.ID = token.PERIOD
		}
	}

	return scanner.curTok
}
