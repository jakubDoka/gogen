package main

import (
	"gogen/commands"
)

func main() {
	name := commands.ParseArgs()
	commands.LoadConfig()

	commands.NewHandler()

	commands.RunCommand(name)

	commands.SaveConfig()
}
