package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name   string `yaml:"name"`
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch"`
}

func (conf *Config) validate() error {
	if conf.Name == "" {
		return fmt.Errorf("missing config name")
	}
	if conf.Repo == "" {
		return fmt.Errorf("missing repo for %s", conf.Name)
	}
	return nil
}

func GetConfigPath() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error fetching user config dir, err is '%s'", err)
	}
	configPath := path.Join(home, "dotsyncer")
	os.MkdirAll(configPath, os.ModePerm)
	configFilePath := path.Join(configPath, "config.yaml")
	return configFilePath, nil
}

func LoadConfig(path string) ([]Config, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var config []Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	for _, conf := range config {
		if err := conf.validate(); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func NewConfig() ([]Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(path, []byte(""), os.ModePerm)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return LoadConfig(path)
}

func EditConfig() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(path, []byte(""), os.ModePerm)
		if err != nil {
			return err
		}
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
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
