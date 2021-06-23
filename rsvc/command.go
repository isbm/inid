package rsvc

import (
	"strings"
)

type RunServiceCommand struct {
	command      string
	args         []string
	isConcurrent bool
}

func NewRunServiceCommand(command string) *RunServiceCommand {
	svc := new(RunServiceCommand)

	commandTokens := strings.Split(command, " ")
	svc.command = commandTokens[0]
	svc.args = commandTokens[1:]
	svc.isConcurrent = false

	return svc
}

func (svc *RunServiceCommand) SetConcurrent(isConcurrent bool) *RunServiceCommand {
	svc.isConcurrent = isConcurrent
	return svc
}

// IsConcurrent returns true of the command supposed to be concurrent
func (svc *RunServiceCommand) IsConcurrent() bool {
	return svc.isConcurrent
}
