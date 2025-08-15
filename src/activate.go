package src

import (
	"app/constants"
	"fmt"
	"strconv"
	"time"

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
	appSettings := Initialize()
	now := time.Now().Unix()
	expiresOn, err := strconv.Atoi(appSettings.Session.AZPimToken.ExpiresOn)
	if err != nil {
		expiresOn = 0
	}
	if now > int64(expiresOn) {
		fmt.Println("ActivatePim: Token expired")
		GetBrowserAndPage(appSettings, PimOptions{
			Headless:        false,
			AppMode:         true,
			KioskMode:       true,
			PreserveSession: true,
			AzurePortalURL:  constants.AZURE_PORTAL_URL,
		})
		appSettings = Initialize()
	}
	fmt.Println("YAY WE DID IT!!!")

}
