package builtin

import "fmt"

type Command func([]string) error

var ErrBuiltinNotFound = fmt.Errorf("builtin not found")

var commands = map[string]Command{
	"cd": cd,
}

func Look(cmd string) (Command, error) {
	if c, ok := commands[cmd]; ok {
		return c, nil
	}
	return nil, ErrBuiltinNotFound
}
