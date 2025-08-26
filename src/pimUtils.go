package src

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"time"

	"app/azuClient"

	terminal "golang.org/x/term"

	"github.com/playwright-community/playwright-go"
)

const userNameTextBox = "//*[@id=\"i0116\"]"
const userNameSubmitButton = "//*[@id=\"idSIButton9\"]"
const enterButton = "Enter"

type PimOptions struct {
	Headless        bool
	AppMode         bool
	KioskMode       bool
	PreserveSession bool
	AzurePortalURL  string
	Username        string
	Password        string
}

func promptForCredentials() (string, string, error) {
	username := os.Getenv("EZPIM_USERNAME")
	var password string
	if len(username) == 0 {
		fmt.Print("Username: ")
		_, err := fmt.Scanln(&username)
		if err != nil {
			return "", "", err
		}
	}
	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", "", err
	}
	password = string(bytePassword)
	fmt.Println("")
	return username, password, nil
}

func handle2FA(page playwright.Page) error {
	var multiFactorCode string
	multifactorLocator := page.GetByPlaceholder("Code")
	fmt.Print("2FA Code: ")
	_, err := fmt.Scanln(&multiFactorCode)
	if err != nil {
		return err
	}
	fmt.Println("")
	if err = multifactorLocator.Fill(multiFactorCode); err != nil {
		return err
	}
	if err = multifactorLocator.Press(enterButton); err != nil {
		return err
	}
	return nil
}

func LaunchBrowserToGetToken(appSettings AppSettings, opts PimOptions) {

	pw, err := playwright.Run(
		&playwright.RunOptions{
			Browsers: []string{"chromium"},
		},
	)
	defer func(pw *playwright.Playwright) {
		_ = pw.Stop()
	}(pw)
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

	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: &opts.Headless,
			Args:     args,
		},
	)
	defer func(browser playwright.Browser, options ...playwright.BrowserCloseOptions) {
		_ = browser.Close(options...)
	}(browser)
	if err != nil {
		panic("could not launch browser: " + err.Error())
	}

	// We need to intercept the request when the user logs in to get their pim token
	page, err := browser.NewPage()
	if err != nil {
		panic("could not create page: " + err.Error())
	}

	_, err = page.Goto(opts.AzurePortalURL)
	if err != nil {
		panic("could not goto: " + err.Error())
	}

	if opts.Username != "" && opts.Password != "" {
		if err = page.Locator(userNameTextBox).Fill(opts.Username); err != nil {
			panic(err)
		}
		if err = page.Locator(userNameSubmitButton).Click(); err != nil {
			panic("could not press email: " + err.Error())
		}
		passwordLocator := page.GetByPlaceholder("Password")
		err = passwordLocator.Fill(opts.Password)
		if err != nil {
			panic("could not fill password: " + err.Error())
		}
		err = passwordLocator.Press(enterButton)
		if err != nil {
			panic("could not press password: " + err.Error())
		}

		if opts.Headless {
			err = handle2FA(page)
			if err != nil {
				panic("could not handle 2FA: " + err.Error())
			}
		}
	}

	// Need to wait for user to auth and then go to the azure portal home page
	err = page.WaitForURL(
		opts.AzurePortalURL+"#home", playwright.PageWaitForURLOptions{
			Timeout: playwright.Float(float64(5 * time.Minute)),
		},
	)
	if err != nil {
		panic("could not wait for URL: " + err.Error())
	}
	sessionData := CaptureSessionData(page)
	for _, session := range sessionData {
		var apt azuClient.AzurePimToken
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
	sessionStorageData, err := page.Evaluate(
		`() => {
		const storage = {};
		for (let i = 0; i < sessionStorage.length; i++) {
			const key = sessionStorage.key(i);
			storage[key] = sessionStorage.getItem(key);
		}
		return storage;
	}`,
	)
	if err != nil {
		return make(map[string]string)
	}

	// Convert the interface{} to map[string]interface{} first, then to map[string]string
	if storageMap, ok := sessionStorageData.(map[string]any); ok {
		result := make(map[string]string)
		for key, value := range storageMap {
			if strValue, ok := value.(string); ok {
				result[key] = strValue
			}
		}
		return result
	}

	return make(map[string]string)
}

// func RestoreSessionData(sessionData string, page playwright.Page) {
// 	expression := `
// 	(data) => {
// 		const dataParsed = JSON.parse(data)
// 		for(const key in dataParsed) {
// 			sessionStorage.setItem(key)
// 		}
// 	}`
// 	page.Evaluate(expression, sessionData)
// }
