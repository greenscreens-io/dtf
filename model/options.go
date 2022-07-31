/**
 * This file represents data structure for *.gst and *.gsf files.
 * Files are managed by Web UI app available in provided service.transfer module
 */
package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Options - Json structure for DTF options
type Options struct {
	Url              string        `json:"url"`
	Key              string        `json:"key"`
	User             string        `json:"user"`
	Password         string        `json:"password"`
	Charset          string        `json:"charset"`
	Library          string        `json:"library"`
	Object           string        `json:"object"`
	Member           string        `json:"member"`
	TplLibrary       string        `json:"libraryTpl"`
	TplObject        string        `json:"objectTpl"`
	TplMember        string        `json:"memberTpl"`
	CodePage         string        `json:"codePage"`
	IgnoreErrors     bool          `json:"ignoreErrors"`
	Truncate         bool          `json:"truncate"`
	DecFloatChar     bool          `json:"decFloatChar"`
	Limit            int           `json:"limit"`
	Skip             int           `json:"skip"`
	FieldDelimiter   string        `json:"fieldDelimiter"`
	LineDelimiter    string        `json:"lineDelimiter"`
	Header           int           `json:"header"`
	PadNumeric       int           `json:"padNumeric"`
	FileFormat       int           `json:"fileFormat"`
	DateFormat       int           `json:"dateFormat"`
	TimeFormat       int           `json:"timeFormat"`
	DateSeparator    int           `json:"dateSeparator"`
	TimeSeparator    int           `json:"timeSeparator"`
	DecimalSeparator int           `json:"decimalSeparator"`
	Mode             int           `json:"mode"`
	Type             int           `json:"type"`
	Strict           bool          `json:"strict"`
	Columns          []string      `json:"columns"`
	Path             string        `json:"path"`
	IFS              string        `json:"ifs"`
	SQL              string        `json:"sql"`
	Streaming        bool          `json:"streaming"`
	Defs             []Description `json:"fieldDefs"`
	fullPath         string
	hasCreds         bool
	// is api key partial services/keycheck?key=API_KEY
	Partial bool
}

func (d *Options) FileMode() int {
	if d.Strict {
		return d.Mode
	}
	return -1
}

// ToBody return structure converted to JSON byte array
func (d *Options) ToBodyRaw() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Options) ToBody() (*bytes.Reader, error) {

	payloadBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payloadBytes), nil
}

// IsValid check if all auth data is ok
func (d *Options) IsValid() bool {
	if d.Partial {
		return len(d.Key) > 0 && len(d.Url) > 0 && len(d.User) > 0 && len(d.Password) > 0
	}
	return len(d.Key) > 0 && len(d.Url) > 0
}

// HasCredentials return true if laoded cached creds
func (d *Options) HasCredentials() bool {
	return d.hasCreds
}

func (d *Options) FullPath() string {
	if len(d.Path) == 0 {
		fullPath, err := filepath.Abs(d.Path)
		if err != nil {
			return fullPath
		}
	}
	return d.Path
}

// Load read file from given path and convert into JSON data
func (d *Options) Load(path string) error {

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &d)
	if err != nil {
		return err
	}

	d.fullPath = fullPath

	if len(d.Key) == 0 {
		d.Key = os.Getenv("DTF_KEY")
	}

	d.hasCreds = d.loadCreds()

	d.initFromEnv()
	d.initFromArgs()

	return nil
}

// Save auth cache only if not found
func (d *Options) Save() {

	if d.HasCredentials() {
		return
	}

	if !d.IsValid() {
		return
	}

	var data AuthData
	err := data.Load(d.Key)
	if err == nil {
		return
	}

	data.Url = d.Url
	data.Key = d.Key
	data.User = d.User
	data.Password = d.Password
	err = data.Save(d.Key)
	if err != nil {
		return
	}

	d.saveSafe()

}

func (d *Options) saveSafe() (bool, error) {

	if len(d.fullPath) == 0 {
		return false, nil
	}

	user := d.User
	pass := d.Password
	d.Password = ""
	d.User = ""
	data, err := json.MarshalIndent(d, "", "  ")
	d.Password = user
	d.User = pass
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(d.fullPath, data, fs.ModePerm)

	return err == nil, err
}

// InitFromEnv init args from env variables
func (d *Options) initFromEnv() {

	env_key := os.Getenv("DTF_KEY")
	env_user := os.Getenv("DTF_USER")
	env_pass := os.Getenv("DTF_PWD")

	if len(env_key) > 0 {
		d.Key = env_key
	}

	if len(env_user) > 0 {
		d.User = env_user
	}

	if len(env_pass) > 0 {
		d.User = env_pass
	}
}

// InitFromArgs init from program start arguments
func (d *Options) initFromArgs() {
	args := os.Args[1:] // without program name
	if len(args) >= 3 {
		if args[1] == "-u" {
			d.User = args[2]
		}
		if args[1] == "-p" {
			d.Password = args[2]
		}
	}

	if len(args) >= 5 {
		if args[3] == "-u" {
			d.User = args[4]
		}
		if args[3] == "-p" {
			d.Password = args[4]
		}
	}
}

// loadCreds load user creds and populate options
func (d *Options) loadCreds() bool {
	var data AuthData
	err := data.Load(d.Key)
	if err != nil {
		fmt.Println(err)
		return false
	}

	d.Key = data.Key
	d.User = data.User
	d.Password = data.Password

	if len(d.Url) == 0 {
		d.Url = data.Url
	}

	return true
}
