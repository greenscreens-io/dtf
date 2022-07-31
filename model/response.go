/**
 * This file represents data structure response from REST API.
 */
package model

import (
	"encoding/json"
	"fmt"
)

var TRACE = false

// ExtJSReponse Rest response format
type ExtJSReponse struct {
	Partial bool   `json:"partial"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Data    []byte `json:"data"`
}

// Load JSON string into structure
func (d *ExtJSReponse) LoadRaw(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *ExtJSReponse) Load(data string) error {
	if TRACE {
		fmt.Print(data)
	}
	return d.LoadRaw([]byte(data))
}
