package lib

import (
	"os"
	"fmt"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

// ExtJSReponse Rest response format
type ExtJSReponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Code    string `json:"code"`
	Type    string `json:"type"`
}

// Options - Json structure for DTF options
type Options struct {
	Url              string   `json:"url"`
	Key              string   `json:"key"`
	User             string   `json:"user"`
	Password         string   `json:"password"`
	Charset          string   `json:"charset"`
	Library          string   `json:"library"`
	Object           string   `json:"object"`
	Member           string   `json:"member"`
	CodePage         string   `json:"codePage"`
	IgnoreErrors     bool     `json:"ignoreErrors"`
	Truncate         bool     `json:"truncate"`
	DecFloatChar     bool     `json:"decFloatChar"`
	Limit            int      `json:"limit"`
	Skip             int      `json:"skip"`
	FieldDelimiter   string   `json:"fieldDelimiter"`
	LineDelimiter    string   `json:"lineDelimiter"`
	Header           int      `json:"header"`
	PadNumeric       int      `json:"padNumeric"`
	FileFormat       int      `json:"fileFormat"`
	DateFormat       int      `json:"dateFormat"`
	TimeFormat       int      `json:"timeFormat"`
	DateSeparator    int      `json:"dateSeparator"`
	TimeSeparator    int      `json:"timeSeparator"`
	DecimalSeparator int      `json:"decimalSeparator"`
	Mode             int      `json:"mode"`
	Columns          []string `json:"columns"`
	Path             string   `json:"path"`
}

// RequestData structure for http client requests
type RequestData struct {
	Url       string `json:"url"`
	Key       string `json:"key"`
	User      string `json:"user"`
	Password  string `json:"password"`
	FileName  string `json:"fileName"`
	FileSize  int64  `json:"fileSize"`
	Path      string `json:"path"`
	Streaming bool   `json:"streaming"`
	file	  *os.File
}

// TRACE to log JSON responses
var TRACE = false


// StoreFilePointer store file pointer fol later to be closed
func (d *RequestData) StoreFilePointer(file *os.File) {
	if d.file != nil {
		d.file.Close()
	}
	d.file = file
}

// CloseFilePointer postponed file close after request
func (d *RequestData) CloseFilePointer() {
	if d.file != nil {
		d.file.Close()
		d.file = nil
	}
}

// ToBody return structure converted to JSON byte array
func (d *RequestData) ToBody() (*bytes.Reader, error) {

	payloadBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payloadBytes), nil
}

// Load JSON string into structure
func (d *ExtJSReponse) Load(data []byte) error {
	if TRACE {
		fmt.Print(string(data[:]))
	}
	return json.Unmarshal(data, d)
}

// ToBody return structure converted to JSON byte array
func (d *Options) ToBody() (*bytes.Reader, error) {

	payloadBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payloadBytes), nil
}

// Load read file from given path and convert into JSON data
func (d *Options) Load(path string) error {

	byteValue, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, &d)
}

// Load will read config file to get service url, api key and login params
func (d *RequestData) Load() error {

	byteValue, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, &d)
}
