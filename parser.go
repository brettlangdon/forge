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
	settings    *Section
	scanner     *Scanner
	cur_tok     token.Token
	cur_section *Section
	previous    []*Section
}

func NewParser(reader io.Reader) *Parser {
	settings := NewSection()
	return &Parser{
		scanner:     NewScanner(reader),
		settings:    settings,
		cur_section: settings,
		previous:    make([]*Section, 0),
	}
}

func (this *Parser) syntaxError(msg string) error {
	msg = fmt.Sprintf(
		"Syntax error line <%d> column <%d>: %s",
		this.cur_tok.Line,
		this.cur_tok.Column,
		msg,
	)
	return errors.New(msg)
}

func (this *Parser) readToken() token.Token {
	this.cur_tok = this.scanner.NextToken()
	return this.cur_tok
}

func (this *Parser) parseReference(starting_section *Section, period bool) (Value, error) {
	name := ""
	if period == false {
		name = this.cur_tok.Literal
	}
	for {
		this.readToken()
		if this.cur_tok.ID == token.PERIOD && period == false {
			period = true
		} else if period && this.cur_tok.ID == token.IDENTIFIER {
			if len(name) > 0 {
				name += "."
			}
			name += this.cur_tok.Literal
			period = false
		} else if this.cur_tok.ID == token.SEMICOLON {
			break
		} else {
			msg := fmt.Sprintf("expected ';' instead found '%s'", this.cur_tok.Literal)
			return nil, this.syntaxError(msg)
		}
	}
	if len(name) == 0 {
		return nil, this.syntaxError(
			fmt.Sprintf("expected IDENTIFIER instead found %s", this.cur_tok.Literal),
		)
	}

	if period {
		return nil, this.syntaxError(fmt.Sprintf("expected IDENTIFIER after PERIOD"))
	}

	value, err := starting_section.Resolve(name)
	if err != nil {
		err = errors.New("Reference error, " + err.Error())
	}
	return value, nil
}

func (this *Parser) parseSetting(name string) error {
	var value Value
	this.readToken()

	read_next := true
	switch this.cur_tok.ID {
	case token.STRING:
		value = NewString(this.cur_tok.Literal)
	case token.BOOLEAN:
		bool_val, err := strconv.ParseBool(this.cur_tok.Literal)
		if err != nil {
			return nil
		}
		value = NewBoolean(bool_val)
	case token.NULL:
		value = NewNull()
	case token.INTEGER:
		int_val, err := strconv.ParseInt(this.cur_tok.Literal, 10, 64)
		if err != nil {
			return err
		}
		value = NewInteger(int_val)
	case token.FLOAT:
		float_val, err := strconv.ParseFloat(this.cur_tok.Literal, 64)
		if err != nil {
			return err
		}
		value = NewFloat(float_val)
	case token.PERIOD:
		reference, err := this.parseReference(this.cur_section, true)
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
			fmt.Sprintf("expected STRING, INTEGER, FLOAT, BOOLEAN or IDENTIFIER, instead found %s", this.cur_tok.ID),
		)
	}

	if read_next {
		this.readToken()
	}
	if this.cur_tok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", this.cur_tok.Literal)
		return this.syntaxError(msg)
	}
	this.readToken()

	this.cur_section.Set(name, value)
	return nil
}

func (this *Parser) parseInclude() error {
	if this.cur_tok.ID != token.STRING {
		msg := fmt.Sprintf("expected STRING instead found '%s'", this.cur_tok.ID)
		return this.syntaxError(msg)
	}
	pattern := this.cur_tok.Literal

	this.readToken()
	if this.cur_tok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", this.cur_tok.Literal)
		return this.syntaxError(msg)
	}

	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	old_scanner := this.scanner
	for _, filename := range filenames {
		reader, err := os.Open(filename)
		if err != nil {
			return err
		}
		// this.cur_section.AddInclude(filename)
		this.scanner = NewScanner(reader)
		this.parse()
	}
	this.scanner = old_scanner
	this.readToken()
	return nil
}

func (this *Parser) parseSection(name string) error {
	section := this.cur_section.AddSection(name)
	this.previous = append(this.previous, this.cur_section)
	this.cur_section = section
	return nil
}

func (this *Parser) endSection() error {
	if len(this.previous) == 0 {
		return this.syntaxError("unexpected section end '}'")
	}

	p_len := len(this.previous)
	previous := this.previous[p_len-1]
	this.previous = this.previous[0 : p_len-1]
	this.cur_section = previous
	return nil
}

func (this *Parser) GetSettings() *Section {
	return this.settings
}

func (this *Parser) parse() error {
	this.readToken()
	for {
		if this.cur_tok.ID == token.EOF {
			break
		}
		tok := this.cur_tok
		this.readToken()
		switch tok.ID {
		case token.COMMENT:
			// this.cur_section.AddComment(tok.Literal)
		case token.INCLUDE:
			this.parseInclude()
		case token.IDENTIFIER:
			if this.cur_tok.ID == token.LBRACKET {
				err := this.parseSection(tok.Literal)
				if err != nil {
					return err
				}
				this.readToken()
			} else if this.cur_tok.ID == token.EQUAL {
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
