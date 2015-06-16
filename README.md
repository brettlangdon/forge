forge
=====

Forge is a configuration syntax and parser.

## Installation

`git get github.com/brettlangdon/forge`

## File format

The format was influenced a lot by nginx configuration file format.

```config
# Global settings
global_key = "string value";

# Sub section
sub_settings {
  sub_int = 500;
  sub_float = 80.80;
  # Sub-Sub Section
  sub_sub_settings {
    sub_sub_sub_settings {
      key = "value";
    }
  }
}

# Second section
second {
  key = "value";
  global_reference = sub_settings.sub_float;
  local_reference = .key;  # References second.key
}
```

For normal settings the format is the key followed by an equal sign followed by the value and lastly ending with a semicolon.
`<key> = <value>;`

Sections (basically a map) is formatted as the section name with the section's settings wrapped in brackets.
`<section> { <key> = <value>; }`

Comments start with a pound sign `#` and end with a newline. A comment can exist on the same line as settings/sections, but the comment must end the line.

## Data types

### Boolean
A boolean value is either `true` or `false` of any case.

`TRUE`, `true`, `True`, `FALSE`, `False`, `false`.

### Null

A null value is allowed as `null` of any case.

`NULL`, `Null`, `null`.

### String
A string value is wrapped by double quotes (single quotes will not work).

`"string value"`, `"single ' quotes ' allowed"`.

As of right now there is no way to escape double quotes within a string's value;

### Number

There are two supported numbers, Integer and Float, both of which are simply numbers with the later having one period.

`500`, `50.56`.

### Section

Sections are essentially maps, that is a setting whose purpose is to hold other settings.
Sections can be used to namespace settings.

`section { setting = "value"; }`.


### References

References are used to refer to previously defined settings. There are two kinds of references, a global reference and a local reference;

The general format for a reference is a mix of identifiers and periods, for example `production.db.name`.

A global reference is a reference which starts looking for its value from the top most section (global section).

A local reference is a reference whose value starts with a period, this reference will start looking for it's value from the current section it is within (local section).

```config
production {
  db {
    name = "forge";
  }
}

development {
  db {
    name = production.db.name;
  }
  db_name = .db.name;
}
```

## API

`github.com/brettlangdon/forge`

* `forge.ParseString(data string) (map[string]interface{}, error)`
* `forge.ParseBytes(data []byte) (map[string]interface{}, error)`
* `forge.ParseFile(filename string) (map[string]interface{}, error)`
* `forge.ParseReader(reader io.Reader) (map[string]interface{}, error)`


## Example

You can see example usage in the `example` folder.

```go
package main

import (
	"fmt"
	"json"

	"github.com/brettlangdon/forge"
)

func main() {
	// Parse the file `example.cfg` as a map[string]interface{}
	settings, err := forge.ParseFile("example.cfg")
	if err != nil {
		panic(err)
	}

	// Convert the settings to JSON for printing
	jsonBytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	// Print the parsed settings
	fmt.Println(string(jsonBytes))
}
```

## Future Plans

The following features are currently on my bucket list for the future:

### Operations/Expressions

Would be nice to have Addition/Subtraction/Multiplication/Division:

```config
whole = 100
half = whole / 2;
double = whole * 2;
one_more = whole + 1;
one_less = whole - 1;
```

Also Concatenation for strings:

```config
domain = "github.com";
username = "brettlangdon";
name = "forge";
repo_url = domain + "/" + username + "/" + name;
```

### API

I'll probably revisit the API, I just threw it together quick, want to make sure it right.

### Documentation

Documentation is a good thing.
