package utils

import (
	"fmt"
	"os/exec"
)

// Provisioner is an interface for running infrastructure commands.
type Provisioner interface {
	Init(directory string, args ...string) (string, error)
	Apply(directory string, args ...string) (string, error)
	Output(directory string) (string, error)
}

// TofuProvisioner implements Provisioner for tofu.
type TofuProvisioner struct{}

func (p *TofuProvisioner) Init(directory string, args ...string) (string, error) {
	cmdArgs := append([]string{"init", "-input=false"}, args...)
	return runCommand("tofu", directory, cmdArgs...)
}

func (p *TofuProvisioner) Apply(directory string, args ...string) (string, error) {
	cmdArgs := append([]string{"apply", "-input=false", "-auto-approve"}, args...)
	return runCommand("tofu", directory, cmdArgs...)
}

func (p *TofuProvisioner) Output(directory string) (string, error) {
	cmdArgs := []string{"output", "-json"}
	return runCommand("tofu", directory, cmdArgs...)
}

// runCommand is a helper function to run a command in a given directory with arguments.
func runCommand(cmdName, directory string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running command %s %v: %v - output: %s", cmdName, args, err, string(output))
	}
	return string(output), nil
}
