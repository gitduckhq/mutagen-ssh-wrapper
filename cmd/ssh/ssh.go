package main

import (
	"fmt"
	"os"

	"ssh-wrapper/cmd"
	"ssh-wrapper/pkg/ssh"
)

func main() {
	args := cmd.ProcessArgs()

	commandName, err := ssh.SshCommandPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: unable to set up SSH invocation:", err)
		os.Exit(1)
	}

	cmd.Grun(commandName, args...)
}
