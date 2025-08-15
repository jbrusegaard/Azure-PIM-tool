package src

import (
	"github.com/playwright-community/playwright-go"
)

type ActivationOptions struct {
	Reason         string
	Duration       int    // Duration in hours
	ActivationType string // Type of activation, e.g., "group", "resource", "role"
	Filter         string // Filter criteria for activation
}

func preflight() {
	// check if playwright is installed
	if _, err := playwright.Run(); err != nil {
		iErr := playwright.Install(&playwright.RunOptions{Browsers: []string{"chromium"}})
		if iErr != nil {
			panic("Failed to install Playwright: " + iErr.Error())
		}
	}
}

func ActivatePim(opts ActivationOptions) {
	preflight()
	pw, err := playwright.Run()
	if err != nil {
		panic("could not start Playwright: " + err.Error())
	}

	pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{})
}
