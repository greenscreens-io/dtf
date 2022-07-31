/**
 * This file represents main library functions called from main module
 * Used for importing data to IBM i through GS server
 */
package transfer

import (
	"errors"
	"greenscreens-io/dtf/model"
	"greenscreens-io/dtf/services"
)

// ToGS - transfer data to GS
// upload file, get token, call transfer
func ToGS(opt *model.Options) error {
	if len(opt.IFS) > 0 {
		return toGSFile(opt)
	}
	return toGSData(opt)
}

// toGSFile upload file to IFS path
func toGSFile(opt *model.Options) error {
	var caller services.HTTPCaller

	caller.Init()

	data := &model.RequestData{
		Key:      opt.Key,
		Url:      opt.Url,
		User:     opt.User,
		Password: opt.Password,
		Path:     opt.IFS,
	}
	_, err := data.InitFilePointer(opt.FullPath())
	if err != nil {
		return err
	}
	defer data.CloseFilePointer()

	resp, err := caller.Upload(data)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Msg)
	}

	return err
}

// toGSData import data to IBM server
func toGSData(opt *model.Options) error {

	var caller services.HTTPCaller

	caller.Init()

	// streaming file SQL/MBR import is not available
	data := &model.RequestData{
		Key:      opt.Key,
		Url:      opt.Url,
		User:     opt.User,
		Password: opt.Password,
	}

	_, err := data.InitFilePointer(opt.FullPath())
	if err != nil {
		return err
	}
	defer data.CloseFilePointer()

	// authorize API key and get session cookie
	respAuth, err := caller.Authorize(data)
	if err != nil {
		return err
	}

	if !respAuth.Success {
		return errors.New(respAuth.Msg)
	}
	defer caller.Logout(data)

	data.Path = opt.FullPath()

	// call service.file/upload -> returns file hashID
	code, msg, err := share(&caller, data, opt.Streaming)
	if err != nil {
		return err
	}

	// call service.transfer/import
	opt.Streaming = false
	respImp, err := caller.Import(opt, code, msg)
	if err != nil {
		return err
	}

	if !respImp.Success {
		return errors.New(respImp.Msg)
	}

	return nil
}

func share(caller *services.HTTPCaller, data *model.RequestData, streaming bool) (string, string, error) {
	if streaming {
		return "", "", nil
	}
	respUpl, err := caller.Share(data)
	if err != nil {
		return "", "", err
	}

	if !respUpl.Success {
		return "", "", errors.New(respUpl.Msg)
	}

	return respUpl.Code, respUpl.Msg, nil
}
