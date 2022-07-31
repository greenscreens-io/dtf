/**
 * This file represents main library functions called from main module
 * Used for exporting data from IBM i through GS server
 */
package transfer

import (
	"errors"
	"greenscreens-io/dtf/model"
	"greenscreens-io/dtf/services"
	"os"
	"path/filepath"
)

// FromGS - transfer data from GS
// call transfer, read downloaded stream
func FromGS(opt *model.Options) error {

	err := verifyFile(opt)
	if err != nil {
		return err
	}

	if len(opt.IFS) > 0 {
		return fromGSFile(opt)
	}

	return fromGSData(opt)
}

// fromGSFile does IFS file transfer
func fromGSFile(opt *model.Options) error {
	var caller services.HTTPCaller

	caller.Init()

	data := &model.RequestData{
		Key:      opt.Key,
		Url:      opt.Url,
		User:     opt.User,
		Password: opt.Password,
		Path:     opt.IFS,
	}
	defer data.CloseFilePointer()

	// call service.file/download -> returns file from IFS
	return caller.Download(data, opt.FullPath(), opt.FileMode())
}

// fromGSData does data transfer
func fromGSData(opt *model.Options) error {

	var caller services.HTTPCaller

	caller.Init()

	data := &model.RequestData{
		Key:      opt.Key,
		Url:      opt.Url,
		User:     opt.User,
		Password: opt.Password,
	}

	// authorize API key and get session cookie
	respAuth, err := caller.Authorize(data)
	if err != nil {
		return err
	}

	if !respAuth.Success {
		return errors.New(respAuth.Msg)
	}

	defer caller.Logout(data)

	// call service.transfer/export -> returns file hashID
	respExp, err := caller.Export(opt)
	if err != nil {
		return err
	}

	if opt.Streaming {
		return nil
	}

	if !respExp.Success {
		return errors.New(respExp.Msg)
	}

	_, file := filepath.Split(opt.FullPath())
	data.Path = opt.FullPath()
	data.FileName = file

	// call service.file/download to get the file ( opt, taskID, fileID, fileMode)
	err = caller.Receive(data, respExp.Code, respExp.Msg, opt.FileMode())

	if err != nil {
		return err
	}

	return nil
}

func verifyFile(opt *model.Options) error {

	path := opt.FullPath()
	_, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		return nil
	}

	if opt.Strict && opt.Mode == 0 {
		return errors.New("file already exists, transfer mode is NEW")
	}

	// overwrite
	if opt.Mode == 2 {
		os.Remove(path)
	}

	return nil
}
