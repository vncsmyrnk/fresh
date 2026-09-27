package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"

	fbuiltin "github.com/vncsmyrnk/fresh/internal/builtin"
	fproc "github.com/vncsmyrnk/fresh/internal/proc"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("fresh is a simple shell implementation built for educational purposes.")
		os.Exit(0)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	onExit := func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			fmt.Fprintln(os.Stderr, "fresh: failed to restore terminal state.")
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT|syscall.SIGKILL)
	go func() {
		<-sigChan
		onExit()
		os.Exit(1)
	}()

	t := term.NewTerminal(os.Stdin, "")
	lastReturnStatus := fproc.StatusCodeSuccess

	for {
		promptStatusCode := ""
		if lastReturnStatus.Failed() {
			promptStatusCode = lastReturnStatus.String()
		}
		prompt := fmt.Sprintf("%s> ", promptStatusCode)
		t.SetPrompt(prompt)

		l, err := t.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, "fresh: failed to read input.")
				os.Exit(1)
			}
			t = term.NewTerminal(os.Stdin, prompt)
			fmt.Println("\r")
			continue
		}

		if l == "" {
			continue
		}

		tokens := strings.Split(l, " ")

		command := tokens[0]
		arguments := tokens[1:]

		if command == "exit" {
			if fbuiltin.IsExit(command) {
				switch len(arguments) {
				case 0:
					onExit()
					os.Exit(int(lastReturnStatus))
				case 1:
					s, err := strconv.Atoi(arguments[0])
					if err != nil {
						fmt.Fprintln(os.Stderr, "fresh: exit expects one or no arguments.")
						continue
					}
					onExit()
					os.Exit(s)
				default:
					fmt.Fprintln(os.Stderr, "fresh: exit expects one or no arguments.")
				}
			}
		}

		if builtinCmd, err := fbuiltin.Lookup(command); err != fbuiltin.ErrBuiltinNotFound {
			if errBuiltinCmd := builtinCmd(arguments); errBuiltinCmd != nil {
				fmt.Fprintf(os.Stderr, "fresh: %s failed: %s\n", command, errBuiltinCmd)
				lastReturnStatus = 1
			}
			continue
		}

		cmd := exec.Command(command, arguments...)
		cmd.Stdin = os.Stdin
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

		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			fmt.Fprintln(os.Stderr, "fresh: failed to restore terminal state.")
		}

		if err := cmd.Run(); err != nil && cmd.Process == nil {
			fmt.Fprintf(os.Stderr, "fresh: unexpected error: %s\n", err)
			lastReturnStatus = 1
		}

		_, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "fresh: failed to enter raw mode: %s\n", err)
			lastReturnStatus = 1
		}

		lastReturnStatus = fproc.StatusCode(cmd.ProcessState.ExitCode())
		sigChan <- nil
		fmt.Print("\r")
	}
}
