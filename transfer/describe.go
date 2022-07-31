/**
 * This file represents main library functions called from main module
 * Used for importing data to IBM i through GS server
 */
package transfer

import (
	"errors"
	"greenscreens-io/dtf/model"
	"greenscreens-io/dtf/services"
	"os"
)

// Describe file structure for importing new tables
func Describe(opt *model.Options, file string) error {

	var caller services.HTTPCaller

	caller.Init()

	// streaming file SQL/MBR import is not available
	opt.Streaming = false
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

	respExp, err := caller.Describe(data, opt.Library, opt.Object)
	if err != nil {
		return err
	}

	if !respExp.Success {
		return errors.New(respExp.Msg)
	}

	tmp, err := respExp.ToRawData()
	if err != nil {
		return err
	}

	err = os.WriteFile(file+"fd", tmp, 0644)
	if err != nil {
		return err
	}

	return nil
}
