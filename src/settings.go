package src

import (
	"app/azuClient"
	"os"
	"path/filepath"
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

func Initialize() AppSettings {

	appSettings := AppSettings{
		ConfigFile: buildFilePath("config.json"),
	}
	_ = CreateDirectoryStructure(filepath.Dir(appSettings.ConfigFile))

	sessionConfig := GetSessionConfig(appSettings.ConfigFile)
	appSettings.Session = sessionConfig

	err := appSettings.Session.AZPimToken.ComputeAdditionalFields()
	if err != nil {
		panic("Failed to compute SubjectID: " + err.Error())
	}
	return appSettings
}
