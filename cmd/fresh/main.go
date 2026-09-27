package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/vncsmyrnk/fresh/internal/builtin"
)

type statusCode uint32

const (
	statusCodeSuccess statusCode = 0
)

func (s statusCode) failed() bool {
	return s > 0
}

func (s statusCode) string() string {
	return fmt.Sprintf("%d", s)
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("fresh is a simple shell implementation built for educational purposes.")
		os.Exit(0)
	}

	lastReturnStatus := statusCodeSuccess
	for {
		promptExitStatus := ""
		if lastReturnStatus.failed() {
			promptExitStatus = lastReturnStatus.string()
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

		if builtinCmd, err := builtin.Look(command); err != builtin.ErrBuiltinNotFound {
			if errBuiltinCmd := builtinCmd(arguments); errBuiltinCmd != nil {
				fmt.Fprintf(os.Stderr, "%s failed: %s\n", command, errBuiltinCmd)
				lastReturnStatus = 1
			}
			continue
		}

		cmd := exec.Command(command, arguments...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "unexpected error: %s\n", err)
			lastReturnStatus = 1
		}

		lastReturnStatus = statusCode(cmd.ProcessState.ExitCode())
	}
}
