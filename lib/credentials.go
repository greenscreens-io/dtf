package lib

import (
	"errors"
	"fmt"
	"greenscreens-io/dtf/model"
	"greenscreens-io/dtf/services"
	"os"
)

// CreateCreds create user credentials file for auto login
func CreateCreds() error {

	var data model.AuthData
	var name string
	var err error

	args := os.Args[2:] // without program name

	if len(args) > 0 {
		data.Url = args[0]
	}

	if len(args) > 1 {
		name = args[1]
		data.Key = args[1]
	}

	if len(args) > 2 {
		data.User = args[2]
	}

	if len(args) > 3 {
		data.Password = args[3]
	}

	if len(args) > 4 {
		name = args[4]
	}

	data.Url = PromptURL(data.Url)
	data.Key = PromptKey(data.Key)

	if len(name) == 0 {
		name = Prompt("Enter alias (optional):", false, false)
	}

	fmt.Println("Verifying API Key...")

	partial, err := KeyCheck(data.Url, data.Key)
	if err != nil {
		fmt.Println(err)
		return errors.New("invalid API Key")
	}

	data.Partial = partial

	if !partial {
		fmt.Println("Verification OK!")
		return data.Save(name)
	}

	data.User = PromptUser(data.User)
	data.Password = PromptPassword(data.Password)

	fmt.Println("Verifying user...")

	err = validate(&data)
	if err != nil {
		return err
	}

	fmt.Println("Verification OK!")

	return data.Save(name)
}

// validate data by authorizing client on remote IBM i system
func validate(data *model.AuthData) error {

	var caller services.HTTPCaller

	caller.Init()

	reqData := &model.RequestData{
		Key:      data.Key,
		Url:      data.Url,
		User:     data.User,
		Password: data.Password,
	}

	respAuth, err := caller.Authorize(reqData)
	if err != nil {
		return err
	}

	if !respAuth.Success {
		return errors.New(respAuth.Msg)
	}
	defer caller.Logout(reqData)

	return nil
}

// validate key on gs server
func KeyCheck(url, key string) (bool, error) {
	var caller services.HTTPCaller

	caller.Init()

	reqData := &model.RequestData{
		Key: key,
		Url: url,
	}

	respData, err := caller.KeyCheck(reqData)
	if err != nil {
		return false, err
	}

	return respData.Partial, nil
}

func FileFormat(url, value string) (string, error) {
	var caller services.HTTPCaller

	caller.Init()

	respData, err := caller.FileFormat(url, value)
	if err != nil {
		return "", err
	}

	return respData, nil

}
