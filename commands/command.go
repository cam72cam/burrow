package commands

import (
	"fmt"

	"github.com/cam72cam/burrow/completion"
	"github.com/derekparker/delve/proctl"
)

func MatchInput(cmd string) []completion.Match {
	res := make([]completion.Match, 0, len(commands))

	for name, def := range commands {
		if len(name) >= len(cmd) && name[:len(cmd)] == cmd {
			res = append(res, completion.Match{Name: name, Complete: def.Complete})
		}
	}
	return res
}

type CmdFn func(p *proctl.DebuggedProcess, args ...string) error

type CommandDef struct {
	Func     CmdFn
	Help     string
	Complete completion.CompleteParamsFunc
}

var commands map[string]CommandDef

func init() {
	commands = map[string]CommandDef{
		"help":   CommandDef{help, "Help Text", helpComplete},
		"exit":   CommandDef{exit, "Exit", nil},
		"quit":   CommandDef{exit, "Exit", nil},
		"quazar": CommandDef{exit, "Exit", nil},
	}
}

func Run(p *proctl.DebuggedProcess, cmd string, params []string) error {
	if c, ok := commands[cmd]; ok {
		return c.Func(p, params...)
	}
	return nil
}

func exit(p *proctl.DebuggedProcess, args ...string) error {
	return ExitCMD
}

func helpComplete(args []string) []string {
	if len(args) == 0 {
		return []string{"TODO usage"}
	}
	arg := args[0]
	possible := MatchInput(arg)
	if len(possible) == 0 {
		return []string{"Unknown command, no help availiable"}
	}

	res := make([]string, len(possible))
	for i, p := range possible { //TODO This should probably be a func on the list
		res[i] = p.Name
	}
	return res
}
func help(p *proctl.DebuggedProcess, args ...string) error {
	for name, def := range commands {
		fmt.Printf("%s: %s\n", name, def.Help)
	}
	return nil
}
