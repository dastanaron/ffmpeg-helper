package main

import (
	_ "embed"
	"fmt"
	"os"

	"gitlab.com/Dastanaron/ffmpeg-helper/commands"
	"gitlab.com/Dastanaron/ffmpeg-helper/helpers"
)

const CUSTOM_COMMANDS_FILE = "./commands.yaml"

//go:embed default.commands.yaml
var defaultCommands string

func main() {
	commandsList, err := commands.ParseFromString(defaultCommands)
	helpers.CheckError("Error parsing default commands", err)

	if _, err := os.Stat(CUSTOM_COMMANDS_FILE); err == nil {
		customCommandsYaml, err := os.ReadFile(CUSTOM_COMMANDS_FILE)
		helpers.CheckError("Error open custom commands file", err)

		customCommands, err := commands.ParseFromBytes(customCommandsYaml)
		helpers.CheckError("Error parsing custom commands", err)

		*commandsList = append(*commandsList, *customCommands...)
	}

	fmt.Println(commandsList)
}
