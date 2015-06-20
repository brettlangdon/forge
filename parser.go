package forge

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/brettlangdon/forge/token"
)

type Parser struct {
	settings   *Section
	scanner    *Scanner
	curTok     token.Token
	curSection *Section
	previous   []*Section
}

func NewParser(reader io.Reader) *Parser {
	settings := NewSection()
	return &Parser{
		scanner:    NewScanner(reader),
		settings:   settings,
		curSection: settings,
		previous:   make([]*Section, 0),
	}
}

func (parser *Parser) syntaxError(msg string) error {
	msg = fmt.Sprintf(
		"Syntax error line <%d> column <%d>: %s",
		parser.curTok.Line,
		parser.curTok.Column,
		msg,
	)
	return errors.New(msg)
}

func (parser *Parser) readToken() token.Token {
	parser.curTok = parser.scanner.NextToken()
	return parser.curTok
}

func (parser *Parser) parseReference(startingSection *Section, period bool) (Value, error) {
	name := ""
	if period == false {
		name = parser.curTok.Literal
	}
	for {
		parser.readToken()
		if parser.curTok.ID == token.PERIOD && period == false {
			period = true
		} else if period && parser.curTok.ID == token.IDENTIFIER {
			if len(name) > 0 {
				name += "."
			}
			name += parser.curTok.Literal
			period = false
		} else if parser.curTok.ID == token.SEMICOLON {
			break
		} else {
			msg := fmt.Sprintf("expected ';' instead found '%s'", parser.curTok.Literal)
			return nil, parser.syntaxError(msg)
		}
	}
	if len(name) == 0 {
		return nil, parser.syntaxError(
			fmt.Sprintf("expected IDENTIFIER instead found %s", parser.curTok.Literal),
		)
	}

	if period {
		return nil, parser.syntaxError(fmt.Sprintf("expected IDENTIFIER after PERIOD"))
	}

	value, err := startingSection.Resolve(name)
	if err != nil {
		err = errors.New("Reference error, " + err.Error())
	}
	return value, nil
}

func (parser *Parser) parseSetting(name string) error {
	var value Value
	parser.readToken()

	readNext := true
	switch parser.curTok.ID {
	case token.STRING:
		value = NewString(parser.curTok.Literal)
	case token.BOOLEAN:
		boolVal, err := strconv.ParseBool(parser.curTok.Literal)
		if err != nil {
			return nil
		}
		value = NewBoolean(boolVal)
	case token.NULL:
		value = NewNull()
	case token.INTEGER:
		intVal, err := strconv.ParseInt(parser.curTok.Literal, 10, 64)
		if err != nil {
			return err
		}
		value = NewInteger(intVal)
	case token.FLOAT:
		floatVal, err := strconv.ParseFloat(parser.curTok.Literal, 64)
		if err != nil {
			return err
		}
		value = NewFloat(floatVal)
	case token.PERIOD:
		reference, err := parser.parseReference(parser.curSection, true)
		if err != nil {
			return err
		}
		value = reference
		readNext = false
	case token.IDENTIFIER:
		reference, err := parser.parseReference(parser.settings, false)
		if err != nil {
			return err
		}
		value = reference
		readNext = false
	default:
		return parser.syntaxError(
			fmt.Sprintf("expected STRING, INTEGER, FLOAT, BOOLEAN or IDENTIFIER, instead found %s", parser.curTok.ID),
		)
	}

	if readNext {
		parser.readToken()
	}
	if parser.curTok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", parser.curTok.Literal)
		return parser.syntaxError(msg)
	}
	parser.readToken()

	parser.curSection.Set(name, value)
	return nil
}

func (parser *Parser) parseInclude() error {
	if parser.curTok.ID != token.STRING {
		msg := fmt.Sprintf("expected STRING instead found '%s'", parser.curTok.ID)
		return parser.syntaxError(msg)
	}
	pattern := parser.curTok.Literal

	parser.readToken()
	if parser.curTok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", parser.curTok.Literal)
		return parser.syntaxError(msg)
	}

	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	oldScanner := parser.scanner
	for _, filename := range filenames {
		reader, err := os.Open(filename)
		if err != nil {
			return err
		}
		// parser.curSection.AddInclude(filename)
		parser.scanner = NewScanner(reader)
		parser.parse()
	}
	parser.scanner = oldScanner
	parser.readToken()
	return nil
}

func (parser *Parser) parseSection(name string) error {
	section := parser.curSection.AddSection(name)
	parser.previous = append(parser.previous, parser.curSection)
	parser.curSection = section
	return nil
}

func (parser *Parser) endSection() error {
	if len(parser.previous) == 0 {
		return parser.syntaxError("unexpected section end '}'")
	}

	pLen := len(parser.previous)
	previous := parser.previous[pLen-1]
	parser.previous = parser.previous[0 : pLen-1]
	parser.curSection = previous
	return nil
}

func (parser *Parser) GetSettings() *Section {
	return parser.settings
}

func (parser *Parser) parse() error {
	parser.readToken()
	for {
		if parser.curTok.ID == token.EOF {
			break
		}
		tok := parser.curTok
		parser.readToken()
		switch tok.ID {
		case token.COMMENT:
			// parser.curSection.AddComment(tok.Literal)
		case token.INCLUDE:
			parser.parseInclude()
		case token.IDENTIFIER:
			if parser.curTok.ID == token.LBRACKET {
				err := parser.parseSection(tok.Literal)
				if err != nil {
					return err
				}
				parser.readToken()
			} else if parser.curTok.ID == token.EQUAL {
				err := parser.parseSetting(tok.Literal)
				if err != nil {
					return err
				}
			}
		case token.RBRACKET:
			err := parser.endSection()
			if err != nil {
				return err
			}
		default:
			return parser.syntaxError(fmt.Sprintf("unexpected token %s", tok))
		}
	}
	return nil
}

func (parser *Parser) Parse() error {
	err := parser.parse()
	if err != nil {
		return err
	}

	if len(parser.previous) > 0 {
		return parser.syntaxError("expected end of section, instead found EOF")
	}

	return nil
}
