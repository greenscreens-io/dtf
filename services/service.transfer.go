/**
 * This file represents REST API for service.transfer
 * Prepare http.Request structure for remote calls
 */
package services

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"greenscreens-io/dtf/model"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

func FileFormat(service, value string) (*http.Request, error) {
	par := b64.URLEncoding.EncodeToString([]byte(value))
	url := service + "/service.transfer/format?q=" + par
	return toRequestSimple(url, "")
}

func KeyCheck(data *model.RequestData) (*http.Request, error) {
	url := data.Url + "/services/keycheck?key=" + data.Key
	return toRequestSimple(url, data.Key)
}

// Authorize to crate session
func Authorize(data *model.RequestData) (*http.Request, error) {
	url := data.Url + "/service.transfer/rest/service/v1/login"
	payload, _ := data.ToBody()
	return toRequest(payload, url, data.Key)
}

func Logout(data *model.RequestData) (*http.Request, error) {
	url := data.Url + "/service.transfer/rest/service/v1/logout"
	return toRequest(nil, url, data.Key)
}

func Describe(data *model.RequestData, lib, obj string) (*http.Request, error) {
	url := data.Url + "/service.transfer/rest/service/v1/describe/" + lib + "/" + obj
	return toRequest(nil, url, data.Key)
}

// Export to create session
func Export(options *model.Options) (*http.Request, error) {
	url := options.Url + "/service.transfer/rest/service/v1/export"
	payload, _ := options.ToBody()
	return toRequest(payload, url, options.Key)
}

// Import to create session
func ImportCached(options *model.Options, taskID string, fileID string) (*http.Request, error) {
	url := "/service.transfer/rest/service/v1/import"
	url = options.Url + url + "?taskID=" + taskID + "&fileID=" + fileID
	payload, _ := options.ToBody()
	return toRequest(payload, url, options.Key)
}

func ImportStreaming(options *model.Options) (*http.Request, error) {
	url := "/service.transfer/rest/service/v1/import"
	return toUploadRequest(options, url)
}

// toRequest - Create http reqeust struct from given params
func toRequestSimple(url string, key string) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("key", key)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func toRequest(payload io.Reader, url string, key string) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("key", key)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// send file to GS server
func toUploadRequest(options *model.Options, urlPart string) (*http.Request, error) {

	fp, fs, err := getFilePointer(options.FullPath())
	if err != nil {
		return nil, err
	}
	stat := *fs

	json, err := options.ToBodyRaw()
	if err != nil {
		return nil, err
	}

	// Reduce number of syscalls when reading from disk.
	bufferedFileReader := bufio.NewReader(fp)

	// Create a pipe for writing from the file and reading to
	// the request concurrently.
	bodyReader, bodyWriter := io.Pipe()
	formWriter := multipart.NewWriter(bodyWriter)

	url := options.Url + urlPart
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("key", options.Key)
	req.Header.Set("Content-Type", formWriter.FormDataContentType())

	// Store the first write error in writeErr.
	var (
		writeErr error
		errOnce  sync.Once
	)
	setErr := func(err error) {
		if err != nil {
			errOnce.Do(func() {
				writeErr = err
				fmt.Println(writeErr)
			})
		}
	}

	go func() {

		_ = formWriter.WriteField("json", string(json))

		partWriter, err := formWriter.CreateFormFile("file", stat.Name())
		setErr(err)

		if counter {
			counter := &HTTPCounter{Total: uint64(stat.Size())}
			defer counter.Finish()
			_, err = io.Copy(partWriter, io.TeeReader(bufferedFileReader, counter))
		} else {
			_, err = io.Copy(partWriter, bufferedFileReader)
		}

		setErr(err)
		setErr(formWriter.Close())
		setErr(bodyWriter.Close())
	}()

	return req, nil
}

func getFilePointer(filePath string) (*os.File, *os.FileInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	return file, &fi, nil
}
