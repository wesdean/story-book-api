package app_config

import (
	"errors"
	"fmt"
	"github.com/wesdean/story-book-api/logging"
	"github.com/wesdean/story-book-api/utils"
	"os/exec"
	"strings"
)

type IntegrationTestConfig struct {
	Hooks  *IntegrationTestHooksConfig
	ApiUrl string
	Logger logging.LoggerConfig
}

type IntegrationTestHooksConfig struct {
	CommandHooks *IntegrationTestCommandHooks
}

type IntegrationTestCommandHooks struct {
	BeforeTest IntegrationTestCommandConfig
	AfterTest  IntegrationTestCommandConfig
}
type IntegrationTestCommandConfig struct {
	Directory string
	Command   string
	Arguments []string
	Disabled  bool
}

func (config *IntegrationTestHooksConfig) RunHook(hookPath string) error {
	if config == nil {
		return nil
	}

	commandConfig, err := utils.GetValueFromStruct(hookPath, config)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "invalid value found at") {
			return err
		}
	}

	if commandConfig == nil {
		return nil
	}

	command, ok := commandConfig.(IntegrationTestCommandConfig)
	if !ok {
		return errors.New("invalid commandConfig config")
	}

	if command.Disabled {
		return nil
	}

	cmd := exec.Command(command.Command, command.Arguments...)
	cmd.Dir = command.Directory
	fmt.Printf("Running command: %s/%s %s",
		command.Directory,
		command.Command,
		strings.Join(command.Arguments, " "))
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
