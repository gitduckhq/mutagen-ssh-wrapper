package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"ssh-wrapper/pkg/process"
)

const (
	MUTAGEN_SSH_CONFIG                 = "MUTAGEN_SSH_CONFIG"
	MUTAGEN_CONNECT_TIMEOUT_IN_SECONDS = 20
)

func ProcessArgs() []string {
	sshConfig := os.Getenv(MUTAGEN_SSH_CONFIG)
	if sshConfig == "" {
		fmt.Fprintln(os.Stderr, "Error: MUTAGEN_SSH_CONFIG environment variable is not set")
		os.Exit(1)
	}

	// Increase the default timeout
	args := os.Args[1:]
	for i, arg := range args {
		if strings.HasPrefix(arg, "-oConnectTimeout=") {
			args[i] = fmt.Sprintf("-oConnectTimeout=%d", MUTAGEN_CONNECT_TIMEOUT_IN_SECONDS)
			break
		}
	}

	return append([]string{fmt.Sprintf("-F%s", sshConfig)}, args...)
}

func WrapRunAndExit(command *exec.Cmd) {
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func Grun(commandName string, args ...string) {
	args = append([]string{strconv.Itoa(os.Getpid()), commandName}, args...)

	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: cannot get current executable path")
		os.Exit(1)
	}

	grun := path.Join(filepath.Dir(self), process.ExecutableName("grun", runtime.GOOS))
	command := exec.CommandContext(context.Background(), grun, args...)
	WrapRunAndExit(command)
}
