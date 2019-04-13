package app_config

import (
	"encoding/json"
	"fmt"
	"github.com/wesdean/story-book-api/logging"
	"github.com/wesdean/story-book-api/utils"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"strings"
)

type Config struct {
	Logger          logging.LoggerConfig
	Models          PackageConfig
	Controllers     PackageConfig
	API             PackageConfig
	IntegrationTest IntegrationTestConfig
}

type PackageConfig struct {
	DatabaseSeed string
	Logger       logging.LoggerConfig
}

type IntegrationTestConfig struct {
	StartDocker *CommandConfig
	StopDocker  *CommandConfig
	ApiUrl      string
	Logger      logging.LoggerConfig
}

type CommandConfig struct {
	Directory string
	Command   string
	Arguments []string
}

func NewConfigFromFile(filename string) (*Config, error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *Config) Validate() error {
	return nil
}

func (config *Config) GetLogger(path string) (*logging.LoggerConfig, error) {
	iLoggerConfig, err := utils.GetValueFromStruct(path, config)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "invalid value found at") {
			return nil, err
		}
	}

	if iLoggerConfig == nil {
		return &config.Logger, nil
	}

	loggerConfig, ok := iLoggerConfig.(logging.LoggerConfig)
	if !ok {
		return &config.Logger, nil
	}

	if loggerConfig.Name.IsZero() {
		loggerConfig.Name = config.Logger.Name
	} else if !config.Logger.Name.IsZero() {
		loggerConfig.Name = null.StringFrom(
			fmt.Sprintf("%s %s", config.Logger.Name.ValueOrZero(), loggerConfig.Name.ValueOrZero()),
		)
	}

	if loggerConfig.LogPath.IsZero() {
		loggerConfig.LogPath = config.Logger.LogPath
	}

	if loggerConfig.Verbose.IsZero() {
		loggerConfig.Verbose = config.Logger.Verbose
	}

	if loggerConfig.ClearLog.IsZero() {
		loggerConfig.ClearLog = config.Logger.ClearLog
	}

	return &loggerConfig, nil
}
