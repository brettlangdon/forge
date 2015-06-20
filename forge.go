// Package forge provides an api to deal with parsing configuration files.
//
// Config file example:
//
//     # example.cfg
//     top_level = "a string";
//     primary {
//       primary_int = 500;
//       sub_section {
//         sub_float = 50.5;  # End of line comment
//       }
//     }
//     secondary {
//       secondary_bool = true;
//       secondary_null = null;
//
//       # Reference other config value
//       local_ref = .secondary_null;
//       global_ref = primary.sub_section.sub_float;
//
//       # Include all files matching the provided pattern
//       include "/etc/app/*.cfg";
//     }
//
//
// Config file format:
//
//     IDENTIFIER: [_a-zA-Z]+
//
//     BOOL: 'true' | 'false'
//     NULL: 'null'
//     INTEGER: [0-9]+
//     FLOAT: INTEGER '.' INTEGER
//     STRING: '"' .* '"'
//     REFERENCE: [IDENTIFIER] ('.' IDENTIFIER)+
//     VALUE: BOOL | NULL | INTEGER | FLOAT | STRING | REFERENCE
//
//     INCLUDE: 'include ' STRING ';'
//     DIRECTIVE: (IDENTIFIER '=' VALUE | INCLUDE) ';'
//     SECTION: IDENTIFIER '{' (DIRECTIVE | SECTION)* '}'
//     COMMENT: '#' .* NEWLINE '\n'
//
//     CONFIG_FILE: (COMMENT | DIRECTIVE | SECTION)*
//
//
// Values
//  * String:
//      Any value enclosed in double quotes (single quotes not allowed) (e.g. "string")
//  * Integer:
//      Any number without decimal places (e.g. 500)
//  * Float:
//      Any number with decimal places (e.g. 500.55)
//  * Boolean:
//      The identifiers 'true' or 'false' of any case (e.g. TRUE, True, true, FALSE, False, false)
//  * Null:
//      The identifier 'null' of any case (e.g. NULL, Null, null)
//  * Global reference:
//      An identifier which may contain periods, the references are resolved from the global
//      section (e.g. global_value, section.sub_section.value)
//  * Local reference:
//      An identifier which main contain periods which starts with a period, the references
//      are resolved from the settings current section (e.g. .value, .sub_section.value)
//
// Directives
//  * Comment:
//      A comment is a pound symbol ('#') followed by any text any which ends with a newline (e.g. '# I am a comment\n')
//      A comment can either be on a line of it's own or at the end of any line. Nothing can come after the comment
//      until after the newline.
//  * Directive:
//      A directive is a setting, a identifier and a value. They are in the format '<identifier> = <value>;'
//      All directives must end in a semicolon. The value can be any of the types defined above.
//  * Section:
//      A section is a grouping of directives under a common name. They are in the format '<section_name> { <directives> }'.
//      All sections must be wrapped in brackets ('{', '}') and must all have a name. They do not end in a semicolon.
//      Sections may be left empty, they do not have to contain any directives.
//  * Include:
//      An include statement tells the config parser to include the contents of another config file where the include
//      statement is defined. Includes are in the format 'include "<pattern>";'. The <pattern> can be any glob
//      like pattern which is compatible with `path.filepath.Match` http://golang.org/pkg/path/filepath/#Match
//
// Note about references.
//
// References are resolved during parsing not access (they do not act as pointers... today).
// Meaning that if I have the following config:
//   setting = "value";
//   ref_to_setting = setting;
//
// During parsing when I get to `ref_to_setting` it will be resolved immediately and it's value will
// be set to `"value"`.
//
// This config example, however, will raise a parsing error since it cannot resolve the reference:
//   ref_to_setting = setting;
//   setting = "value";
//
// Since `setting` does not exist yet.
// This is mostly due to the naive implementation of references today. How they work may change in
// future versions or they may disappear entirely.
//
package forge

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func ParseBytes(data []byte) (*Section, error) {
	return ParseReader(bytes.NewReader(data))
}

func ParseFile(filename string) (*Section, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return ParseReader(reader)
}

func ParseReader(reader io.Reader) (*Section, error) {
	parser := NewParser(reader)
	err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return parser.GetSettings(), nil
}

func ParseString(data string) (*Section, error) {
	return ParseReader(strings.NewReader(data))
}
