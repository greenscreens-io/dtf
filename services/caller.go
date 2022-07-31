package services

import (
	"crypto/tls"
	"errors"
	"greenscreens-io/dtf/model"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
)

// should transfer counter be used
const counter = true

// HTTPCaller struct for holding http client and cookie
type HTTPCaller struct {
	Jar    *cookiejar.Jar
	Client *http.Client
}

// Init initializes HTTP client structure
func (h *HTTPCaller) Init() error {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	h.Jar = jar

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       h.Jar,
		Transport: tr,
	}

	h.Client = client

	return nil
}

// fileExists chack if file exists
func (h *HTTPCaller) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// callRaw makes remote http request and get reponse in string format
func (h *HTTPCaller) callRaw(req *http.Request) ([]byte, error) {

	res, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// call makes remote http request and get reponse in string format
func (h *HTTPCaller) call(req *http.Request) (*model.ExtJSReponse, error) {

	data, err := h.callRaw(req)
	if err != nil {
		return nil, err
	}

	resp := &model.ExtJSReponse{}
	err = resp.LoadRaw(data)
	return resp, err
}

func (h *HTTPCaller) validate(path string, fileMode int) error {
	// allow fileMode -1 as no check
	if fileMode > 2 {
		return errors.New("invalid file mode")
	}

	// 0 - NEW - only if file does not exist
	// 1 - APPEND - append to file
	// 2 - OVERWRITE - only if file exist
	isFound := h.fileExists(path)
	if isFound && fileMode == 0 {
		return errors.New("file already exist and file mode is NEW")
	}

	if !isFound && fileMode == 1 {
		return errors.New("file does not exist and file mode is APPEND")
	}

	if !isFound && fileMode == 2 {
		return errors.New("file does not exist and file mode is OVERWRITE")
	}

	// remove file if overwrite
	if isFound && fileMode != 1 {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// callDownload makes remote http request and get reponse in string format
func (h *HTTPCaller) callDownload(req *http.Request, path string, fileMode int) error {

	err := h.validate(path, fileMode)
	if err != nil {
		return err
	}

	mode := os.O_APPEND | os.O_CREATE | os.O_WRONLY

	// 0 - new, 1 - append,  2- overwrite
	/*
		switch (fileMode) {
		case 1:
			mode = os.O_APPEND|os.O_CREATE|os.O_WRONLY
			break;
		case 2:
			mode = os.O_CREATE|os.O_WRONLY
			break;
		}
	*/

	out, err := os.OpenFile(path, mode, 0600)
	//out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	res, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		msg := res.Status
		b, err := io.ReadAll(res.Body)
		if err == nil && len(b) > 0 {
			msg = res.Status + " -> " + string(b)
		}
		return errors.New(msg)
	}

	sizeStr := res.Header.Get("Content-Size")
	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		size = 0
	}

	counter := &HTTPCounter{Total: uint64(size)}
	defer counter.Finish()
	_, err = io.Copy(out, io.TeeReader(res.Body, counter))
	if err != nil {
		return err
	}

	return nil
}

// Share file to service.file
func (h *HTTPCaller) Share(data *model.RequestData) (*model.ExtJSReponse, error) {
	req, err := Share(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

// Upload IFS file to service.file
func (h *HTTPCaller) Upload(data *model.RequestData) (*model.ExtJSReponse, error) {
	req, err := Upload(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

// Download IFS file from service.file
func (h *HTTPCaller) Download(data *model.RequestData, localFile string, fileMode int) error {
	req, err := Download(data)
	if err != nil {
		return err
	}
	return h.callDownload(req, localFile, fileMode)
}

// Receive exported file from service.file
func (h *HTTPCaller) Receive(data *model.RequestData, taskID string, fileID string, fileMode int) error {
	req, err := Receive(data, taskID, fileID)
	if err != nil {
		return err
	}
	return h.callDownload(req, data.Path, fileMode)
}

// Authorize service.transfer
func (h *HTTPCaller) Authorize(data *model.RequestData) (*model.ExtJSReponse, error) {
	req, err := Authorize(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

func (h *HTTPCaller) KeyCheck(data *model.RequestData) (*model.ExtJSReponse, error) {
	req, err := KeyCheck(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

func (h *HTTPCaller) FileFormat(service, value string) (string, error) {
	req, err := FileFormat(service, value)
	if err != nil {
		return "", err
	}
	bin, err := h.callRaw(req)
	if err != nil {
		return "", err
	}
	return string(bin), nil
}

// Logout user session
func (h *HTTPCaller) Logout(data *model.RequestData) (*model.ExtJSReponse, error) {
	req, err := Logout(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

func (h *HTTPCaller) Describe(data *model.RequestData, lib, obj string) (*model.DescriptionResponse, error) {
	req, err := Describe(data, lib, obj)
	if err != nil {
		return nil, err
	}

	raw, err := h.callRaw(req)
	if err != nil {
		return nil, err
	}

	resp := &model.DescriptionResponse{}
	err = resp.LoadRaw(raw)
	return resp, err
}

// Export file with service.transfer
func (h *HTTPCaller) Export(options *model.Options) (*model.ExtJSReponse, error) {
	req, err := Export(options)
	if err != nil {
		return nil, err
	}

	if options.Streaming {
		return nil, h.callDownload(req, options.FullPath(), options.FileMode())
	}

	return h.call(req)
}

// Import file with service.transfer
func (h *HTTPCaller) Import(options *model.Options, taskID, fileID string) (*model.ExtJSReponse, error) {

	var err error
	var req *http.Request

	if options.Streaming {
		req, err = ImportStreaming(options)
	} else {
		req, err = ImportCached(options, taskID, fileID)
	}

	if err != nil {
		return nil, err
	}

	return h.call(req)
}
