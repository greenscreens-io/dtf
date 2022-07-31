/**
 * This file represents data structure used for JSON REST request
 */
package model

import (
	"bytes"
	"encoding/json"
	"os"
)

// RequestData structure for http client requests
type RequestData struct {
	Url       string        `json:"url"`
	Key       string        `json:"key"`
	User      string        `json:"user"`
	Password  string        `json:"password"`
	FileName  string        `json:"fileName"`
	FileSize  int64         `json:"fileSize"`
	Path      string        `json:"path"`
	Streaming bool          `json:"streaming"`
	Data      []Description `json:"fieldDefs"`
	file      *os.File
}

// GetFilePointer input / output file
func (d *RequestData) GetFilePointer() *os.File {
	return d.file
}

// InitFilePointer initialize data source or target
func (d *RequestData) InitFilePointer(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	d.StoreFilePointer(file)

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	d.FileSize = fi.Size()
	d.FileName = fi.Name()
	return file, nil
}

// StoreFilePointer store file pointer for later to be closed
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
