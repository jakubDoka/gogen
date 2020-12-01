package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// ValidFormat it is printed on multiple places and it may change
const ValidFormat = "TemplateName<TemplateArguments..., Prefix>"

// CheckError exits and prints error message if there is an error
func CheckError(message string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("ERROR: %s \nMATTER: %s \n", message, err)
	os.Exit(-1)
}

// Terminate program on user triggered error
func Terminate(message string) {
	fmt.Printf("ERROR: %s \n", message)
	os.Exit(1)
}

// IsAccessable returns whethere file is accessable
func IsAccessable(path string) bool {
	stats, err := os.Stat(path)
	return err != nil || stats.Mode().Perm()&(1<<(uint(7))) != 0
}

// Exit exits with success message
func Exit(message string) {
	fmt.Println(message)
	os.Exit(0)
}

// Warming displays warming if error is not nil
func Warming(message string, err error) bool {
	if err == nil {
		return false
	}
	fmt.Printf("WARMING: %s\nMATTER: %s\n", message, err)
	return true
}

// Confirm gives option to user to confirm his dessision
func Confirm() {
	fmt.Println("Please confirm (y/n):")
	var input string
	fmt.Scanln(&input)

	if len(input) == 0 || input[:1] != "y" {
		Exit("Newer mind...")
	}
}

// GetErrorLine finds line where error happened
func GetErrorLine(path, name, content string) int {
	count := strings.Count(content, name)
	if count > 2 {
		count = 2
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	line := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, name) && strings.Contains(text, Cf.TemplateStart) {
			count--
		}
		if count == 0 {
			break
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return line
}

// SyntaxErr ...
type SyntaxErr int

// Syntax err variants
const (
	IncorrectArgs SyntaxErr = iota
	IncorrectRules
	IncorrectNamePrefix
	MissingTemplateArgs
	Duplicate
	Undefined
)

// HandleSyntaxError prints relevant error message, its like this just to have them all together
func HandleSyntaxError(adr ErrData) {
	er := func(str string) { SyntaxError(str, adr.path, adr.line) }
	switch adr.err {
	case MissingTemplateArgs:
		er("missing template arguments, template is redundant")
	case IncorrectRules:
		er("the template rules signature is incorrect, example signature: " + ValidFormat)
	case IncorrectArgs:
		er("you supplied incorrect amount of arguments")
	case Duplicate:
		er("duplicate template")
	case Undefined:
		er("template is undefined")
	case IncorrectNamePrefix:
		er("Name prefix is incorrect, correct syntax: " + ValidFormat)
	}
}

// SyntaxError formats error and prints it
func SyntaxError(message, path string, line int) {
	Terminate(fmt.Sprintf("%s : %d >> ", path, line) + message)
}
