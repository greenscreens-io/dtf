package main

import (
	"fmt"
	"greenscreens-io/dtf/install"
	"greenscreens-io/dtf/lib"
	"greenscreens-io/dtf/model"
	"greenscreens-io/dtf/transfer"
	"greenscreens-io/dtf/winapi"
	"os"
	"strings"
)

func main() {

	winapi.Guitocons()
	winapi.SetConsoleTitle("Green Screens DTF")

	i := doMain()
	os.Exit(i)
}

func doMain() int {

	var err error
	args := os.Args[1:] // without program name

	if len(args) == 0 {
		lib.PrintInstructions()
		return 10
	}
	file := args[0]
	cmd := strings.ToLower(file)

	switch cmd {
	case "install":
		return doInstall()
	case "web":
		return doWeb(args[1])
	case "auth":
		return doAuth()
	case "edit":
		file = args[1]
	case "describe":
		file = args[1]
	}

	// load transfer options
	var options model.Options
	err = options.Load(file)
	if err != nil {
		fmt.Println(err)
		return 20
	}

	fdf, _ := model.LoadFDF(file)
	if fdf != nil {
		options.Defs = *fdf
	}

	if cmd == "edit" {
		return doEdit(&options, file)
	}

	fmt.Println("Validating API Key")
	partial, err := lib.KeyCheck(options.Url, options.Key)
	if err != nil {
		fmt.Println(err)
		return 28
	}
	options.Partial = partial

	ensureOptions(&options)

	if !options.IsValid() {
		fmt.Println("Auth data not complete")
		return 30
	}

	if cmd == "describe" {
		return doDescribe(&options, file)
	}

	options.Save()

	if isImport(file) {
		return doImport(&options, file)
	}

	if isExport(file) {
		return doExport(&options, file)
	}

	fmt.Printf("Transfer file format not recognized :%s", file)
	return 5

}

func isImport(file string) bool {
	return strings.HasSuffix(file, ".gst")
}

func isExport(file string) bool {
	return strings.HasSuffix(file, ".gsf")
}

func ensureOptions(options *model.Options) {
	options.Url = lib.PromptURL(options.Url)
	options.Key = lib.PromptAlias(options.Key)
	if options.Partial {
		options.User = lib.PromptUser(options.User)
		options.Password = lib.PromptPassword(options.Password)
	} else {
		options.User = ""
		options.Password = ""
	}

}

func doInstall() int {
	err := install.Install()
	if err != nil {
		fmt.Println(err)
		return 12

	}
	return 0
}

func doWeb(data string) int {
	err := lib.WebCall(data)
	if err != nil {
		fmt.Println(err)
		fmt.Println(data)
		fmt.Println("Press enter to continue...")
		var val string
		fmt.Scanln(&val)
		return 14

	}
	return 0
}

func doAuth() int {
	err := lib.CreateCreds()
	if err != nil {
		fmt.Println(err)
		return 15

	}
	return 0
}

func doEdit(options *model.Options, file string) int {
	err := lib.EditConfig(options, file)
	if err != nil {
		fmt.Println(err)
		return 25
	}
	return 0
}

func doDescribe(options *model.Options, file string) int {
	fmt.Println("Processing transfer to server...")
	err := transfer.Describe(options, file)
	if err != nil {
		fmt.Println(err)
		return 50
	}
	return 0
}

func doImport(options *model.Options, file string) int {
	fmt.Println("Processing transfer to server...")
	err := transfer.ToGS(options)
	if err != nil {
		fmt.Println(err)
		return 40
	}
	return 0
}

// preformat parse ouput path if there are date/time format bracket string segments
func preformat(options *model.Options) {

	formated := strings.Contains(options.Path, "[")

	if !formated {
		return
	}

	val, err := lib.FileFormat(options.Url, options.Path)
	if err != nil {
		fmt.Println(err)
	} else {
		options.Path = val
	}

}

func doExport(options *model.Options, file string) int {

	fmt.Println("Processing transfer from server...")
	preformat(options)

	err := transfer.FromGS(options)
	if err != nil {
		fmt.Println(err)
		return 50
	}
	return 0
}
