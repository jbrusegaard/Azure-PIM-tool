package src

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
)

type SessionConfig struct {
	// SessionData map[string]string
	AZPimToken AzurePimToken
}

func GetSessionConfig(filename string) SessionConfig {
	data, err := GetOrCreateFile(filename, "{}")
	if err != nil {
		panic(err)
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
		defer file.Close()
		// write default content to the file
		if _, err = file.WriteString(defaultContent); err != nil {
			log.Errorf("Error creating file: %s", err)
			return []byte{}, err
		}
		return []byte(defaultContent), nil
	} else {
		// read the file contents
		f, err2 := os.ReadFile(path)
		if err2 != nil {
			log.Errorf("Error reading file: %s", err2)
			return []byte{}, err2
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
		log.Errorf("Error writing to file: %s", err)
	}
	return err
}

func MarshalAndWriteFileContents(path string, contents any) error {
	data, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		log.Errorf("Error marshalling contents: %s", err)
		return err
	}
	return WriteFileContents(path, data)
}

func RemoveFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	err := os.Remove(path)
	if err != nil {
		log.Errorf("Error deleting file: %s", err)
	}
	return err
}

func CreateDirectoryStructure(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Errorf("Error creating directory structure: %s", err)
		return err
	}
	return nil
}
