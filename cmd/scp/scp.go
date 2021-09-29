package main

import (
	"fmt"
	"os"

	"ssh-wrapper/pkg/ssh"
	"ssh-wrapper/cmd"
)

func main() {
	args := cmd.ProcessArgs()

	commandName, err := ssh.ScpCommandPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: unable to set up SCP invocation:", err)
		os.Exit(1)
	}

	cmd.Grun(commandName, args...)
}
