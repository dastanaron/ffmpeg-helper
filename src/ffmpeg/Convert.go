package ffmpeg

import (
	"errors"
	"os/exec"
)

type Converter struct {
	ProgressChannel chan int
	Command         string
	InputFile       string
	OutputFile      string
	CmdOutput       []byte
}

type Command struct {
	Name string
	Cmd  string
}

func NewConverter(srcFilePath string, outputFilePath string) Converter {
	return Converter{
		ProgressChannel: make(chan int),
		Command:         "",
		InputFile:       srcFilePath,
		OutputFile:      outputFilePath,
	}
}

func (conv *Converter) SelectCommand(name string, commands []Command) Converter {
	command := ""
	for _, v := range commands {
		if v.Name == name {
			command = v.Cmd
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
