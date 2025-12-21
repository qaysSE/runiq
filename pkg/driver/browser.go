package driver

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"runiq/pkg/engine"
)

func mask(ctx context.Context) {
	chromedp.Run(ctx, chromedp.Evaluate(`Object.defineProperty(navigator,'webdriver',{get:()=>undefined})`, nil))
}

func Navigate(url string) string {
	engine.EnsureBrowser()
	tCtx, cancel := context.WithTimeout(engine.GlobalTabContext, 45*time.Second)
	defer cancel() 
	mask(tCtx)
	if err := chromedp.Run(tCtx, chromedp.Navigate(url), chromedp.WaitVisible("body", chromedp.ByQuery)); err != nil {
		return err.Error()
	}
	return "Navigated " + url
}

func Screenshot() (string, string) {
	engine.EnsureBrowser()
	tCtx, cancel := context.WithTimeout(engine.GlobalTabContext, 15*time.Second)
	defer cancel() 
	var b []byte
	chromedp.Run(tCtx, chromedp.CaptureScreenshot(&b))
	home, _ := os.UserHomeDir()
	p := filepath.Join(home, "Desktop", "runiq_web.png")
	os.WriteFile(p, b, 0644)
	return "Saved " + p, base64.StdEncoding.EncodeToString(b)
}

func Click(sel string) string {
	engine.EnsureBrowser()
	tCtx, cancel := context.WithTimeout(engine.GlobalTabContext, 15*time.Second)
	defer cancel()
	mask(tCtx)
	if err := chromedp.Run(tCtx, chromedp.Click(sel)); err != nil {
		return err.Error()
	}
	return "Clicked " + sel
}

func Type(sel, txt string) string {
	engine.EnsureBrowser()
	tCtx, cancel := context.WithTimeout(engine.GlobalTabContext, 15*time.Second)
	defer cancel()
	mask(tCtx)
	if err := chromedp.Run(tCtx, chromedp.SendKeys(sel, txt)); err != nil {
		return err.Error()
	}
	return "Typed " + txt
}

func Inspect() string {
	engine.EnsureBrowser()
	tCtx, cancel := context.WithTimeout(engine.GlobalTabContext, 15*time.Second)
	defer cancel()
	var res string
	chromedp.Run(tCtx, chromedp.Evaluate(`Array.from(document.querySelectorAll('input,button,a[href]')).slice(0,50).map(e=>e.tagName).join('\n')`, &res))
	return res
}
