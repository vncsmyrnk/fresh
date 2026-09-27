package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	fbuiltin "github.com/vncsmyrnk/fresh/internal/builtin"
	fproc "github.com/vncsmyrnk/fresh/internal/proc"
)

const (
	statusCodeSuccess statusCode = 0
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("fresh is a simple shell implementation built for educational purposes.")
		os.Exit(0)
	}

	lastReturnStatus := fproc.StatusCodeSuccess
	for {
		promptExitStatus := ""
		if lastReturnStatus.Failed() {
			promptExitStatus = lastReturnStatus.String()
		}
		fmt.Printf("%s> ", promptExitStatus)

		b := make([]byte, 256)
		_, _ = os.Stdin.Read(b)

		prompt := string(b)
		promptTrimmed := strings.Split(prompt, "\n")[0]
		tokens := strings.Split(promptTrimmed, " ")

		command := tokens[0]
		arguments := tokens[1:]

		if command == "" {
			continue
		}

		if builtinCmd, err := fbuiltin.Look(command); err != fbuiltin.ErrBuiltinNotFound {
			if errBuiltinCmd := builtinCmd(arguments); errBuiltinCmd != nil {
				fmt.Fprintf(os.Stderr, "fresh: %s failed: %s\n", command, errBuiltinCmd)
				lastReturnStatus = 1
			}
			continue
		}

		cmd := exec.Command(command, arguments...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil && cmd.Process == nil {
			fmt.Fprintf(os.Stderr, "fresh: unexpected error: %s\n", err)
			lastReturnStatus = 1
		}

		lastReturnStatus = fproc.StatusCode(cmd.ProcessState.ExitCode())
	}
}
