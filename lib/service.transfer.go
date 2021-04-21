package lib

import (
	"net/http"
)

// Authorize to crate session - do not use for now
func Authorize(data *RequestData) (*http.Request, error) {

	url := data.Url + "/service.transfer/rest/service/v1/login"

	payload, _ := data.ToBody()

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("key", data.Key)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Export to create session - do not use for now
func Export(options *Options) (*http.Request, error) {

	url := options.Url + "/service.transfer/rest/service/v1/export"

	payload, _ := options.ToBody()

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("key", options.Key)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Import to create session - do not use for now
func Import(options *Options, taskID string, fileID string) (*http.Request, error) {

	url := options.Url + "/service.transfer/rest/service/v1/import?taskID=" + taskID + "&fileID=" + fileID

	payload, _ := options.ToBody()

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("key", options.Key)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
