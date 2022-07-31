/**
 * This file represents REST API for service.file
 * Prepare http.Request structure for remote calls
 */
package services

import (
	"bufio"
	"fmt"
	"greenscreens-io/dtf/model"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
)

// uploadRequest prepare httpRequest for sending a file to GS server
func uploadRequest(data *model.RequestData, urlPart string) (*http.Request, error) {

	data.Streaming = true

	// Reduce number of syscalls when reading from disk.
	bufferedFileReader := bufio.NewReader(data.GetFilePointer())

	// Create a pipe for writing from the file and reading to
	// the request concurrently.
	bodyReader, bodyWriter := io.Pipe()
	formWriter := multipart.NewWriter(bodyWriter)

	url := data.Url + urlPart
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("key", data.Key)
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
		_ = formWriter.WriteField("streaming", "true")
		_ = formWriter.WriteField("user", data.User)
		_ = formWriter.WriteField("password", data.Password)
		_ = formWriter.WriteField("path", data.Path)
		_ = formWriter.WriteField("fileSize", strconv.FormatInt(data.FileSize, 10))

		partWriter, err := formWriter.CreateFormFile("file", data.FileName)
		setErr(err)

		if counter {
			counter := &HTTPCounter{Total: uint64(data.FileSize)}
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

// Share file to GS Server, response returns file hash
func Share(data *model.RequestData) (*http.Request, error) {
	return uploadRequest(data, "/service.file/rest/lite/v1/share")
}

// Upload send IFS file to GS Server, response returns file hash
func Upload(data *model.RequestData) (*http.Request, error) {
	return uploadRequest(data, "/service.file/rest/lite/v1/upload")
}

// Download IFS file from GS Server
func Download(data *model.RequestData) (*http.Request, error) {

	data.Streaming = true
	url := data.Url + "/service.file/rest/lite/v1/download"

	body, _ := data.ToBody()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("key", data.Key)

	return req, nil
}

// Receive shared file from GS Server
func Receive(data *model.RequestData, taskID string, fileID string) (*http.Request, error) {

	url := data.Url + "/service.file/rest/lite/v1/get/" + taskID + "/" + fileID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("key", data.Key)

	return req, nil
}
