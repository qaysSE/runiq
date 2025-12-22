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

// Notify displays a system notification or logs to console
func Notify(text string) {
	switch runtime.GOOS {
	case "darwin":
		text = strings.ReplaceAll(text, "\"", "\\\"")
		exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "Runiq Agent"`, text)).Run()
	case "windows":
		fmt.Fprintf(os.Stderr, "🔔 [NOTIFICATION] %s\n", text)
	default:
		exec.Command("notify-send", "Runiq Agent", text).Run()
	}
}

// AskUser requests permission
func AskUser(action, details string) bool {
	if runtime.GOOS == "darwin" {
		action = strings.ReplaceAll(action, "\"", "\\\"")
		details = strings.ReplaceAll(details, "\"", "\\\"")
		script := fmt.Sprintf(`display dialog "%s\n\n%s" buttons {"Deny", "Allow"} default button "Allow" cancel button "Deny" with title "Runiq Security" with icon caution`, action, details)
		return exec.Command("osascript", "-e", script).Run() == nil
	} else if runtime.GOOS == "windows" {
        vbsScript := fmt.Sprintf(`result = MsgBox("%s" & vbCrLf & vbCrLf & "%s", vbYesNo + vbQuestion, "Runiq Security")
If result = vbYes Then WScript.Quit(0) Else WScript.Quit(1)`, action, details)
        tmp := filepath.Join(os.TempDir(), "runiq_prompt.vbs")
        os.WriteFile(tmp, []byte(vbsScript), 0600)
        err := exec.Command("wscript", tmp).Run()
        return err == nil 
    }
	fmt.Fprintf(os.Stderr, "⚠️ [SECURITY CHECK] Auto-Allowing on Linux: %s\n", action)
	return true
}

func EnsureBrowser() error {
	if GlobalTabContext != nil && GlobalTabContext.Err() == nil {
		return nil
	}
	
	home, _ := os.UserHomeDir()
	userDataDir := filepath.Join(home, ".runiq-agent-profile")
	
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // ALWAYS OPEN (Visible)
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("restore-on-startup", false), // Prevent "Restore pages?" popup
		chromedp.UserDataDir(userDataDir),
		chromedp.WindowSize(1400, 900),
	)
	
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	GlobalTabContext, CancelBrowser = chromedp.NewContext(allocContext)
	
	// Start with a clean blank page
	return chromedp.Run(GlobalTabContext, chromedp.Navigate("about:blank"))
}

// CloseBrowser completely shuts down Chrome
func CloseBrowser() {
	if CancelBrowser != nil {
		CancelBrowser() // Kills the Chrome process
		GlobalTabContext = nil
		CancelBrowser = nil
	}
}
