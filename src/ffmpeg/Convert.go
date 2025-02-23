package ffmpeg

import (
	"errors"
	"os/exec"

	"gitlab.com/Dastanaron/ffmpeg-helper/commands"
)

type Converter struct {
	ProgressChannel chan int
	Command         string
	InputFile       string
	OutputFile      string
	CmdOutput       []byte
}

func NewConverter(srcFilePath string, outputFilePath string) Converter {
	return Converter{
		ProgressChannel: make(chan int),
		Command:         "",
		InputFile:       srcFilePath,
		OutputFile:      outputFilePath,
	}
}

func (conv *Converter) SelectCommand(name string, commands []commands.Command) Converter {
	command := ""
	for _, v := range commands {
		if v.Name == name {
			command = v.CMD
		}
	}

	conv.Command = command

	return *conv
}

func (conv *Converter) Execute() error {
	if conv.Command == "" {
		return errors.New("command is empty")
	}
	cmd := exec.Command("", conv.Command)

	ProgressPipe(cmd, conv.ProgressChannel)
	return nil
}
