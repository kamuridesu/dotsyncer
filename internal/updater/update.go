package updater

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/kamuridesu/dotsyncer/internal/config"
	"github.com/kamuridesu/dotsyncer/internal/git"
)

func updateConfig(wg *sync.WaitGroup, conf config.Config, userHomeDir, textMessage string, doPush bool) error {
	defer wg.Done()
	fmt.Printf("[%s] Updating\n", conf.Name)
	folder := path.Join(userHomeDir, ".config", conf.Name)
	branch := conf.Branch
	if branch == "" {
		branch = "main"
	}
	err := git.CloneOrPull(folder, conf.Repo, branch)
	if err != nil {
		return fmt.Errorf("[%s] failed to update, error is %s", conf.Name, err)
	}
	if doPush {
		hasChanges, err := git.HasChanges(folder)
		if err != nil {
			return fmt.Errorf("[%s] failed to track changes, err is %s", conf.Name, err)
		}
		if !hasChanges {
			fmt.Printf("[%s] Working on a clean tree\n", conf.Name)
			return nil
		}
		fmt.Printf("[%s] Pushing changes to remote\n", conf.Name)
		err = git.Push(folder, branch, textMessage)
		if err != nil {
			return fmt.Errorf("[%s] failed to push changes, err is %s", conf.Name, err)
		}
	}
	return nil

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
	wg := new(sync.WaitGroup)
	wg.Add(len(configs))
	for _, conf := range configs {
		go updateConfig(wg, conf, userHomeDir, textMessage, doPush)
	}
	wg.Wait()
	return nil
}
