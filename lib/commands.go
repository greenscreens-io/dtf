/**
 * This file represents main library functions called from main module
 */
package lib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"greenscreens-io/dtf/model"
	"io/fs"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// WebCall process Browser handler request - used by Web DTF UI when saving document
func WebCall(data string) error {

	if len(data) == 0 {
		return errors.New("arguments not provided")
	}

	args := strings.Split(strings.TrimPrefix(data, "gs-dtf:"), "/")
	if len(args) < 2 {
		return errors.New("invalid arguments")
	}

	raw, err := base64.URLEncoding.DecodeString(args[0])
	if err != nil {
		return err
	}

	file, err := base64.URLEncoding.DecodeString(args[1])
	if err != nil {
		return err
	}

	var options model.Options
	err = json.Unmarshal(raw, &options)
	if err != nil {
		return err
	}

	ioutil.WriteFile(string(file), []byte(raw), fs.ModePerm)

	return nil
}

// OpenBrowser open url with default browser
func OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// EditConfig - open browser DTF UI with .gsf/.gst file data
func EditConfig(opt *model.Options, fileName string) error {
	data, err := opt.ToBodyRaw()
	if err != nil {
		return err
	}

	fullPath, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	fullPath = base64.StdEncoding.EncodeToString([]byte(fullPath))
	str := base64.StdEncoding.EncodeToString(data)
	ext := filepath.Ext(fileName)
	url := opt.Url + "/dtf/config" + ext + "#" + str + "/" + fullPath
	return OpenBrowser(url)
}
