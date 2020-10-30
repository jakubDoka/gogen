package main

import (
	"gogen/commands"
)

func main() {
	commands.LoadConfig()

	commands.NewHandler()

	commands.RunCommand()

	commands.SaveConfig()
}
