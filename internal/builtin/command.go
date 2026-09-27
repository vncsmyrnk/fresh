package builtin

import "fmt"

type Command func([]string) error

var ErrBuiltinNotFound = fmt.Errorf("builtin not found")

var commands = map[string]Command{
	"cd":   cd,
	"exit": exit,
}

func Lookup(cmd string) (Command, error) {
	if c, ok := commands[cmd]; ok {
		return c, nil
	}
	return nil, ErrBuiltinNotFound
}

func IsExit(cmd string) bool {
	return cmd == "exit"
}
