package ffmpeg

import (
	"errors"
	"os/exec"

	"gitlab.com/Dastanaron/ffmpeg-helper/commands"
)

type Converter struct {
	ProgressChannel chan int
	Command         string
	CommandArgs     []string
	InputFile       string
	OutputFile      string
	CmdOutput       []byte
}

func NewConverter(srcFilePath string, outputFilePath string, command commands.Command) Converter {
	command.ApplyPaths(srcFilePath, outputFilePath)
	cmd, cmdArgs := command.SplitCommand()

	return Converter{
		ProgressChannel: make(chan int),
		Command:         cmd,
		CommandArgs:     cmdArgs,
		InputFile:       srcFilePath,
		OutputFile:      outputFilePath,
	}
}

func (conv *Converter) Execute() error {
	if conv.Command == "" {
		return errors.New("command is empty")
	}
	cmd := exec.Command(conv.Command, conv.CommandArgs...)

	ProgressPipe(cmd, conv.ProgressChannel)
	return nil
}
