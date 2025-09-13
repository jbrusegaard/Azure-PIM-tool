package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"app/azuClient"
	"app/constants"
	"app/log"
	spinner2 "github.com/charmbracelet/bubbles/spinner"

	charmlog "github.com/charmbracelet/log"
	"github.com/playwright-community/playwright-go"
)

var Debugging = false

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
		exitWithError(log.GetLogger(), "Could not save config file", err.Error())
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
		exitWithError(log.GetLogger(), "could not determine user home directory", err.Error())
	}
	return home + "/.ezpim/" + filename
}

func loadSessionFromFile() AppSettings {
	appSettings := AppSettings{
		ConfigFile: buildFilePath("config.json"),
	}
	dirErr := CreateDirectoryStructure(filepath.Dir(appSettings.ConfigFile))
	if dirErr != nil {
		exitWithError(log.GetLogger(), "Failed to create config directory", dirErr.Error())
	}

	sessionConfig := GetSessionConfig(appSettings.ConfigFile)
	appSettings.Session = sessionConfig

	return appSettings
}

func preflight(logger *charmlog.Logger) {
	// check if playwright is installed
	if _, err := playwright.Run(); err != nil {
		iErr := playwright.Install(&playwright.RunOptions{Browsers: []string{"chromium"}})
		if iErr != nil {
			l := logger.With("PRE-FLIGHT", "Playwright")
			exitWithError(l, "Failed to install Playwright.", iErr.Error())
		}
	}
}

func Initialize(logger *charmlog.Logger, opts InitOpts) AppSettings {
	// Set up preflight checks, like installing playwright if not installed
	preflight(logger)

	// Load session from file
	appSettings := loadSessionFromFile()

	now := time.Now().Unix()

	expiresOn, err := strconv.Atoi(appSettings.Session.AZPimToken.ExpiresOn)
	if err != nil {
		expiresOn = 0
	}

	if now > int64(expiresOn) {
		logger.Info("Token expired. Please login to get new token")

		var username, password string
		if !opts.Interactive {
			username, password, err = promptForCredentials()
			if err != nil {
				logger.Warn("Failed to get credentials. You will need to manually login to get new token")
			} else {
				fmt.Println()
				logger.Info("Successfully captured credentials")
			}

		}

		spinner := StartSpinner("Starting login process", spinner2.Points)
		defer func() {
			if err3 := spinner.ReleaseTerminal(); err3 != nil {
				logger.Warn("Failed to release terminal from spinner")
			}
			spinner.Quit()
		}()

		LaunchBrowserToGetToken(
			appSettings,
			PimOptions{
				Headless:        opts.Headless,
				AppMode:         true,
				KioskMode:       true,
				PreserveSession: true,
				AzurePortalURL:  constants.AzurePortalUrl,
				Username:        username,
				Password:        password,
			},
			logger,
			spinner,
		)
		appSettings = loadSessionFromFile()
	}

	err = appSettings.Session.AZPimToken.ComputeAdditionalFields()
	if err != nil {
		exitWithError(logger, "Failed to compute SubjectID", err.Error())
	}

	return appSettings
}

func setDebugging(debug bool) {
	Debugging = debug
}

func exitWithError(logger *charmlog.Logger, basicError, debugError string) {
	if Debugging {
		logger.Error(debugError)
		panic(fmt.Errorf(debugError))
	} else {
		logger.Fatal(basicError)
	}
}
