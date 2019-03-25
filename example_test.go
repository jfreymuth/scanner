package scanner_test

import (
	"fmt"

	"github.com/jfreymuth/scanner"
)

func Example() {
	const input = `

	// single line comment

	/*
	multi line
	comment
	*/
	
	exampleString = "This is an example \u263a"
	exampleNumber = 12.5
	exampleList = ["test", 2]
	
	`

	sc := scanner.FromString(input)
	for !sc.End() {
		name := sc.Ident()
		sc.Demand("=")
		value := parseValue(sc)
		fmt.Printf("%s: %v\n", name, value)
	}

	if sc.Err() != nil {
		fmt.Println("An error occurred:")
		if err, ok := sc.Err().(*scanner.Error); ok {
			fmt.Println(err.PositionIndicator())
		} else {
			fmt.Println(sc.Err())
		}
	}

	// Output:
	// exampleString: This is an example ☺
	// exampleNumber: 12.5
	// exampleList: [test 2]
}

func parseValue(sc *scanner.Scanner) interface{} {
	switch {
	case sc.Is("\""):
		return sc.String()
	case sc.Eat("["):
		var list []interface{}
		for !sc.Eat("]") {
			list = append(list, parseValue(sc))
			if !sc.Eat(",") {
				sc.Demand("]")
				break
			}
		}
		return list
	case sc.IsInt():
		return sc.Int()
	case sc.IsFloat():
		return sc.Float()
	default:
		sc.Fail("invalid value")
		return nil
	}
}
