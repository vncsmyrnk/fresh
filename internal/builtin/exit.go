package builtin

func exit(args []string) error {
	return nil
}

var _ Command = exit
