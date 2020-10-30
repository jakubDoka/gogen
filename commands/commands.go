package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//gen Worker<string, struct{}, 30, R>
//gen Worker<string, []Template, 30, T>

// CommandHandler global state, we need just one so let it be global
var CommandHandler Handler

// Command is encapsulation for any command line variation
type Command struct {
	help                          []string
	name, description, argsStruct string
	run                           func()
}

// GetRestriction returns minimal and maximal amount of arguments you can pass to command
func (c *Command) GetRestriction() (min, max int) {
	for _, a := range strings.Split(c.argsStruct, " ") {
		if a[0] == '<' {
			min++
			max++
		} else if a[0] == '[' {
			max++
		}
	}

	return
}

// InitCommand Adds repetition to command description
func InitCommand(c Command) Command {
	c.description += ",for more info use command 'gogen help " + c.name + "'"
	return c
}

// Handler handles commands
type Handler struct {
	commands map[string]Command
}

//RegisterCommand ...
func (h *Handler) RegisterCommand(c Command) {
	h.commands[c.name] = InitCommand(c)
}

// RunCommand ...
func (h *Handler) RunCommand(name string) {
	command, ok := h.commands[name]
	if !ok {
		Terminate("Sorry i don't know this command, try 'gogen help'")
	}
	min, max := command.GetRestriction()

	if len(Args) < min {
		Terminate("too few arguments, argument structure: " + command.argsStruct)
	} else if len(Args) > max {
		Terminate("Too match arguments, argument structure: " + command.argsStruct)
	}

	command.run()
}

// NewHandler registers all commands
func NewHandler() {
	handler := Handler{commands: map[string]Command{}}

	handler.RegisterCommand(Command{
		name:        "help",
		description: "it does just what you see now, and also shows help about other commands",
		help: []string{
			"- shows all commands and descriptions",
			"<commandName> - shows more descriptions about command",
		},
		argsStruct: "[commandName]",
		run: func() {
			if len(Args) == 0 {
				PrintAllCommands()
			}

			command, ok := CommandHandler.commands[Args[0]]
			if !ok {
				Terminate("Only think i know is that COMMAND DOES NOT EXIST!")
			}
			PrintCommandHelp(command)
		},
	})

	handler.RegisterCommand(Command{
		name:        "add",
		description: "adds current or specified directory to template paths, for more info use 'help add' ",
		help: []string{
			"- adds current working directory to templates",
			"<dirName> - adds desired directory to templates",
			"<dirName> -rm - removes an directory, omitting dir works too",
			"info - shows all registered directories",
			"clear - deletes all directories",
		},
		argsStruct: "[dirName/info/clear] -rm",
		run: func() {
			var dir string

			if len(Args) == 0 {
				dir = WDir
			} else {
				dir = Args[0]
			}

			if dir == "clear" {
				Confirm()
				Cf.TemplateDirectories = map[string]bool{}
				Cf.TemplateFiles = map[string]bool{}
				fmt.Println("Big boy big CLEANUP is DONE")
				return
			} else if dir == "info" {
				fmt.Println("Template paths:")
				for k := range Cf.TemplateDirectories {
					fmt.Printf("\t%s\n", k)
				}
				Exit("Thats all i have so far.")
			} else if !Exists(dir) {
				Terminate("How em i supposed to add DIRECTORY that DOES NOT EXIST?!")
			} else if !IsAccessable(dir) {
				Terminate("Ops, DIRECTORY you gave me CANNOT BE ACCESSED!")
			}

			stats, err := os.Stat(dir)
			CheckError("unable to determinate whether path points to directory or file", err)

			dir, _ = filepath.Abs(dir)

			isDir := stats.IsDir()
			var ok bool

			if isDir {
				_, ok = Cf.TemplateDirectories[dir]
			} else {
				_, ok = Cf.TemplateFiles[dir]
			}

			if Labels["-rm"] {
				if !ok {
					Terminate("No such directory/file in the list.")
				}
				if isDir {
					delete(Cf.TemplateDirectories, dir)
				} else {
					delete(Cf.TemplateFiles, dir)
				}
				fmt.Printf("%s removed from templates", dir)
			} else {
				if ok {
					Terminate("I already have this directory/file in list.")
				}

				if d, ok := IsIncluded(dir, Cf.TemplateDirectories); ok {
					Terminate("this file is already included by " + d)
				}

				if isDir {
					FilterIncluded(dir, Cf.TemplateDirectories)
					FilterIncluded(dir, Cf.TemplateFiles)

					Cf.TemplateDirectories[dir] = true
				} else {

					Cf.TemplateFiles[dir] = true
				}

				fmt.Printf("%s added to templates", dir)
			}
		},
	})

	handler.RegisterCommand(Command{
		name:        "gen",
		description: "updates all templates you annotated and creates new ones if needed",
		help: []string{
			"- updates templates in current directory, command creates new file for templates",
			"-r - updates all files recursively from working directory",
			"<dirName> - updates templates in desired directory, -r also works",
		},
		argsStruct: "[dirName] -r",
		run: func() {
			fmt.Println("Parsing templates...")
			rec := make(chan []Template, 0)
			cores := runtime.NumCPU()
			threads := make([]chan string, cores)
			for i := range threads {
				threads[i] = TWorker(rec, func(dir string) []Template {
					return ParseTemplatesInDir(dir)
				}, func(str string) bool {
					return str == ""
				}, false)
			}
			dirs := []string{}
			for k := range Cf.TemplateDirectories {
				dirs = append(dirs, GetDirList(k, Cf.TemplateDirMaxDepth)...)
			}

			for i, dir := range dirs {
				threads[i%cores] <- dir
			}

			for f := range Cf.TemplateFiles {
				a := ParseTemplateFile(f, false)
				for i, v := range a {
					Templates[v.name] = &a[i]
				}
			}

			for range dirs {
				a := <-rec
				for i, v := range a {
					if _, ok := Templates[v.name]; ok && v.err == -1 {
						v.err = Duplicate
					}
					HandleSyntaxError(v.ErrData)
					Templates[v.name] = &a[i]
				}
			}

			for _, t := range threads {
				t <- ""
			}

			fmt.Println("Generating code...")

			if len(Args) == 1 {
				dirs = GetDirList(Args[0], Cf.WorkingDirMaxDepth)
			} else {
				dirs = GetDirList(WDir, Cf.WorkingDirMaxDepth)
			}

			if Labels["-r"] {
				rec := make(chan struct{}, 0)
				for i := range threads {
					threads[i] = RWorker(rec, func(dir string) struct{} {
						CreateTemplatesInDir(dir)
						return struct{}{}
					}, func(str string) bool {
						return str == ""
					}, false)
				}
				for i, d := range dirs {
					threads[i%cores] <- d
				}
				for _, t := range threads {
					t <- ""
				}

				for _, d := range dirs {
					CreateTemplatesInDir(d)
				}
			} else {
				CreateTemplatesInDir(WDir)
			}

			fmt.Println("Finished.")
		},
	})

	CommandHandler = handler
}
