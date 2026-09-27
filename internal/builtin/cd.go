package builtin

import (
	"fmt"
	"os"
)

func cd(args []string) error {
	if len(args) == 0 || len(args) > 1 {
		return fmt.Errorf("you must specify one path")
	}
	return os.Chdir(args[0])
}

var _ Command = cd
