package alert

import (
	"errors"
	"log"
	"os/exec"
)

type CmdAlerter struct {
	Alert
}

func (c *CmdAlerter) Validate() (bool, error) {
	if len(c.AlertCommand) < 1 {
		return false, errors.New("alert command is required")
	}
	if len(c.ResolutionCommand) < 1 {
		return false, errors.New("resolution command is required")
	}
	return true, nil
}

func (c *CmdAlerter) SendAlert() {
	execute(c.AlertCommand)
}

func (c *CmdAlerter) SendResolution() {
	execute(c.ResolutionCommand)
}

func execute(command []string) {
	execution := exec.Command(command[0], command[1:]...)
	output, err := execution.CombinedOutput()

	if err != nil {
		log.Printf("Command error: %s\n", err.Error())
	}
	log.Printf("Command output: %s\n", output)
}
