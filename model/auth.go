/**
 * This file represents data structure used for encrypted cached authorized API Key.
 * Data contains API Key,  user credentials and service URL
 */
package model

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/denisbrodbeck/machineid"
)

// AuthData structure encrypted and stored in API cached key
type AuthData struct {
	Url      string `json:"url"`
	Uid      string `json:"uid"`
	Winuser  string `json:"winuser"`
	Key      string `json:"key"`
	User     string `json:"user"`
	Password string `json:"password"`
	Partial  bool
}

const appID = "dtf.exe"

// getFile retrieve cached API key path,
// if not specified and only one file exist, use that
func (d *AuthData) getFile(apiKey string) string {
	usr, _ := user.Current()
	dirName := filepath.Join(usr.HomeDir, ".dtf")

	if len(apiKey) > 0 {
		return filepath.Join(dirName, apiKey+".dtfkey")
	}

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return ""
	}

	if len(files) == 1 {
		fileName := files[0].Name()
		//newKey := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
		return filepath.Join(dirName, fileName)
	}

	return ""
}

// Load encrypted cached authorized API key
func (d *AuthData) Load(name string) error {
	id, _ := machineid.ID()
	usr, _ := user.Current()
	key := Protect(appID, id)

	fileName := d.getFile(name)
	byteValue, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	byteValue, err = DecryptMessageRaw(key, byteValue)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &d)
	if err != nil {
		return err
	}

	if d.Uid != id {
		return fmt.Errorf("invalid authentication file : machine not matched")
	}

	if d.Winuser != usr.Uid {
		return fmt.Errorf("invalid authentication file : user not matched")
	}

	return nil
}

// Save authorized API key
func (d *AuthData) Save(name string) error {
	id, _ := machineid.ID()
	usr, _ := user.Current()

	d.Uid = id
	d.Winuser = usr.Uid

	payloadBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}

	dirName := filepath.Join(usr.HomeDir, ".dtf")
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	fileName := filepath.Join(dirName, name+".dtfkey")

	key := Protect(appID, id)
	payloadBytes, err = EncryptMessageRaw(key, payloadBytes)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, payloadBytes, fs.ModePerm)
}
