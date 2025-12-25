package driver

import (
	"context"
	"time"

	"runiq/pkg/engine"

	"github.com/chromedp/chromedp"
)

func Navigate(url string) string {
	if err := engine.EnsureBrowser(); err != nil {
		return "Error starting browser: " + err.Error()
	}

	// 60s Timeout for loading pages (Standard for real web use)
	ctx, cancel := context.WithTimeout(engine.GlobalTabContext, 60*time.Second)
	defer cancel()

	chromedp.Run(ctx, chromedp.Evaluate(`Object.defineProperty(navigator,'webdriver',{get:()=>undefined})`, nil))

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// WaitVisible is good for production to ensure page loaded
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)

	if err != nil {
		return "Navigation failed: " + err.Error()
	}
	return "Navigated to " + url
}

// (Keep your other tools: Click, Type, etc. here)
func Click(sel string) string      { /* ... restore from previous working version ... */ return "" }
func Type(sel, txt string) string  { /* ... restore from previous working version ... */ return "" }
func Screenshot() (string, string) { /* ... restore from previous working version ... */ return "", "" }
func Inspect() string              { /* ... restore from previous working version ... */ return "" }
