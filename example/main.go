package main

import (
	"encoding/json"
	"fmt"

	"github.com/brettlangdon/forge"
)

func main() {
	settings, err := forge.ParseFile("example.cfg")
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
