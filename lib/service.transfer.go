package lib

import (
	"bytes"
	"net/http"
)

// Authorize to crate session
func Authorize(data *RequestData) (*http.Request, error) {
	url := data.Url + "/service.transfer/rest/service/v1/login"
	payload, _ := data.ToBody()
	return toRequest(payload, url, data.Key)
}

// Export to create session
func Export(options *Options) (*http.Request, error) {
	url := options.Url + "/service.transfer/rest/service/v1/export"
	payload, _ := options.ToBody()
	return toRequest(payload, url, options.Key)
}

// Import to create session
func Import(options *Options, taskID string, fileID string) (*http.Request, error) {
	url := options.Url + "/service.transfer/rest/service/v1/import?taskID=" + taskID + "&fileID=" + fileID
	payload, _ := options.ToBody()
	return toRequest(payload, url, options.Key)
}

// toRequest - Create http reqeust struct from given params
func toRequest(payload *bytes.Reader, url string, key string) (*http.Request, error) {

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("key", key)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
