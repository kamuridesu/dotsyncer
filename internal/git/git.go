package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(executable string, command string) error {
	cmd := exec.Command(executable, strings.Split(command, " ")...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()
	if err := cmd.ProcessState.ExitCode(); err != 0 {
		return fmt.Errorf("command %s failed with exit code %d", executable, err)
	}
	return nil
}

func HasChanges(folder string) (bool, error) {
	oldCwd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	defer os.Chdir(oldCwd)
	os.Chdir(folder)
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Trim(string(output), " ") != "", nil
}

func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func Reset(folder string) error {
	oldCwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldCwd)
	os.Chdir(folder)
	err = run("git", "reset")
	if err != nil {
		return err
	}
	return nil
}

func Push(folder, branch, message string) error {
	oldCwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldCwd)
	os.Chdir(folder)
	err = run("git", "add .")
	if err != nil {
		return err
	}
	if err := Commit(message); err != nil {
		return Reset(folder)
	}
	return run("git", fmt.Sprintf("push origin %s", branch))
}

func Clone(folder, repo, branch string) error {
	return run("git", fmt.Sprintf("clone %s %s --recursive --branch %s", repo, folder, branch))
}

func Pull(folder string) error {
	oldCwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldCwd)
	os.Chdir(folder)
	return run("git", "pull")
}

func CloneOrPull(folder, repo, branch string) error {
	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		return Clone(folder, repo, branch)
	}
	return Pull(folder)
}
