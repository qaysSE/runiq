package driver

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func TypeGlobal(text string) string {
	if runtime.GOOS != "darwin" { return "Mac only" }
	script := fmt.Sprintf(`tell application "System Events" to keystroke "%s"`, text)
	exec.Command("osascript", "-e", script).Run()
	return "Typed: " + text
}

func ScreenshotDesktop() (string, string) {
	if runtime.GOOS != "darwin" { return "Mac only", "" }
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, "Desktop", "nexus_desktop.png")
	exec.Command("screencapture", "-x", "-r", path).Run()
	b, _ := ioutil.ReadFile(path)
	return "Captured " + path, base64.StdEncoding.EncodeToString(b)
}
