package internal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearConsole() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin", "freebsd":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		return fmt.Errorf("unsupported OS")
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
