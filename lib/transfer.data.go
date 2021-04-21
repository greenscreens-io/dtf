package lib

import(
	"errors"
	"path/filepath"
)

// ToGS - transfer data to GS
// upload file, get token, call transfer
func ToGS(opt *Options) error {

	var caller HTTPCaller

	caller.Init()

	data := &RequestData{
		Key:      opt.Key,
		Url:      opt.Url,
		User:     opt.User,
		Password: opt.Password,
	}
	defer data.CloseFilePointer()

	// authorize API key and get session cookie
	respAuth, err := caller.Authorize(data)
	if err != nil {
		return err
	}
	if respAuth.Success == false {
		return errors.New(respAuth.Msg)
	}

	data.Path = opt.Path
	// call service.file/upload -> returns file hashID
	respUpl, err := caller.Share(data)
	if err != nil {
		return err
	}

	if respUpl.Success == false {
		return errors.New(respUpl.Msg)
	}

	// call service.transfer/import -> returns file hashID
	respImp, err := caller.Import(opt, respUpl.Code, respUpl.Msg)
	if err != nil {
		return err
	}
	if respImp.Success == false {
		return errors.New(respImp.Msg)
	}

	return nil
}

// FromGS - transfer data from GS
// call transfer, read downloaded stream
func FromGS(opt *Options) error {

	var caller HTTPCaller

	caller.Init()

	data := &RequestData{
		Key:      opt.Key,
		Url:      opt.Url,
		User:     opt.User,
		Password: opt.Password,
	}
	defer data.CloseFilePointer()

	// authorize API key and get session cookie
	respAuth, err := caller.Authorize(data)
	if err != nil {
		return err
	}
	if respAuth.Success == false {
		return errors.New(respAuth.Msg)
	}

	// call service.transfer/export -> returns file hashID
	respExp, err := caller.Export(opt)
	if err != nil {
		return err
	}
	if respExp.Success == false {
		return errors.New(respExp.Msg)
	}

	_, file := filepath.Split(opt.Path)
	data.Path = opt.Path
	data.FileName = file

	// call service.file/download to get the file ( opt, taskID, fileID)
	err = caller.Receive(data, respExp.Code, respExp.Msg)
	if err != nil {
		return err
	}

	return nil
}
