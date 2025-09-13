package src

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"time"

	"app/azuClient"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

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

func LaunchBrowserToGetToken(appSettings AppSettings, opts PimOptions, logger *log.Logger, spinner *tea.Program) {

	pw, err := playwright.Run(
		&playwright.RunOptions{
			Browsers: []string{"chromium"},
		},
	)
	defer func(pw *playwright.Playwright) {
		_ = pw.Stop()
	}(pw)
	if err != nil {
		exitWithError(logger, "could not start Playwright", fmt.Sprintf("could not start Playwright: %s", err.Error()), spinner)
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

	spinner.Send(UpdateMessageMsg{NewMessage: "Launching browser"})

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
		exitWithError(logger, "could not launch browser", fmt.Sprintf("could not launch browser: %s", err.Error()), spinner)
	}

	// We need to intercept the request when the user logs in to get their pim token
	page, err := browser.NewPage()
	if err != nil {
		exitWithError(logger, "could not create page", fmt.Sprintf("could not create page: %s", err.Error()), spinner)
	}

	spinner.Send(UpdateMessageMsg{NewMessage: "Navigating to Azure portal"})
	_, err = page.Goto(opts.AzurePortalURL)
	if err != nil {
		exitWithError(logger, fmt.Sprintf("could not goto %s", opts.AzurePortalURL), fmt.Sprintf("couldn'd go to: %s: %s", opts.AzurePortalURL, err.Error()), spinner)
	}

	if opts.Username != "" && opts.Password != "" {
		spinner.Send(UpdateMessageMsg{NewMessage: "Filling username"})
		if err = page.Locator(userNameTextBox).Fill(opts.Username); err != nil {
			exitWithError(logger, "could not fill username", fmt.Sprintf("could not fill username [%s]: %s", opts.Username, err.Error()), spinner)
		}
		if err = page.Locator(userNameSubmitButton).Click(); err != nil {
			exitWithError(logger, "could not submit username", fmt.Sprintf("could not submit username [%s]: %s", opts.Username, err.Error()), spinner)
		}

		spinner.Send(UpdateMessageMsg{NewMessage: "Filling password"})
		passwordLocator := page.GetByPlaceholder("Password")
		err = passwordLocator.Fill(opts.Password)
		if err != nil {
			exitWithError(logger, "could not fill password", fmt.Sprintf("could not fill password: %s", err.Error()), spinner)
		}
		err = passwordLocator.Press(enterButton)
		if err != nil {
			exitWithError(logger, "could not press password", fmt.Sprintf("could not press password: %s", err.Error()), spinner)
		}

		spinner.Send(UpdateMessageMsg{NewMessage: "Handling MFA"})
		if opts.Headless {
			err = handle2FA(page)
			if err != nil {
				exitWithError(logger, "could not handle 2FA", fmt.Sprintf("could not handle 2FA: %s", err.Error()), spinner)
			}
		}
	}

	spinner.Send(UpdateMessageMsg{NewMessage: "Waiting for dashboard"})

	// Need to wait for user to auth and then go to the azure portal home page
	err = page.WaitForURL(
		opts.AzurePortalURL+"#home", playwright.PageWaitForURLOptions{
			Timeout: playwright.Float(float64(5 * time.Minute)),
		},
	)
	if err != nil {
		exitWithError(logger, "could not wait for URL", "could not wait for URL: "+err.Error(), spinner)
	}

	spinner.Send(UpdateMessageMsg{NewMessage: "Capturing session data"})
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
