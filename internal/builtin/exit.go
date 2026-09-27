package builtin

import (
	"fmt"
	"os"
	"strconv"
)

func exit(args []string) error {
	s := 0
	if len(args) > 0 {
		var err error
		s, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("status code must be an integer")
		}
	}
	os.Exit(int(s))
	return nil // unreachable
}

var _ Command = exit
