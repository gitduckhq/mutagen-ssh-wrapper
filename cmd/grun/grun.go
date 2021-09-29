package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/sync/errgroup"
	"ssh-wrapper/cmd"
)

func main() {
	if len(os.Args) < 2 {
		if len(os.Args) == 1 {
			fmt.Fprintln(os.Stderr, "Error: missing arguments: parent PID and command")
		} else {
			fmt.Fprintln(os.Stderr, "Error: missing command argument")
		}
		os.Exit(1)
	}

	commandName, args := os.Args[2], os.Args[3:]

	parentPid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: parent PID is not a number:", parentPid)
		os.Exit(1)
	}

	parentProcess, err := os.FindProcess(parentPid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot find process PID %d: %v", parentPid, err)
		os.Exit(1)
	}

	g, ctx := errgroup.WithContext(context.Background())

	command := exec.CommandContext(ctx, commandName, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Wait for the child process to finish
	g.Go(func() error {
		cmd.WrapRunAndExit(command)
		return nil
	})

	// Wait for the parent process to finish. If so, kill the child process
	g.Go(func() error {
		_, err := parentProcess.Wait()
		command.Process.Kill()
		return err
	})

	g.Wait()
	os.Exit(1)
}
