package forge_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/brettlangdon/forge"
)

var testConfigBytes = []byte(`
# Global stuff
global = "global value";
# Primary stuff
primary {
  string = "primary string value";
  string_with_quote = "some \"quoted\" str\\ing";
  single = 'hello world';
  single_with_quote = '\'hello\' "world"';
  integer = 500;
  float = 80.80;
  negative = -50;
  boolean = true;
  not_true = FALSE;
  nothing = NULL;
  # Reference secondary._under (which hasn't been defined yet)
  sec_ref = secondary._under;
   # Primary-sub stuff
  sub {
      key = "primary sub key value";
      include "./test_include.cfg";
  }
}

secondary {
  another = "secondary another value";
  global_reference = global;
  primary_sub_key = primary.sub.key;
  another_again = .another;  # References secondary.another
  _under = 50;
  path = $PATH;
}
`)

var testConfigString = string(testConfigBytes)
var testConfigReader = bytes.NewReader(testConfigBytes)

var expectedPath = os.Getenv("PATH")

func assertEqual(a interface{}, b interface{}, t *testing.T) {
	if a != b {
		t.Fatal(fmt.Sprintf("'%v' != '%v'", a, b))
	}
}

func assertDirectives(values map[string]interface{}, t *testing.T) {
	// Global
	assertEqual(values["global"], "global value", t)

	// Primary
	primary := values["primary"].(map[string]interface{})
	assertEqual(primary["string"], "primary string value", t)
	assertEqual(primary["string_with_quote"], "some \"quoted\" str\\ing", t)
	assertEqual(primary["single"], "hello world", t)
	assertEqual(primary["single_with_quote"], "'hello' \"world\"", t)
	assertEqual(primary["integer"], int64(500), t)
	assertEqual(primary["float"], float64(80.80), t)
	assertEqual(primary["negative"], int64(-50), t)
	assertEqual(primary["boolean"], true, t)
	assertEqual(primary["not_true"], false, t)
	assertEqual(primary["nothing"], nil, t)
	assertEqual(primary["sec_ref"], int64(50), t)

	// Primary Sub
	sub := primary["sub"].(map[string]interface{})
	assertEqual(sub["key"], "primary sub key value", t)
	assertEqual(sub["included_setting"], "primary sub included_setting value", t)

	// Secondary
	secondary := values["secondary"].(map[string]interface{})
	assertEqual(secondary["another"], "secondary another value", t)
	assertEqual(secondary["global_reference"], "global value", t)
	assertEqual(secondary["primary_sub_key"], "primary sub key value", t)
	assertEqual(secondary["another_again"], "secondary another value", t)
	assertEqual(secondary["_under"], int64(50), t)
	assertEqual(secondary["path"], expectedPath, t)
}

func TestParseBytes(t *testing.T) {
	settings, err := forge.ParseBytes(testConfigBytes)
	if err != nil {
		t.Fatal(err)
	}

	values := settings.ToMap()
	assertDirectives(values, t)
}

func TestParseString(t *testing.T) {
	settings, err := forge.ParseString(testConfigString)
	if err != nil {
		t.Fatal(err)
	}
	values := settings.ToMap()
	assertDirectives(values, t)
}

func TestParseReader(t *testing.T) {
	settings, err := forge.ParseReader(testConfigReader)
	if err != nil {
		t.Fatal(err)
	}
	values := settings.ToMap()
	assertDirectives(values, t)
}

func TestParseFile(t *testing.T) {
	settings, err := forge.ParseFile("./test.cfg")
	if err != nil {
		t.Fatal(err)
	}
	values := settings.ToMap()
	assertDirectives(values, t)
}
