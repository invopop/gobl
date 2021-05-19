package c14n_test

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/invopop/gobl/c14n"
)

func TestJSONToArray(t *testing.T) {
	// some random data found on the internets
	data := `{
		"actors": [
			{
				"name": "Tom Cruise",
				"age": 56,
				"born_at": "Syracuse, NY",
				"birthdate": "July 3, 1962",
				"photo": "https://jsonformatter.org/img/tom-cruise.jpg",
				"wife": null,
				"weight": 67.5,
				"has_children": true,
				"has_grey_hair": false,
				"children": [
					"Suri",
					"Isabella Jane",
					"Connor"
				],
				"icon": "🤩"
			},
			{
				"name": "Robert Downey Jr.",
				"age": 53,
				"born_at": "New York City, NY",
				"birthdate": "April 4, 1965",
				"photo": "https://jsonformatter.org/img/Robert-Downey-Jr.jpg",
				"wife": {
					"name": "Susan Downey",
					"age": 35
				},
				"weight": 77.1,
				"has_children": true,
				"has_grey_hair": false,
				"children": [
					"Indio Falconer",
					"Avri Roel",
					"Exton Elias"
				]
			}
		]
	}`
	r := strings.NewReader(data)
	obj, err := c14n.UnmarshalJSON(r)
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
		return
	}

	d, err := obj.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect JSON marshal error: %v", err.Error())
		return
	}
	s := fmt.Sprintf("%x", sha256.Sum256(d))
	if s != "f35a55c7bba2df8438802603db442976a8238ceb0a610d1eea38cae1b9fd9013" {
		t.Logf("marshaled data:\n%v\n", string(d))
		t.Errorf("unexpected sum, please check marshaled data, got: %v", s)
	}
}
