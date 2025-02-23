package commands

import "gopkg.in/yaml.v3"

type Command struct {
	Name string `yaml:"name"`
	CMD  string `yaml:"cmd"`
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
