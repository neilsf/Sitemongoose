package alert

import (
	"log"
	"os/exec"
)

type CmdAlerter struct {
	Alert
}

func execute(command []string) {
	execution := exec.Command(command[0], command[1:]...)
	_, err := execution.Output()

	if err != nil {
		log.Printf("Command error: %s\n", err.Error())
	}
}

func (c *CmdAlerter) SendAlert() {
	execute(c.AlertCommand)
}

func (c *CmdAlerter) SendResolution() {
	execute(c.ResolutionCommand)
}
