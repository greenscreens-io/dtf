package main

import (
	"fmt"
	lib "greenscreens-io/dtf/lib"
	"os"
	"strings"
)

func print() {

	fmt.Println("**********************************************")
	fmt.Println("* Copyright: Green Screens Ltd. 2016 - 2021  *")
	fmt.Println("* Contact: info@greenscreens.io              *")
	fmt.Println("*                                            *")
	fmt.Println("* Data Transfer Client                       *")
	fmt.Println("* Version : 1.0.0.                           *")
	fmt.Println("**********************************************")

	fmt.Println("  Call with path to *.gsf or *.gst to transfer data")
	fmt.Println("    dtf data.gsf")
	fmt.Println("")

	fmt.Println("  Call with \"install\" to register file extensions")
	fmt.Println("    dtf install")
}

func main() {

	/*
	lib.TRACE = true
	test3()
	*/

	i := doMain()
	os.Exit(i)
}

func doMain() int {

	var err error
	args := os.Args[1:] // without program name

	if len(args) < 1 {
		print()
		return 1
	}

	if args[0] == "install" {
		lib.Install()
		return 0
	}

	// load transfer options
	var options lib.Options
	err = options.Load(args[0])
	if err != nil {
		fmt.Println(err)
		return 2
	}

	if len(args) >= 3 {
		if args[1] == "-u" {
			options.User = args[2]
		}
		if args[1] == "-p" {
			options.Password = args[2]
		}
	}

	if len(args) >= 5 {
		if args[3] == "-u" {
			options.User = args[4]
		}
		if args[3] == "-p" {
			options.Password = args[4]
		}
	}

	// if it is upload type
	sfx := strings.HasSuffix(args[0], ".gst")
	if sfx {
		err = lib.ToGS(&options)
		if err != nil {
			fmt.Println(err)
			return 3
		}
		return 0
	}

	// if it is download type
	sfx = strings.HasSuffix(args[0], ".gsf")
	if sfx {
		err = lib.FromGS(&options)
		if err != nil {
			fmt.Println(err)
			return 4
		}
		return 0
	}

	print()
	return 0
}
