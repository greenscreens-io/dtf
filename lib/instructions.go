package lib

import (
	"fmt"
	"time"
)

func PrintInstructions() {

	t := time.Now()
	fmt.Println("")
	fmt.Println("**********************************************")
	fmt.Printf("* Copyright: Green Screens Ltd. 2016 - %d  *\n", t.Year())
	fmt.Println("* Contact: info@greenscreens.io              *")
	fmt.Println("*                                            *")
	fmt.Println("* Data Transfer Client                       *")
	fmt.Println("* Version : 1.0.0.                           *")
	fmt.Println("**********************************************")

	fmt.Println("  Call with \"install\" to register file extensions")
	fmt.Println("    dtf INSTALL")
	fmt.Println("")

	fmt.Println("  Call with \"AUTH\" to authorize user for auto login")
	fmt.Println("    dtf AUTH [SERVICE_URL] [API_KEY] [IBM_USER] [IBM_PASSWORD] [ALIAS]")
	fmt.Println("")
	fmt.Println("  For console entry use")
	fmt.Println("   dtf AUTH")
	fmt.Println("")

	fmt.Println("  Call with path to *.gsf or *.gst to transfer data")
	fmt.Println("")
	fmt.Println("    dtf data.gsf")
	fmt.Println("")

	fmt.Println("  Optionally, provide user and password for IBM i server if not stored inside transfer file")
	fmt.Println("    dtf data.gst -u IBMUSER -p IBMPASSWORD")
	fmt.Println("")

	fmt.Println("  Call with \"EDIT\" to open configuration in Web UI console. ")
	fmt.Println("    dtf EDIT data.gst")
	fmt.Println("")

}
