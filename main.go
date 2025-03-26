package main

import (
	"fmt"
	"os"
	"strings"

	"flag"

	"github.com/kamuridesu/dotsyncer/internal/config"
	"github.com/kamuridesu/dotsyncer/internal/updater"
)

func argparser() error {
	flag.Parse()
	verb := strings.Join(flag.Args(), " ")
	if verb == "edit" {
		return config.EditConfig()
	}
	push := false
	if verb == "push" {
		push = true
	}
	return sync(push)
}

func sync(push bool) error {
	conf, err := config.NewConfig()
	if err != nil {
		return err
	}
	if conf == nil {
		fmt.Println("Config file is empty!")
		os.Exit(1)
	}
	err = updater.Update(conf, push)
	if err != nil {
		return fmt.Errorf("error while syncinc, err is %s", err)
	}
	return nil
}

func main() {
	err := argparser()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
