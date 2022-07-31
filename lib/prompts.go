package lib

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"
)

var MATCH_URL = getExpression(`^https?:\/\/`)
var MATCH_KEY = getExpression(`\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func getExpression(rule string) *regexp.Regexp {
	rxp, _ := regexp.Compile(rule)
	return rxp
}

// Prompt is generic funnction to read from console input
func Prompt(msg string, loop, hidden bool) string {

	var val string

	for {
		fmt.Println(msg)
		if hidden {
			val = GetPassword()
		} else {
			fmt.Scanln(&val)
		}
		if !loop {
			break
		}
		if len(val) > 0 {
			break
		}
	}

	return val
}

// PromptURL ask for URL entry and validate enterd url
func PromptURL(url string) string {

	matched := false
	value := url

	if !MATCH_URL.Match([]byte(value)) {
		for {
			fmt.Println("Invalid url! Must start with http:// or https://")
			value = Prompt("Enter Service URL:", true, false)
			value = strings.ToLower(value)
			matched = MATCH_URL.Match([]byte(value))
			if matched {
				break
			}
		}
	}

	return strings.ToLower(value)
}

// PromptKey ask for UUIDv4 format entry for API key
func PromptKey(key string) string {

	matched := false
	value := key

	if !MATCH_KEY.Match([]byte(value)) {
		for {
			fmt.Println("Invalid API Key format!")
			value = Prompt("Enter API Key:", true, false)
			matched = MATCH_KEY.Match([]byte(value))
			if matched {
				break
			}
		}
	}

	return value
}

// PromptAlias ask for API key alias
func PromptAlias(key string) string {

	matched := false
	value := key

	if len(value) == 0 {
		for {
			fmt.Println("Invalid API Key format!")
			value = Prompt("Enter API Key:", true, false)
			matched = len(value) > 0
			if matched {
				break
			}
		}
	}

	return value
}

// PromptUser ask for IBM i user, max length is 10
func PromptUser(user string) string {

	matched := false
	value := user

	if len(value) == 0 {
		for {
			fmt.Println("Invalid IBM user! Maximum length is 10 characters.")
			value = Prompt("Enter IBM user name:", true, false)
			matched = len(value) < 11
			if matched {
				break
			}
		}
	}

	return strings.ToUpper(value)

}

// PromptPassword ask for password entry, input is hidden
func PromptPassword(password string) string {

	matched := false
	value := password

	if len(value) == 0 {
		for {
			fmt.Println("Invalid password! Not set")
			value = Prompt("Enter user password", true, true)
			matched = len(value) > 0
			if matched {
				break
			}
		}
	}

	return value

}

// GetPassword read password from console input
func GetPassword() string {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	t := term.NewTerminal(os.Stdin, ">")
	pass, err := t.ReadPassword(">")
	if err != nil {
		panic(err)
	}
	return pass
}
