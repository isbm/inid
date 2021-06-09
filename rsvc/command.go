package rsvc

import (
	"strings"
)

type RunitServiceCommand struct {
	command      string
	args         []string
	isConcurrent bool
}

func NewRunitServiceCommand(command string) *RunitServiceCommand {
	svc := new(RunitServiceCommand)

	commandTokens := strings.Split(command, " ")
	svc.command = commandTokens[0]
	svc.args = commandTokens[1:]
	svc.isConcurrent = false

	return svc
}

func (svc *RunitServiceCommand) SetConcurrent(isConcurrent bool) *RunitServiceCommand {
	svc.isConcurrent = isConcurrent
	return svc
}

// IsConcurrent returns true of the command supposed to be concurrent
func (svc *RunitServiceCommand) IsConcurrent() bool {
	return svc.isConcurrent
}
