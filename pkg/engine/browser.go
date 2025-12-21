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
		// Mac Native Notification
		text = strings.ReplaceAll(text, "\"", "\\\"")
		exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "Runiq Agent"`, text)).Run()
	case "windows":
		// Windows: Log to stderr (v1.0 simple support)
		fmt.Fprintf(os.Stderr, "🔔 [NOTIFICATION] %s\n", text)
	default:
		// Linux: Try notify-send or log
		if err := exec.Command("notify-send", "Runiq Agent", text).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "🔔 [NOTIFICATION] %s\n", text)
		}
	}
}

// AskUser requests permission. 
// On Mac: Uses Native GUI. 
// On Windows: Uses VBScript for Native GUI.
func AskUser(action, details string) bool {
	if runtime.GOOS == "darwin" {
		// Mac Native Popup
		action = strings.ReplaceAll(action, "\"", "\\\"")
		details = strings.ReplaceAll(details, "\"", "\\\"")
		script := fmt.Sprintf(`display dialog "%s\n\n%s" buttons {"Deny", "Allow"} default button "Allow" cancel button "Deny" with title "Runiq Security" with icon caution`, action, details)
		return exec.Command("osascript", "-e", script).Run() == nil
	} else if runtime.GOOS == "windows" {
        // Windows Native Popup using VBScript
        vbsScript := fmt.Sprintf(`result = MsgBox("%s" & vbCrLf & vbCrLf & "%s", vbYesNo + vbQuestion, "Runiq Security")
If result = vbYes Then WScript.Quit(0) Else WScript.Quit(1)`, action, details)
        
        tmp := filepath.Join(os.TempDir(), "runiq_prompt.vbs")
        os.WriteFile(tmp, []byte(vbsScript), 0600)
        
        // Run it
        err := exec.Command("wscript", tmp).Run()
        return err == nil // Exit code 0 means Yes
    }

	// Linux / Headless Fallback: Auto-Allow with Log
	fmt.Fprintf(os.Stderr, "⚠️ [SECURITY CHECK] Auto-Allowing on Linux: %s\n", action)
	return true
}

func EnsureBrowser() error {
	if GlobalTabContext != nil && GlobalTabContext.Err() == nil {
		return nil
	}
	
	home, err := os.UserHomeDir()
	if err != nil { home = "." }

    // Universal path joining
	userDataDir := filepath.Join(home, ".runiq-agent-profile")
	
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
