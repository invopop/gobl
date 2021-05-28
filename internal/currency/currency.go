package currency

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/invopop/gobl/currency"
)

type CurrencyDoc struct {
	XMLName xml.Name       `xml:"ISO_4217"`
	Table   *CurrencyTable `xml:"CcyTbl"`
}

type CurrencyTable struct {
	Rows []*CurrencyDef `xml:"CcyNtry"`
}

type CurrencyDef struct {
	Name    string `xml:"CcyNm"`      // name of the currency
	Country string `xml:"CtryNm"`     // name of the country it belongs to
	Code    string `xml:"Ccy"`        // three-letter currency code
	Num     string `xml:"CcyNbr"`     // three-digit currency code
	Units   string `xml:"CcyMnrUnts"` // how many cents are used for the currency
}

// GenerateCodes is a special tool function used to convert the source XML
// data into an array of currency definitions.
func GenerateCodes() error {
	data, err := ioutil.ReadFile("./internal/currency/data/iso4217.xml")
	if err != nil {
		return err
	}

	d := new(CurrencyDoc)
	if err := xml.Unmarshal(data, d); err != nil {
		return err
	}

	f, err := os.Create("./currency/codes.go")
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.ParseFiles("./internal/currency/codes.go.tmpl")
	if err != nil {
		return err
	}

	codeSet := make(map[string]int)
	rows := make([]*currency.Def, 0)
	for _, r := range d.Table.Rows {
		u, err := strconv.Atoi(r.Units)
		if err != nil {
			fmt.Printf("skipping row: %v - %v\n", r.Name, r.Code)
			continue
		}
		if _, ok := codeSet[r.Code]; ok {
			codeSet[r.Code]++
		} else {
			codeSet[r.Code] = 1
		}
		def := &currency.Def{
			Name:    r.Name,
			Country: r.Country,
			Code:    currency.Code(r.Code),
			Num:     r.Num,
			Units:   u,
		}
		rows = append(rows, def)
	}

	// Move codeSet to array and sort
	codes := make([]string, len(codeSet))
	i := 0
	for k := range codeSet {
		codes[i] = k
		i++
	}
	sort.Strings(codes)

	fields := make(map[string]interface{})
	fields["Rows"] = rows
	fields["Codes"] = codes
	fields["Date"] = time.Now().UTC().String()
	if err := tmpl.Execute(f, fields); err != nil {
		return err
	}

	return nil
}
