package lib

import (
	"io"
	"os"
	"bytes"
	"strconv"
	"net/http"
	"mime/multipart"
)

// send file to GS server
func send(data  *RequestData, urlPart string) (*http.Request, error) {

	file, err := os.Open(data.Path)
	if err != nil {
		return nil, err
	}
	data.StoreFilePointer(file)
	// defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data.FileSize = fi.Size()
	data.FileName = fi.Name()
	data.Streaming = true

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("user", data.User)
	_ = writer.WriteField("password", data.Password)
	_ = writer.WriteField("fileSize", strconv.FormatInt(data.FileSize, 10))

	part, err := writer.CreateFormFile("file", data.FileName)
	if err != nil {
		return nil, err
	}

	counter := &HTTPCounter{Total:uint64(data.FileSize)}
	defer counter.Finish()
	_, err = io.Copy(part, io.TeeReader(file, counter))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	url := data.Url + urlPart
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("key", data.Key)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

// Share file to GS Server, response returns file hash
func Share(data  *RequestData) (*http.Request, error) {
	return send(data, "/service.file/rest/lite/v1/share")
}

// Upload send IFS file to GS Server, response returns file hash
func Upload(data  *RequestData) (*http.Request, error) {
	return send(data, "/service.file/rest/lite/v1/upload")
}


// Download IFS file from GS Server
func Download(data  *RequestData, taskID string, fileID string) (*http.Request, error) {

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
func Receive(data  *RequestData, taskID string, fileID string) (*http.Request, error) {

	url := data.Url + "/service.file/rest/lite/v1/get/" + taskID + "/" + fileID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("key", data.Key)

	return req, nil
}
