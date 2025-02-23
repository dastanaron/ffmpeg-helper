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

	commandsList, err := commands.ParseFromString(defaultCommands)
	helpers.CheckError("Error parsing default commands", err)

	if _, err := os.Stat(CUSTOM_COMMANDS_FILE); err == nil {
		customCommandsYaml, err := os.ReadFile(CUSTOM_COMMANDS_FILE)
		helpers.CheckError("Error open custom commands file", err)

		customCommands, err := commands.ParseFromBytes(customCommandsYaml)
		helpers.CheckError("Error parsing custom commands", err)

		*commandsList = append(*commandsList, *customCommands...)
	}

	app := tview.NewApplication()

	menu := tview.NewFlex()

	list := tview.NewList()
	list.SetTitle("Select command")
	list.SetBorder(true).SetBorderColor(tcell.ColorGray)

	menu.AddItem(list, 0, 1, false)

	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetTitle("Description")
	textView.SetBorder(true).SetBorderColor(tcell.ColorGreenYellow)

	for i, v := range *commandsList {
		list.AddItem(v.Name, "", 0, func() {
			app.Stop()

			command := (*commandsList)[list.GetCurrentItem()]

			bar := progressbar.Default(100)

			converter := ffmpeg.NewConverter(inputFilePath, outputFilePath, command)

			c := make(chan uint)

			go func(c chan uint, converter ffmpeg.Converter, bar *progressbar.ProgressBar) {
				for progressPercent := range converter.ProgressChannel {
					bar.Set(progressPercent)
				}

				c <- 1
			}(c, converter, bar)

			converter.Execute()

			end := <-c

			if end == 1 {
				os.Exit(0)
			}

		})
		if i == 0 {
			textView.SetText(v.Description)
		}
	}

	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	menu.AddItem(textView, 0, 1, false)

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index > len(*commandsList)-1 {
			return
		}

		description := (*commandsList)[index].Description
		textView.SetText(description)
	})

	if err := app.SetRoot(menu, true).SetFocus(list).Run(); err != nil {
		log.Fatal(err)
	}
}
