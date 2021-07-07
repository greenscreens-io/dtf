package lib

import (
	"os"
	"log"
	"strings"
	"syscall"
	"io/ioutil"
	"path/filepath"
	"greenscreens-io/dtf/winapi"
)

const template = `Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\GreenScreens.DTF]
@="Data Transfer Definition File"

[HKEY_CLASSES_ROOT\GreenScreens.DTF\shell]

[HKEY_CLASSES_ROOT\GreenScreens.DTF\shell\open]

[HKEY_CLASSES_ROOT\GreenScreens.DTF\shell\open\command]
@="#PATH#\\dtf.exe \"%1\""

[HKEY_CLASSES_ROOT\.gsf]
"Content Type"="text/plain"
"PerceivedType"="text"
@="GreenScreens.DTF"

[HKEY_CLASSES_ROOT\.gst]
"Content Type"="text/plain"
"PerceivedType"="text"
@="GreenScreens.DTF"

`

// Install file extensions for this program
func Install() {

	dir := getDir()

	regFile := filepath.Join(dir, "dtf.reg");
	//template := loadReg()
	generate(template, regFile, "#PATH#", dir)

	err := winapi.ShellExecuteNowait(0, "", "regedit.exe", "/s \""+regFile+"\"", "", syscall.SW_HIDE)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(regFile)
}

// get current program dir
func getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// generate files
func generate(data, to, key, value string) {
	path := strings.Replace(value, "\\", "\\\\", -1)
	text := strings.Replace(string(data), key, path, -1)
	ioutil.WriteFile(to, []byte(text), 0644)
}
