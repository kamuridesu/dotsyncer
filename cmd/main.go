package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kamuridesu/dotsyncer/internal/config"
	"github.com/kamuridesu/dotsyncer/internal/updater"
)

type Args struct {
	edit    *bool
	push    *bool
	message *string
}

func argparser() error {
	args := &Args{}
	args.message = flag.String("message", "", "commit message to be used")
	args.edit = flag.Bool("edit", false, "edit configs")
	args.push = flag.Bool("push", false, "push changes")
	flag.Parse()
	if *args.edit {
		return config.EditConfig()
	}
	return sync(*args.push, args.message)
}

func sync(push bool, message *string) error {
	conf, err := config.NewConfig()
	if err != nil {
		return err
	}
	if conf == nil {
		fmt.Println("Config file is empty! Use `-edit` to add your configs")
		os.Exit(1)
	}
	err = updater.Update(conf, push, message)
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
