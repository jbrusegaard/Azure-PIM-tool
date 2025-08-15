package src

import (
	"encoding/json"
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

func GetBrowserAndPage(appSettings AppSettings, opts PimOptions) {

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
	page.WaitForURL(opts.AzurePortalURL+"#home", playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(float64(5 * time.Minute)),
	})
	sessionData := CaptureSessionData(page)
	for _, session := range sessionData {
		var apt AzurePimToken
		err := json.Unmarshal([]byte(session), &apt)
		if err != nil {
			continue
		}
		if apt.TokenType == "Bearer" && apt.CredentialType == "AccessToken" && apt.Secret != "" {
			appSettings.SavePIMToken(apt)
			break
		}
	}
}

func CaptureSessionData(page playwright.Page) map[string]string {
	// Get session storage
	// TODO fix this function
	sessionStorageData, err := page.Evaluate("() => sessionStorage")
	if err != nil {
		return make(map[string]string)
	}
	return sessionStorageData.(map[string]string)
}

func RestoreSessionData(sessionData string, page playwright.Page) {
	expression := `
	(data) => {
		const dataParsed = JSON.parse(data)
		for(const key in dataParsed) {
			sessionStorage.setItem(key)
		}
	}`
	page.Evaluate(expression, sessionData)
}
