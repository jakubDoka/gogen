package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//Templates stores all templates
var Templates = map[string]*Template{}

// Dirs is list of dirs in working directory
var Dirs []string

var WDir string

// Labels Stores all labels inputted as arguments like -r and so on
var Labels map[string]bool

// Args are all ordered arguments, no labels
var Args []string

// RunCommand ...
func RunCommand(arg string) {

	CommandHandler.RunCommand(arg)
}

// LoadWorkingDirectory ...
func LoadWorkingDirectory() {
	var err error
	WDir, err = os.Getwd()
	CheckError("cannot access working directory", err)
}

// ParseArgs sorts args to labels, arguments, and other
func ParseArgs() string {
	if len(os.Args) == 1 {
		Terminate("Hello there! Sey 'gogen help' to see what i offer.")
	}
	Labels = map[string]bool{}
	for _, arg := range os.Args[2:] {

		if len(arg) > 1 && arg[:2] == "--" {
			// TODO add these
			continue
		} else if len(arg) > 0 && arg[:1] == "-" {
			Labels[arg] = true
			continue
		}
		Args = append(Args, arg)
	}

	return os.Args[1]
}

// GetDirList lists all directories in range
func GetDirList(path string, limit int) (dirs []string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
			if len(dirs) == limit {
				Warming("not continuing to the deeper directories",
					fmt.Errorf("exceeded limit of %d, if you want to go deeper, set it in config"))
				return errors.New("going too deep")
			}
		}
		return err
	})

	return dirs
}
