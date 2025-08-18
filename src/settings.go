package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"app/azuClient"
	"app/constants"

	"github.com/charmbracelet/log"
	"github.com/playwright-community/playwright-go"
	"golang.org/x/crypto/ssh/terminal"
)

type AppSettings struct {
	ConfigFile string
	Session    SessionConfig
}

func (a *AppSettings) SaveSettings() {
	err := MarshalAndWriteFileContents(a.ConfigFile, a.Session)
	if err != nil {
		panic(err)
	}
}

func (a *AppSettings) SavePIMToken(token azuClient.AzurePimToken) {
	a.Session.AZPimToken = token
	a.SaveSettings()
}

// func (a *AppSettings) SaveSessionData(sessionData string, url string) {
//	a.Session.SessionData[url] = sessionData
//	a.SaveSettings()
// }

func buildFilePath(filename string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home + "/.ezpim/" + filename
}

func loadSessionFromFile() AppSettings {
	appSettings := AppSettings{
		ConfigFile: buildFilePath("config.json"),
	}
	_ = CreateDirectoryStructure(filepath.Dir(appSettings.ConfigFile))

	sessionConfig := GetSessionConfig(appSettings.ConfigFile)
	appSettings.Session = sessionConfig

	return appSettings
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

func promptForCredentials() (string, string, error) {
	var username, password string
	fmt.Print("Username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		return "", "", err
	}
	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", "", err
	}
	password = string(bytePassword)
	return username, password, nil
}

func Initialize(logger *log.Logger) AppSettings {
	preflight()
	appSettings := loadSessionFromFile()
	now := time.Now().Unix()
	expiresOn, err := strconv.Atoi(appSettings.Session.AZPimToken.ExpiresOn)
	if err != nil {
		expiresOn = 0
	}
	if now > int64(expiresOn) {
		log.Info("Token expired. Please login to get new token")
		log.Info("Launching browser to get new token")
		headless := false
		username, password, err := promptForCredentials()
		if err != nil {
			logger.Warn("Failed to get credentials. You will need to manually login to get new token")
		} else {
			logger.Info("Successfully retrieved credentials")
			headless = true
		}
		LaunchBrowserToGetToken(
			appSettings, PimOptions{
				Headless:        headless,
				AppMode:         true,
				KioskMode:       true,
				PreserveSession: true,
				AzurePortalURL:  constants.AzurePortalUrl,
				Username:        username,
				Password:        password,
			},
		)
		appSettings = loadSessionFromFile()
	}

	err = appSettings.Session.AZPimToken.ComputeAdditionalFields()
	if err != nil {
		panic("Failed to compute SubjectID: " + err.Error())
	}

	return appSettings
}
