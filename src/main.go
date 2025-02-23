package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/dastanaron/ffmpeg-helper/commands"
	"github.com/dastanaron/ffmpeg-helper/ffmpeg"
	"github.com/dastanaron/ffmpeg-helper/helpers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schollz/progressbar/v3"
)

const COMMANDS_FILE = "commands.yaml"
const APP_NAME = "ffmpeg-helper"
const APP_VERSION = "1.0.0"

//go:embed default.commands.yaml
var defaultCommands []byte

type AppConfig struct {
	InputFilePath  string
	OutputFilePath string
	CommandsList   *[]commands.Command
}

func main() {
	defaultRuntimeCommands()

	runtimeArgs := os.Args[1:]

	if len(runtimeArgs) < 2 {
		log.Fatal("Required arguments inputFile and outputFile")
	}

	inputFilePath := runtimeArgs[0]
	outputFilePath := runtimeArgs[1]

	if inputFilePath == "" || outputFilePath == "" {
		log.Fatal("Required arguments inputFile and outputFile")
	}

	commandsList, err := loadCommands()
	helpers.CheckError("Error loading commands", err)

	config := &AppConfig{
		InputFilePath:  inputFilePath,
		OutputFilePath: outputFilePath,
		CommandsList:   commandsList,
	}

	app := tview.NewApplication()
	setupUI(app, config)

	if err := app.EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}
}

func loadCommands() (*[]commands.Command, error) {
	homeDir, err := os.UserHomeDir()
	helpers.CheckError("No homedir", err)
	appConfigPath := fmt.Sprintf("%s/.config/%s", homeDir, APP_NAME)
	appCommandsFile := fmt.Sprintf("%s/%s", appConfigPath, COMMANDS_FILE)

	_, err = os.Stat(appCommandsFile)

	if err != nil {
		err = os.MkdirAll(appConfigPath, 0744)
		helpers.CheckError("Cannot create config dir", err)

		err = os.WriteFile(appCommandsFile, defaultCommands, 0644)
		helpers.CheckError("Cannot create commands file", err)
	}

	commandsYaml, err := os.ReadFile(appCommandsFile)
	helpers.CheckError("Error reading commands file", err)

	commandsList, err := commands.ParseFromBytes(commandsYaml)
	helpers.CheckError("Error parsing commands", err)

	return commandsList, nil
}

func setupUI(app *tview.Application, config *AppConfig) {
	menu := tview.NewFlex()

	list := tview.NewList()
	list.SetTitle("Select command")
	list.SetBorder(true).SetBorderColor(tcell.ColorGray)

	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetTitle("Description")
	textView.SetBorder(true).SetBorderColor(tcell.ColorGreenYellow)

	for index, v := range *config.CommandsList {
		list.AddItem(v.Name, "", 0, func() {
			app.Stop()
			runCommand(config, index)
		})
		if index == 0 {
			textView.SetText(v.Description)
		}
	}

	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	menu.AddItem(list, 0, 1, false)
	menu.AddItem(textView, 0, 1, false)

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index > len(*config.CommandsList)-1 {
			return
		}
		textView.SetText((*config.CommandsList)[index].Description)
	})

	app.SetRoot(menu, true)
	app.SetFocus(list)
}

func runCommand(config *AppConfig, index int) {
	command := (*config.CommandsList)[index]

	bar := progressbar.Default(100)
	converter := ffmpeg.NewConverter(config.InputFilePath, config.OutputFilePath, command)

	done := make(chan uint)

	go func(c chan uint, converter ffmpeg.Converter, bar *progressbar.ProgressBar) {
		for progressPercent := range converter.ProgressChannel {
			bar.Set(progressPercent)
		}
		c <- 1
	}(done, converter, bar)

	converter.Execute()

	<-done
	os.Exit(0)
}

func defaultRuntimeCommands() {
	runtimeArgs := os.Args[1:]

	if runtimeArgs[0] == "--version" {
		fmt.Println(APP_NAME, "version:", APP_VERSION)
		os.Exit(0)
	}

	if runtimeArgs[0] == "--help" {
		fmt.Println(`
      Run ./ffmpeg-helper inputFile outputFile
      Custom commands ~/.config/ffmpeg-helper/commands.yaml
      More see https://github.com/dastanaron/ffmpeg-helper
      `)
		os.Exit(0)
	}
}
