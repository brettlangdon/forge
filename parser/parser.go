package parser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/brettlangdon/forge/config"
	"github.com/brettlangdon/forge/token"
)

type Parser struct {
	settings    config.SectionValue
	tokenizer   *token.Tokenizer
	cur_tok     token.Token
	cur_section config.SectionValue
	previous    []config.SectionValue
}

func (this *Parser) SyntaxError(msg string) error {
	msg = fmt.Sprintf(
		"Syntax error line <%d> column <%d>: %s",
		this.cur_tok.Line,
		this.cur_tok.Column,
		msg,
	)
	return errors.New(msg)
}

func (this *Parser) ReferenceTypeError(names []string, expected config.ConfigType, actual config.ConfigType) error {
	reference := strings.Join(names, ".")
	msg := fmt.Sprintf(
		"Reference type error, '%s', expected type %s instead got %s",
		reference,
		expected,
		actual,
	)
	return errors.New(msg)
}

func (this *Parser) ReferenceMissingError(names []string, searching string) error {
	reference := strings.Join(names, ".")
	msg := fmt.Sprintf(
		"Reference missing error, '%s' does not have key '%s'",
		reference,
		searching,
	)
	return errors.New(msg)
}

func (this *Parser) readToken() token.Token {
	this.cur_tok = this.tokenizer.NextToken()
	return this.cur_tok
}

func (this *Parser) parseReference(starting_section config.SectionValue, period bool) (config.ConfigValue, error) {
	names := []string{}
	if period == false {
		names = append(names, this.cur_tok.Literal)
	}
	for {
		this.readToken()
		if this.cur_tok.ID == token.PERIOD && period == false {
			period = true
		} else if period && this.cur_tok.ID == token.IDENTIFIER {
			names = append(names, this.cur_tok.Literal)
			period = false
		} else if this.cur_tok.ID == token.SEMICOLON {
			break
		} else {
			msg := fmt.Sprintf("expected ';' instead found '%s'", this.cur_tok.Literal)
			return nil, this.SyntaxError(msg)
		}
	}
	if len(names) == 0 {
		return nil, this.SyntaxError(
			fmt.Sprintf("expected IDENTIFIER instead found %s", this.cur_tok.Literal),
		)
	}

	if period {
		return nil, this.SyntaxError(fmt.Sprintf("expected IDENTIFIER after PERIOD"))
	}

	var reference config.ConfigValue
	reference = starting_section
	visited := []string{}
	for {
		if len(names) == 0 {
			break
		}
		if reference.GetType() != config.SECTION {
			return nil, this.ReferenceTypeError(visited, config.SECTION, reference.GetType())
		}
		name := names[0]
		names = names[1:]
		section := reference.(config.SectionValue)
		if section.Contains(name) == false {
			return nil, this.ReferenceMissingError(visited, name)
		}
		reference = section.Get(name)
		visited = append(visited, name)
	}

	return reference, nil
}

func (this *Parser) parseSetting(name string) error {
	var value config.ConfigValue
	this.readToken()

	read_next := true
	switch this.cur_tok.ID {
	case token.STRING:
		value = config.StringValue{
			Name:  name,
			Value: this.cur_tok.Literal,
		}
	case token.BOOLEAN:
		bool_val, err := strconv.ParseBool(this.cur_tok.Literal)
		if err != nil {
			return nil
		}
		value = config.BooleanValue{
			Name:  name,
			Value: bool_val,
		}
	case token.NULL:
		value = config.NullValue{
			Name:  name,
			Value: nil,
		}
	case token.INTEGER:
		int_val, err := strconv.ParseInt(this.cur_tok.Literal, 10, 64)
		if err != nil {
			return err
		}
		value = config.IntegerValue{
			Name:  name,
			Value: int_val,
		}
	case token.FLOAT:
		float_val, err := strconv.ParseFloat(this.cur_tok.Literal, 64)
		if err != nil {
			return err
		}
		value = config.FloatValue{
			Name:  name,
			Value: float_val,
		}
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
		return this.SyntaxError(
			fmt.Sprintf("expected STRING, INTEGER or FLOAT, instead found %s", this.cur_tok.ID),
		)
	}

	if read_next {
		this.readToken()
	}
	if this.cur_tok.ID != token.SEMICOLON {
		msg := fmt.Sprintf("expected ';' instead found '%s'", this.cur_tok.Literal)
		return this.SyntaxError(msg)
	}

	this.cur_section.Set(name, value)
	return nil
}

func (this *Parser) parseSection(name string) error {
	section := config.SectionValue{
		Name:  name,
		Value: make(map[string]config.ConfigValue),
	}
	this.cur_section.Set(name, section)
	this.previous = append(this.previous, this.cur_section)
	this.cur_section = section
	return nil
}

func (this *Parser) endSection() error {
	if len(this.previous) == 0 {
		return this.SyntaxError("unexpected section end '}'")
	}

	p_len := len(this.previous)
	previous := this.previous[p_len-1]
	this.previous = this.previous[0 : p_len-1]
	this.cur_section = previous
	return nil
}

func (this *Parser) Parse() error {
	this.readToken()
	for {
		if this.cur_tok.ID == token.EOF {
			break
		}
		tok := this.cur_tok
		this.readToken()
		switch tok.ID {
		case token.IDENTIFIER:
			if this.cur_tok.ID == token.LBRACKET {
				err := this.parseSection(tok.Literal)
				if err != nil {
					return err
				}
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
		}
	}
	return nil
}

func ParseFile(filename string) (settings *config.SectionValue, err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return settings, err
	}
	return ParseReader(reader)
}

func ParseReader(reader io.Reader) (*config.SectionValue, error) {
	settings := config.SectionValue{
		Value: make(map[string]config.ConfigValue),
	}
	parser := &Parser{
		tokenizer:   token.NewTokenizer(reader),
		settings:    settings,
		cur_section: settings,
		previous:    make([]config.SectionValue, 0),
	}
	err := parser.Parse()
	if err != nil {
		return nil, err
	}

	if len(parser.previous) > 0 {
		return nil, parser.SyntaxError("expected end of section, instead found EOF")
	}

	return &settings, nil
}
