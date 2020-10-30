package commands

import (
	"os"
	"path/filepath"
)

//Templates stores all templates
var Templates = map[string]*Template{}

// Dirs is list of dirs in working directory
var Dirs []string

// Labels Stores all labels inputted as arguments like -r and so on
var Labels map[string]bool

// Args are all ordered arguments, no labels
var Args []string

// RunCommand ...
func RunCommand(arg string) {

	CommandHandler.RunCommand(arg)
}

// LoadDirList ...
func LoadDirList() {
	dir, err := os.Getwd()
	CheckError("cannot access working directory", err)
	Dirs = GetDirList(dir)
}

// ParseArgs sorts args to labels, arguments, and other
func ParseArgs() string {
	if len(os.Args) == 1 {
		Terminate("Hello there! I em Golang code generator, but in short just gogen. Sey 'gogen help' to see what i offer.")
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
func GetDirList(path string) (dirs []string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return err
	})

	return dirs
}
