package updater

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kamuridesu/dotsyncer/internal/config"
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

func commit(message string) error {
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

func push(folder, branch, message string) error {
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
	commit(message)
	return run("git", fmt.Sprintf("push origin %s", branch))
}

func clone(folder, repo, branch string) error {
	return run("git", fmt.Sprintf("clone %s %s --recursive --branch %s", repo, folder, branch))
}

func pull(folder string) error {
	oldCwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldCwd)
	os.Chdir(folder)
	return run("git", "pull")
}

func cloneOrPull(folder, repo, branch string) error {
	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		return clone(folder, repo, branch)
	}
	return pull(folder)
}

func Update(configs []config.Config, doPush bool, message *string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	textMessage := `fix: updated via dotsyncer`
	if message != nil && *message != "" {
		textMessage = *message
	}
	for _, conf := range configs {
		fmt.Printf("Updating %s\n", conf.Name)
		folder := path.Join(userHomeDir, ".config", conf.Name)
		branch := conf.Branch
		if branch == "" {
			branch = "main"
		}
		err := cloneOrPull(folder, conf.Repo, branch)
		if err != nil {
			return fmt.Errorf("failed to update %s, error is %s", conf.Name, err)
		}
		if doPush {
			fmt.Printf("Pushing changes to remote for %s\n", conf.Name)
			err := push(folder, branch, textMessage)
			if err != nil {
				return fmt.Errorf("failed to push changes to %s, err is %s", conf.Name, err)
			}
		}
	}
	return nil
}
