package c14n_test

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/c14n"
)

func ExampleCanonicalJSON() {
	d := `{ "foo":"bar", "c": 123.4, "a": 56, "b": 0.0, "y":null}`
	r := strings.NewReader(d)
	res, err := c14n.CanonicalJSON(r)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v\n", string(res))
	// Output:
	// {"a":56,"b":0.0E0,"c":1.234E2,"foo":"bar"}

}
