package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aleury/goini/parser"
)

func main() {
	name := "test"
	input := `
key=abcdefg

[user]
name=Adam Eury
age=35
job=Software Engineer
email=adam@test.com

[address]
street=1800 Test Lane
city=Testy Mctestersonville
state=North Carolina
zip=90210
`
	file := parser.Parse(name, input)
	data, err := json.MarshalIndent(file, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
