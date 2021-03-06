package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/cam72cam/burrow/attached"
	"github.com/cam72cam/burrow/commands"
	"github.com/cam72cam/burrow/display"
	nc "github.com/gbin/goncurses"
)

func init() {

	runtime.LockOSThread()
}

func main() {

	p, err := attached.Launch(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Exit := func(status int) {
		err := p.Kill()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(status)
	}

	fn, err := display.Init()
	if err != nil {
		fmt.Println("Error initializing display: %v", err)
		Exit(1)
	}
	Exit = func(status int) {
		fn()
		err := p.Kill()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(status)
	}

	defer func() {
		if r := recover(); r != nil {
			fn()
			time.Sleep(time.Millisecond)
			fmt.Println(r)
			fmt.Printf("%s\n", debug.Stack())
		}
		err := p.Kill()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(-1)
	}()

	for {
		in := display.NextInput()
		if in.String() == ":" {
			cmd := display.NextCommand(p, commands.MatchInput)
			if cmd != nil {
				if l, err := strconv.Atoi(cmd.Name); err == nil {
					display.Curr.GoToLine(l)
					continue
				}
				err := commands.Run(p, cmd.Name, cmd.Params)
				switch err {
				case nil:
					continue
				case commands.ExitCMD:
					Exit(0)
				default:
					o := display.NewOutput()
					o.Printf("%v for pid %d", err, p.PID())
					o.Redraw()
					o.Close()
				}
			}
		} else {
			switch in {
			case nc.KEY_UP:
				display.Curr.MoveCursor(-1)
			case nc.KEY_DOWN:
				display.Curr.MoveCursor(1)
			case nc.KEY_PAGEUP:
				display.Curr.Scroll(-display.Curr.Size())
			case nc.KEY_PAGEDOWN:
				display.Curr.Scroll(display.Curr.Size())
			default:
				switch in.String() {
				case "/":
					sstr = display.SearchInput()
					l := display.Curr.FindNext(sstr)
					display.Curr.GoToLine(l)
				case "n":
					l := display.Curr.FindNext(sstr)
					display.Curr.GoToLine(l)
				case "i":
					p.StepAll()
					out := display.NewContextView()
					out.LoadContext(p)
				case "s":
					p.NextAll()
					out := display.NewContextView()
					out.LoadContext(p)
				}
			}
		}
	}
	Exit(0)
}

var sstr = ""
