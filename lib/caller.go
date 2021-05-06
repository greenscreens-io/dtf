package lib

import (
	"os"
	"io"
	"errors"
	"strconv"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

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

	client := &http.Client{
		Jar: h.Jar,
	}

	h.Client = client

	return nil
}

// call makes remote http request and get reponse in string format
func (h *HTTPCaller) call(req *http.Request) (*ExtJSReponse, error) {

	resp := &ExtJSReponse{}

	res, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(string(body[:]))
	}

	err = resp.Load(body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *HTTPCaller) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// callDownload makes remote http request and get reponse in string format
func (h *HTTPCaller) callDownload(req *http.Request, path string, fileMode int) (error) {

	// allow fileMode -1 as no check
	if fileMode > 2 {
		return errors.New("invalid file mode")
	}

	// 0 - NEW - only if file does not exist
	// 1 - APPEND - append to file
	// 2 - OVERWRITE - only if file exist
	isFound := h.fileExists(path);
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
	if isFound && fileMode == 2 {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	out, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	//out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return err
	}

	res, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	sizeStr := res.Header.Get("Content-Size")
	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err !=  nil {
		size = 0
	}

	counter := &HTTPCounter{Total:uint64(size)}
	defer counter.Finish()
	_, err = io.Copy(out, io.TeeReader(res.Body, counter))
	if err != nil {
		return err
	}

	return nil
}

// callUpload makes remote http request and get reponse in string format
func (h *HTTPCaller) callUpload(req *http.Request, path string) (error) {
	return nil
}

// CallRaw makes remote http request and get reponse in string format
func (h *HTTPCaller) callRaw(req *http.Request) (string, error) {

	res, err := h.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New(string(body[:]))
	}

	return string(body[:]), nil
}

// Share file to service.file
func (h *HTTPCaller) Share(data *RequestData) (*ExtJSReponse, error) {
	req, err := Share(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

// Upload IFS file to service.file
func (h *HTTPCaller) Upload(data *RequestData) (*ExtJSReponse, error) {
	req, err := Upload(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

// Download IFS file from service.file
func (h *HTTPCaller) Download(data *RequestData, taskID string, fileID string, fileMode int) (error) {
	req, err := Download(data, taskID, fileID)
	if err != nil {
		return err
	}
	return h.callDownload(req, data.Path, fileMode)
}

// Receive exported file from service.file
func (h *HTTPCaller) Receive(data *RequestData, taskID string, fileID string, fileMode int) (error) {
	req, err := Receive(data, taskID, fileID)
	if err != nil {
		return err
	}
	return h.callDownload(req, data.Path, fileMode)
}

// Authorize service.transfer
func (h *HTTPCaller) Authorize(data *RequestData) (*ExtJSReponse, error) {
	req, err := Authorize(data)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

// Export file with service.transfer
func (h *HTTPCaller) Export(options *Options) (*ExtJSReponse, error) {
	req, err := Export(options)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}

// Import file with service.transfer
func (h *HTTPCaller) Import(options *Options, taskID, fileID string) (*ExtJSReponse, error) {
	req, err := Import(options, taskID, fileID)
	if err != nil {
		return nil, err
	}
	return h.call(req)
}
