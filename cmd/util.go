package cmd

import (
	"fmt"
	"os"
	"os/user"
)

func getClusterSelector() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("runned-by=%s", fmt.Sprintf("%s_%s", usr.Username, hostname)), nil
}
