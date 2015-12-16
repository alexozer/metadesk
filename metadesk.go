package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/alexozer/metadesk/server"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		fail(err.Error())
	}
	defer srv.Close()

	root := server.NewDesktopTree(server.NewBspwm())
	cmdline := server.NewCmdline(root, srv)

	if err = runConfig(); err != nil {
		fail(err.Error())
	}

	for {
		if cmdline.Exec() != nil {
			fmt.Println("Warning: command failed to execute completely\n")
		}

		srv.UpdateSubscribers()
	}
}

func runConfig() error {
	if len(os.Args) < 3 || os.Args[1] != "-c" {
		return errors.New("Usage: metadesk -c <CONFIG_PATH>")
	}

	if exec.Command(os.Args[2]).Start() != nil {
		return fmt.Errorf("Failed to execute '%s'", os.Args[2])
	}

	return nil
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
