/**
 * This file represents data structure response from REST API.
 */
package model

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// ExtJSReponse Rest response format
type DescriptionResponse struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg"`
	Code    string        `json:"code"`
	Type    string        `json:"type"`
	Data    []Description `json:"data"`
}

type Description struct {
	Name  string            `json:"name"`
	Items []DescriptionItem `json:"items"`
}

type DescriptionItem struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Len  int    `json:"len"`
}

// Load JSON string into structure
func (d *DescriptionResponse) LoadRaw(data []byte) error {
	return json.Unmarshal(data, d)
}

// return received description into savable data
func (d *DescriptionResponse) ToRawData() ([]byte, error) {
	return json.Marshal(d.Data)
}

func LoadFDF(file string) (*[]Description, error) {

	fullPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadFile(fullPath + "fd")
	if err != nil {
		return nil, nil
	}

	var defs []Description
	err = json.Unmarshal(byteValue, &defs)
	if err != nil {
		return nil, err
	}
	return &defs, nil
}
