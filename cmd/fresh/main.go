package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	fbuiltin "github.com/vncsmyrnk/fresh/internal/builtin"
	fproc "github.com/vncsmyrnk/fresh/internal/proc"
)

const (
	promptInitialSizeBytes = 1
	promptLimitSizeBytes   = 1024
	promptReallocFactor    = 2
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

		i := 1
		bPrompt := make([]byte, 0, promptLimitSizeBytes)
		for {
			b := make([]byte, i)
			n, err := os.Stdin.Read(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fresh: failed to read input.")
				continue
			} else if b[len(b)-1] == byte(10) || n < len(b) {
				bPrompt = append(bPrompt, b[:n]...)
				break
			}
			i *= promptReallocFactor
			if i > promptLimitSizeBytes {
				fmt.Fprintf(os.Stderr, "fresh: prompt size exceeded.")
				continue
			}
			bPrompt = append(bPrompt, b...)
		}

		prompt := string(bPrompt)
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

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)
		go func() {
			<-sigChan
			if cmd.Process != nil && cmd.ProcessState == nil {
				if err := cmd.Process.Signal(os.Interrupt); err != nil {
					fmt.Fprintf(os.Stderr, "fresh: process interruption failed: %s", err)
				}
				fmt.Println()
			}
		}()

		if err := cmd.Run(); err != nil && cmd.Process == nil {
			fmt.Fprintf(os.Stderr, "fresh: unexpected error: %s\n", err)
			lastReturnStatus = 1
		}

		lastReturnStatus = fproc.StatusCode(cmd.ProcessState.ExitCode())
		sigChan <- nil
	}
}
