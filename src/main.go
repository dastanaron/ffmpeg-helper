package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schollz/progressbar/v3"
	"gitlab.com/Dastanaron/ffmpeg-helper/commands"
	"gitlab.com/Dastanaron/ffmpeg-helper/ffmpeg"
	"gitlab.com/Dastanaron/ffmpeg-helper/helpers"
)

const CUSTOM_COMMANDS_FILE = "./commands.yaml"

//go:embed default.commands.yaml
var defaultCommands string

type AppConfig struct {
	InputFilePath  string
	OutputFilePath string
	CommandsList   *[]commands.Command
}

func main() {
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
	commandsList, err := commands.ParseFromString(defaultCommands)
	helpers.CheckError("Error parsing default commands", err)

	if _, err := os.Stat(CUSTOM_COMMANDS_FILE); err == nil {
		customCommandsYaml, err := os.ReadFile(CUSTOM_COMMANDS_FILE)
		helpers.CheckError("Error reading custom commands file", err)

		customCommands, err := commands.ParseFromBytes(customCommandsYaml)
		helpers.CheckError("Error parsing custom commands", err)

		*commandsList = append(*commandsList, *customCommands...)
	}

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
