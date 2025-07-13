package command

import "github.com/yamaki-87/mailbot/src/consts"

type Command interface {
	Execute(content string) (string, error)
}

var commandMap = map[string]Command{
	consts.HELPCOMMAND: &HelpCommand{},
}

func HandleCommand(input string) (string, error) {
	if cmd, ok := commandMap[input]; ok {
		return cmd.Execute(input)
	}

	return "", nil
}
