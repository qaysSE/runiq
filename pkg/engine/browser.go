package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chromedp/chromedp"
)

var (
	GlobalTabContext context.Context
	CancelBrowser    context.CancelFunc
)

func Notify(text string) {
	if runtime.GOOS == "darwin" {
		// Basic escape
		text = strings.ReplaceAll(text, "\"", "\\\"")
		exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "Nexus Agent"`, text)).Run()
	}
}

func AskUser(action, details string) bool {
	if runtime.GOOS == "darwin" {
		// ESCAPE QUOTES to prevent AppleScript errors
		action = strings.ReplaceAll(action, "\"", "\\\"")
		details = strings.ReplaceAll(details, "\"", "\\\"")
		
		script := fmt.Sprintf(`display dialog "%s\n\n%s" buttons {"Deny", "Allow"} default button "Allow" cancel button "Deny" with title "Nexus Security" with icon caution`, action, details)
		return exec.Command("osascript", "-e", script).Run() == nil
	}
	return true 
}

func EnsureBrowser() error {
	if GlobalTabContext != nil && GlobalTabContext.Err() == nil {
		return nil
	}
	home, _ := os.UserHomeDir()
	userDataDir := filepath.Join(home, ".nexus-agent-profile")
	
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.UserDataDir(userDataDir),
		chromedp.WindowSize(1400, 900),
	)
	
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	GlobalTabContext, CancelBrowser = chromedp.NewContext(allocContext)
	return chromedp.Run(GlobalTabContext, chromedp.Navigate("about:blank"))
}
