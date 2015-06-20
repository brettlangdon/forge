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

func (this *Parser) syntaxError(msg string) error {
	msg = fmt.Sprintf(
		"Syntax error line <%d> column <%d>: %s",
		this.curTok.Line,
		this.curTok.Column,
		msg,
	)
	return errors.New(msg)
}

func (this *Parser) readToken() token.Token {
	this.curTok = this.scanner.NextToken()
	return this.curTok
}

func (this *Parser) parseReference(startingSection *Section, period bool) (Value, error) {
	name := ""
	if period == false {
		name = this.curTok.Literal
	}
	for {
		this.readToken()
		if this.curTok.ID == token.PERIOD && period == false {
			period = true
		} else if period && this.curTok.ID == token.IDENTIFIER {
			if len(name) > 0 {
				name += "."
			}
			name += this.curTok.Literal
			period = false
		} else if this.curTok.ID == token.SEMICOLON {
			break
		} else {
			msg := fmt.Sprintf("expected ';' instead found '%s'", this.curTok.Literal)
			return nil, this.syntaxError(msg)
		}
	}
	if len(name) == 0 {
		return nil, this.syntaxError(
			fmt.Sprintf("expected IDENTIFIER instead found %s", this.curTok.Literal),
		)
	}

	if period {
		return nil, this.syntaxError(fmt.Sprintf("expected IDENTIFIER after PERIOD"))
	}

	value, err := startingSection.Resolve(name)
	if err != nil {
		err = errors.New("Reference error, " + err.Error())
	}
	return value, nil
}

func (this *Parser) parseSetting(name string) error {
	var value Value
	this.readToken()

	read_next := true
	switch this.curTok.ID {
	case token.STRING:
		value = NewString(this.curTok.Literal)
	case token.BOOLEAN:
		boolVal, err := strconv.ParseBool(this.curTok.Literal)
		if err != nil {
			return nil
		}
		value = NewBoolean(boolVal)
	case token.NULL:
		value = NewNull()
	case token.INTEGER:
		intVal, err := strconv.ParseInt(this.curTok.Literal, 10, 64)
		if err != nil {
			return err
		}
		value = NewInteger(intVal)
	case token.FLOAT:
		floatVal, err := strconv.ParseFloat(this.curTok.Literal, 64)
		if err != nil {
			return err
		}
		value = NewFloat(floatVal)
	case token.PERIOD:
		reference, err := this.parseReference(this.curSection, true)
		if err != nil {
			return err
		}
		value = reference
		read_next = false
	case token.IDENTIFIER:
		reference, err := this.parseReference(this.settings, false)
		if err != nil {
			return err
		}
		value = reference
		read_next = false
	default:
		return this.syntaxError(
			fmt.Sprintf("expected STRING, INTEGER, FLOAT, BOOLEAN or IDENTIFIER, instead found %s", this.curTok.ID),
		)
	}

	if read_next {
		this.readToken()
	}
	if this.curTok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", this.curTok.Literal)
		return this.syntaxError(msg)
	}
	this.readToken()

	this.curSection.Set(name, value)
	return nil
}

func (this *Parser) parseInclude() error {
	if this.curTok.ID != token.STRING {
		msg := fmt.Sprintf("expected STRING instead found '%s'", this.curTok.ID)
		return this.syntaxError(msg)
	}
	pattern := this.curTok.Literal

	this.readToken()
	if this.curTok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", this.curTok.Literal)
		return this.syntaxError(msg)
	}

	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	oldScanner := this.scanner
	for _, filename := range filenames {
		reader, err := os.Open(filename)
		if err != nil {
			return err
		}
		// this.curSection.AddInclude(filename)
		this.scanner = NewScanner(reader)
		this.parse()
	}
	this.scanner = oldScanner
	this.readToken()
	return nil
}

func (this *Parser) parseSection(name string) error {
	section := this.curSection.AddSection(name)
	this.previous = append(this.previous, this.curSection)
	this.curSection = section
	return nil
}

func (this *Parser) endSection() error {
	if len(this.previous) == 0 {
		return this.syntaxError("unexpected section end '}'")
	}

	pLen := len(this.previous)
	previous := this.previous[pLen-1]
	this.previous = this.previous[0 : pLen-1]
	this.curSection = previous
	return nil
}

func (this *Parser) GetSettings() *Section {
	return this.settings
}

func (this *Parser) parse() error {
	this.readToken()
	for {
		if this.curTok.ID == token.EOF {
			break
		}
		tok := this.curTok
		this.readToken()
		switch tok.ID {
		case token.COMMENT:
			// this.curSection.AddComment(tok.Literal)
		case token.INCLUDE:
			this.parseInclude()
		case token.IDENTIFIER:
			if this.curTok.ID == token.LBRACKET {
				err := this.parseSection(tok.Literal)
				if err != nil {
					return err
				}
				this.readToken()
			} else if this.curTok.ID == token.EQUAL {
				err := this.parseSetting(tok.Literal)
				if err != nil {
					return err
				}
			}
		case token.RBRACKET:
			err := this.endSection()
			if err != nil {
				return err
			}
		default:
			return this.syntaxError(fmt.Sprintf("unexpected token %s", tok))
		}
	}
	return nil
}

func (this *Parser) Parse() error {
	err := this.parse()
	if err != nil {
		return err
	}

	if len(this.previous) > 0 {
		return this.syntaxError("expected end of section, instead found EOF")
	}

	return nil
}
