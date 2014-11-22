package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cam72cam/burrow/attached"
	"github.com/cam72cam/burrow/completion"
	"github.com/cam72cam/burrow/display"
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

type CmdFn func(p *attached.Process, args ...string) error

type CommandDef struct {
	Func     CmdFn
	Help     string
	Complete completion.CompleteParamsFunc
}

var commands map[string]CommandDef

func init() {
	commands = map[string]CommandDef{
		"help":   CommandDef{help, "Help Text", helpComplete},
		"break":  CommandDef{breakpt, "Break at file.go:line", pointComplete},
		"clear":  CommandDef{clearpt, "Clear break at file.go:line", pointComplete},
		"breaks": CommandDef{breakpts, "Show all breakpoints", nil},
		"show":   CommandDef{showFile, "Show file", nil},
		"exit":   CommandDef{exit, "Exit", nil},
		"quit":   CommandDef{exit, "Exit", nil},
	}
}

func Run(p *attached.Process, cmd string, params []string) error {
	if c, ok := commands[cmd]; ok {
		return c.Func(p, params...)
	}
	return nil
}

func exit(p *attached.Process, args ...string) error {
	return ExitCMD
}

func showFile(p *attached.Process, args ...string) error {
	out := display.NewOutput()
	defer out.Close()

	for _, f := range p.Funcs() {
		out.Printf("%s", f)
	}
	return nil
}

func breakpts(p *attached.Process, args ...string) error {
	return nil
}

func pointArgs(args []string) (attached.Point, error) {
	if len(args) != 1 {
		return attached.Point{}, fmt.Errorf("Invalid number of arguments")
	}
	sp := strings.Split(args[0], ":")
	switch len(sp) {
	case 1:
		addr, err := strconv.Atoi(args[0]) //TODO Atoi is insufficient, need uint64 support
		if err != nil {
			return attached.Point{}, err
		}
		return attached.Point{Addr: uint64(addr)}, nil
	case 2:
		//TODO check file exists
		file := sp[0]
		line, err := strconv.Atoi(sp[1])
		if err != nil {
			return attached.Point{}, err
		}
		return attached.Point{File: file, Line: line}, nil
	default:
		return attached.Point{}, fmt.Errorf("Expected File:Line or Address")
	}
}
func pointComplete(args []string) []string {
	if len(args) != 1 {
		return nil
	}
	sp := strings.Split(args[0], ":")
	switch len(sp) {
	case 1: //func or addr or file
		str := sp[0]
		_, err := strconv.Atoi(str)
		if err == nil { //addr
			//TODO check valid
			return nil
		}
		res := make([]string, 0)
		/*for _, f := range p.Funcs() { //TODO
			if strings.HasPrefix(f, str) {
				res = append(res, f)
			}
		}*/

		partial := filepath.Base(str)
		base := filepath.Dir(str)
		if filepath.Base(base) == partial {
			partial = ""
		}
		dir, _ := ioutil.ReadDir(base)
		if base == "." {
			base = ""
		} else {
			base += "/"
		}
		if partial == "." {
			partial = ""
		}
		var fromFile bool
		for _, f := range dir {
			if strings.HasPrefix(f.Name(), partial) && (strings.HasSuffix(f.Name(), ".go") || f.IsDir()) {
				str = base + f.Name()
				if f.IsDir() {
					str += "/"
				}
				res = append(res, str)
				fromFile = true
			}
		}
		if len(res) == 1 && fromFile && strings.HasSuffix(res[0], ".go") {
			res[0] = res[0] + ":"
		}
		return res
	case 2: //File:line
		return nil
	default:
		return nil
	}
}
func breakpt(p *attached.Process, args ...string) error {
	pt, err := pointArgs(args)
	if err != nil {
		return err
	}
	pt, err = p.Break(pt)
	if err != nil {
		display.ShowPoint(pt)
	}
	return err
}
func clearpt(p *attached.Process, args ...string) error {
	pt, err := pointArgs(args)
	if err != nil {
		return err
	}
	pt, err = p.Clear(pt)
	if err != nil {
		display.HidePoint(pt)
	}
	return err
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
func help(p *attached.Process, args ...string) error {
	out := display.NewOutput()
	defer out.Close()

	for name, def := range commands {
		out.Printf("%s: %s", name, def.Help)
	}
	return nil
}
