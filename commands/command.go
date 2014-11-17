package commands

import (
	"fmt"
	"strings"

	"github.com/derekparker/delve/proctl"
)

var history string

func init() {
	commands = []Command{
		Command{"help", help, "Help Text"},
		Command{"exit", exit, "Exit"},
		Command{"quit", exit, "Exit"},
		Command{"quazar", exit, "Exit"},
	}
}

type CmdFn func(p *proctl.DebuggedProcess, args ...string) error

type Command struct {
	Name string
	Func CmdFn
	Help string
}

var commands []Command

func Match(cmd string) (res *Command) {
	maxlen := len(cmd)
	if len(cmd) == 0 {
		return nil
	}
	for i := 1; i <= maxlen; i++ {
		curr := make([]Command, 0, len(commands))
		for _, c := range commands {
			if len(c.Name) >= maxlen && c.Name[:i] == cmd[:i] {
				curr = append(curr, c)
			}
		}
		if len(curr) == 1 {
			return &curr[0]
		}
		if i == maxlen {
			for _, c := range curr {
				if len(c.Name) == maxlen {
					return &c
				}
			}
			for _, c := range curr {
				fmt.Printf("%s, ", c.Name)
			}
		}
	}
	return
}

func Run(p *proctl.DebuggedProcess, line string) error {
	split := strings.Split(line, " ")
	cmd := Match(split[0])
	if cmd == nil {
		return nil
	}
	return cmd.Func(p, split[1:]...)
}

func exit(p *proctl.DebuggedProcess, args ...string) error {
	return ExitCMD
}

func help(p *proctl.DebuggedProcess, args ...string) error {
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Help)
	}
	return nil
}
