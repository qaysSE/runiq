package driver

import (
	"fmt"
	"os/exec"
	"runtime"
)

func LaunchApp(name string) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-a", name)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", name)
	default:
		cmd = exec.Command("xdg-open", name)
	}
	if err := cmd.Start(); err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("Launched %s", name)
}
