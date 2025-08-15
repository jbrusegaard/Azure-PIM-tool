package azureclient

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
)

type PimOptions struct {
	Headless        bool
	AppMode         bool
	KioskMode       bool
	PreserveSession bool
	AzurePortalURL  string
}

// Use playwright to get the pim token
// user will need to login using the launched browser
// then get the pim token
// then return the pim token
// then use the pim token to get the pim data
// then return the pim data
// then use the pim data to get the pim data
// then return the pim data
func GetPimToken(opts PimOptions) {

	pw, err := playwright.Run(
		&playwright.RunOptions{
			Browsers: []string{"chromium"},
		},
	)
	if err != nil {
		panic("could not start Playwright: " + err.Error())
	}

	var args = []string{
		"--start-maximized",
		"--disable-extensions",
		"--no-experiments",
		"--hide-crash-restore-bubble",
		"--window-name='EZ PIM'",
		"--disable-infobars",
	}

	if opts.Headless {
		args = append(args, "--headless")
	}
	if opts.AppMode {
		args = append(args, "--app=https://azure.com/", "--force-app-mode")
	}
	if opts.KioskMode {
		args = append(args, "--kiosk")
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: &opts.Headless,
		Args:     args,
	})
	if err != nil {
		panic("could not launch browser: " + err.Error())
	}

	// We need to intercept the request when the user logs in to get their pim token
	page, err := browser.NewPage()
	if err != nil {
		panic("could not create page: " + err.Error())
	}

	page.Goto(opts.AzurePortalURL)

	// Need to wait for user to auth and then go to the azure portal home page
	page.WaitForURL(opts.AzurePortalURL + "#home")
	time.Sleep(30 * time.Second)

	// Get session storage
	sessionStorageData, err := page.Evaluate("() => JSON.stringify(sessionStorage)")
	if err != nil {
		// Handle error
	}
	fmt.Printf("Session Storage: %s\n", sessionStorageData)

	time.Sleep(2 * time.Minute)

	fmt.Print("we finished")
}
