package main

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func marshal(c echo.Context) func(interface{}) ([]byte, error) {
	if c.QueryParam("indent") != "true" {
		return json.Marshal
	}
	return func(i interface{}) ([]byte, error) {
		return json.MarshalIndent(i, "", "\t")
	}
}
