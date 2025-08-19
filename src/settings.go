package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"app/azuClient"
	"app/constants"

	"github.com/charmbracelet/log"
	"github.com/playwright-community/playwright-go"
)

type AppSettings struct {
	ConfigFile string
	Session    SessionConfig
}

type InitOpts struct {
	Interactive bool
	Headless    bool
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

func Initialize(logger *log.Logger, opts InitOpts) AppSettings {
	preflight()
	appSettings := loadSessionFromFile()
	now := time.Now().Unix()
	expiresOn, err := strconv.Atoi(appSettings.Session.AZPimToken.ExpiresOn)
	if err != nil {
		expiresOn = 0
	}
	if now > int64(expiresOn) {
		logger.Info("Token expired. Please login to get new token")
		logger.Info("Launching browser to get new token")
		var username, password string
		if !opts.Interactive {
			username, password, err = promptForCredentials()
		}
		if err != nil {
			logger.Warn("Failed to get credentials. You will need to manually login to get new token")
		} else {
			fmt.Println()
			logger.Info("Successfully retrieved credentials")
		}
		LaunchBrowserToGetToken(
			appSettings, PimOptions{
				Headless:        false,
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
