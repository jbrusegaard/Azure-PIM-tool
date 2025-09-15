package src

import (
	"encoding/json"
	"fmt"
	"os"

	"app/azuClient"
	"app/log"
)

type SessionConfig struct {
	// SessionData map[string]string
	AZPimToken azuClient.AzurePimToken
}

func GetSessionConfig(filename string) SessionConfig {
	data, err := GetOrCreateFile(filename, "{}")
	if err != nil {
		exitWithError(log.GetLogger(), "Could not get or create config file", err.Error())
	}
	var sessionConfig SessionConfig
	err = json.Unmarshal(data, &sessionConfig)
	if err != nil {
		return SessionConfig{}
	}
	return sessionConfig
}

func GetOrCreateFile(path string, defaultContent string) ([]byte, error) {
	// check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create the file
		file, err2 := os.Create(path)
		if err2 != nil {
			return nil, err2
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				exitWithError(log.GetLogger(), "Could not close file", err.Error())
				return
			}
		}(file)
		// write default content to the file
		if _, err = file.WriteString(defaultContent); err != nil {
			return []byte{}, fmt.Errorf("could not write default content to file: %w", err)
		}
		return []byte(defaultContent), nil
	} else {
		// read the file contents
		f, err2 := os.ReadFile(path)
		if err2 != nil {
			return []byte{}, fmt.Errorf("could not read file: %w", err2)
		}
		return f, nil
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func WriteFileContents(path string, contents []byte) error {
	err := os.WriteFile(path, contents, 0644)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}

func MarshalAndWriteFileContents(path string, contents any) error {
	data, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal contents: %w", err)
	}
	return WriteFileContents(path, data)
}

func RemoveFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}
	return err
}

func CreateDirectoryStructure(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create directory structure: %w", err)
	}
	return nil
}
