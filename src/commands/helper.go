package commands

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Command struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	CMD         string `yaml:"cmd"`
}

type CommandYaml struct {
	Commands []Command `yaml:"commands"`
}

func ParseFromString(yamlData string) (*[]Command, error) {
	return ParseFromBytes([]byte(yamlData))
}

func ParseFromBytes(yamlData []byte) (*[]Command, error) {
	var data CommandYaml

	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}

	return &data.Commands, nil
}

func (cmd *Command) ApplyPaths(inputFile, outputFile string) *Command {
	command := strings.ReplaceAll(cmd.CMD, "{if}", inputFile)
	command = strings.ReplaceAll(command, "{of}", outputFile)
	cmd.CMD = command
	return cmd
}

func (cmd *Command) SplitCommand() (string, []string) {
	return SplitCommand(cmd.CMD)
}

func SplitCommand(cmd string) (string, []string) {
	var args []string
	var buffer string
	var inQuotes bool

	for _, char := range cmd {
		switch char {
		case ' ':
			if !inQuotes {
				if buffer != "" {
					args = append(args, buffer)
					buffer = ""
				}
			} else {
				buffer += string(char)
			}
		case '"':
			inQuotes = !inQuotes
		default:
			buffer += string(char)
		}
	}

	if buffer != "" {
		args = append(args, buffer)
	}

	return args[0], args[1:]
}
