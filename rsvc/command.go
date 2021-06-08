package rsvc

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-cmd/cmd"
)

type RunitServiceCommand struct {
	command      string
	args         []string
	isConcurrent bool
	cmd          *cmd.Cmd
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

func (svc *RunitServiceCommand) Start() {
	svc.cmd = cmd.NewCmd(svc.command, svc.args...)
	if svc.isConcurrent {
		x := svc.cmd.Start()
		fmt.Println(spew.Sprintf("Concurrent: %s", spew.Sdump(x)))
	} else {
		x := <-svc.cmd.Start()
		fmt.Println(spew.Sprintf("Concurrent: %s", spew.Sdump(x)))
	}
}
