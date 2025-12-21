package driver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func resolvePath(path string) string {
	if filepath.IsAbs(path) { return path }
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Desktop", path)
}

func isSafePath(path string) (bool, string) {
	home, err := os.UserHomeDir()
	if err != nil { return false, "" }
	fullPath := resolvePath(path)
	absPath, err := filepath.Abs(fullPath)
	if err != nil { return false, "" }
	if strings.HasPrefix(absPath, home) { return true, absPath }
	return false, ""
}

func ListFiles(path string) string {
	if path == "" || path == "." {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, "Desktop")
	}
	safe, absPath := isSafePath(path)
	if !safe { return "Access Denied" }
	
	files, err := ioutil.ReadDir(absPath)
	if err != nil { return "Error: " + err.Error() }
	
	var res strings.Builder
	// Fixed: We explicitly use fmt here
	res.WriteString(fmt.Sprintf("Files in %s:\n", absPath))
	for _, f := range files {
		res.WriteString(f.Name() + "\n")
	}
	return res.String()
}

func ReadFile(path string) string {
	safe, absPath := isSafePath(path)
	if !safe { return "Access Denied" }
	content, err := ioutil.ReadFile(absPath)
	if err != nil { return "Error: " + err.Error() }
	s := string(content)
	if len(s) > 5000 { return s[:5000] }
	return s
}

func WriteFile(path, content string) string {
	safe, absPath := isSafePath(path)
	if !safe { return "Access Denied" }
	dir := filepath.Dir(absPath)
	os.MkdirAll(dir, 0755)
	if err := ioutil.WriteFile(absPath, []byte(content), 0644); err != nil {
		return "Error: " + err.Error()
	}
	// Fixed: We explicitly use fmt here
	return fmt.Sprintf("Success! Wrote to %s", absPath)
}
