package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tmc/keyring"
	"io/ioutil"
	"os"
	"runtime"
)

var config *Config

type Config struct {
	JiraHost    string `json:"jiraHost"`
	Username    string `json:"username"`
	Password    string
	KeyringItem string `json:"keyringItem"`
	Editor string `json:"editor"`
}

func checkConfig(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Config does not exist, please create '%s' with the following content:\n", filename)
		fmt.Println("{")
		fmt.Println("  \"jiraHost\": \"https://jira.your-domain.com\",")
		fmt.Println("  \"username\": \"john.doe@your-domain.com\",")
		fmt.Println("  \"keyringItem\": \"jira\",")
		fmt.Println("  \"editor\": \"vim\"")
		fmt.Println("}")
		return err
	}
	return nil
}

// Determine config base path.
// On windows, this is %APPDATA%\\jira\\config.json
// On unix, this is ${XDG_CONFIG_HOME}/jira-cli/config.json (which usually
// is ${HOME}/.config/jira-cli/config.json)
func configFile() (string, error) {
	configPath := ""
	configFile := os.Getenv("JIRA_CONFIG_PATH")
	if len(configFile) > 0 {
		return configFile, nil
	}
	if runtime.GOOS == "windows" {
		configPath = os.Getenv("APPDATA")
		if len(configPath) > 0 {
			return fmt.Sprintf("%s/jira-cli/config.json", configPath), nil
		}
		return "", errors.New("Cannot detect config file!")
	}
	configPath = os.Getenv("XDG_CONFIG_HOME")
	if len(configPath) > 0 {
		return fmt.Sprintf("%s/jira-cli/config.json", configPath), nil
	}
	homeDir := os.Getenv("HOME")
	if len(homeDir) > 0 {
		return fmt.Sprintf("%s/.config/jira-cli/config.json", homeDir), nil
	}
	return "", errors.New("Cannot detect config file!")
}

func readConfig() error {
	// Determine config path
	filename, err := configFile()
	if err != nil {
		return err
	}

	err = checkConfig(filename)
	if err != nil {
		return err
	}

	// read config of file
	config = &Config{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	err = config.GetPassword()

	return err
}

func (c *Config) GetPassword() error {
	pwd, err := keyring.Get(c.KeyringItem, c.Username)
	if err == nil {
		c.Password = pwd
	}
	return err
}
