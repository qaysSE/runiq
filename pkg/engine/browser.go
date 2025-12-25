package engine

import (
	"context"
	"log"
	"sync"

	"github.com/chromedp/chromedp"
)

var (
	GlobalAllocContext context.Context
	GlobalAllocCancel  context.CancelFunc
	GlobalTabContext   context.Context
	GlobalTabCancel    context.CancelFunc
	browserMutex       sync.Mutex
)

// AskUser is stubbed for simplicity (Keep your full version if you want)
func AskUser(action, details string) bool { return true }

func EnsureBrowser() error {
	browserMutex.Lock()
	defer browserMutex.Unlock()

	// 1. If alive, return
	if GlobalTabContext != nil && GlobalTabContext.Err() == nil {
		return nil
	}

	// 2. Clean up old mess
	if GlobalAllocCancel != nil {
		GlobalAllocCancel()
	}

	// 3. LOGGING: Print Chrome errors to STDERR so we see them in the test
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.WindowSize(1400, 900),
	)

	// 4. Create Allocator
	GlobalAllocContext, GlobalAllocCancel = chromedp.NewExecAllocator(context.Background(), opts...)

	// 5. Create Tab (But DO NOT Navigate yet)
	GlobalTabContext, GlobalTabCancel = chromedp.NewContext(GlobalAllocContext, chromedp.WithLogf(log.Printf))

	return nil
}

func CloseBrowser() {
	if GlobalAllocCancel != nil {
		GlobalAllocCancel()
	}
	GlobalAllocContext = nil
	GlobalAllocCancel = nil
	GlobalTabContext = nil
	GlobalTabCancel = nil
}
