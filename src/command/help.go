package command

import (
	"os"

	"github.com/yamaki-87/mailbot/src/config"
)

type HelpCommand struct{}

func (h *HelpCommand) Execute(content string) (string, error) {
	helpMessagePath := config.GetConfig().MessageTmpl.Help
	helpContent, err := os.ReadFile(helpMessagePath)

	if err != nil {
		return "", err
	}

	return string(helpContent), nil
}
