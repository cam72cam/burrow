package display

import (
	"strings"

	"github.com/cam72cam/burrow/history"
	nc "github.com/gbin/goncurses"
)

/// Func of current params
/// Returns possible completions to last param
type CompleteParamsFunc func(curr []string) []string

type Match struct {
	Name string
	CompleteParamsFunc
}

type MatchFunc func(string) []Match

type Command struct {
	Name   string
	Params []string
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func NextCommand(fn MatchFunc) *Command {
	var input string
	var saved string
	histdex := 0
	pos := 0

	redraw := func() {
		_, x := prompt.MaxYX()
		//prompt.Clear()
		prompt.HLine(0, 0, ' ', x)
		prompt.MoveAddChar(0, 0, ':')
		FIXME := strings.Repeat(" ", x-len(input)-2) + "_" //Forces ncurses to move the cursor properly for whatever reason
		prompt.MovePrint(0, 1, input+FIXME)
		prompt.Move(0, pos+1)
		prompt.Refresh()
	}
	redraw()

	for {
		k := prompt.GetChar()
		switch k {
		case nc.KEY_RETURN:
			e := history.NewEntry(input)
			if e != nil {
				history.Add(*e)
				return &Command{Name: e.Name, Params: e.Args}
			}
			return nil
		case nc.KEY_TAB:

		case nc.KEY_UP:
			if histdex == 0 {
				saved = input
			}
			if history.Size() > -histdex {
				histdex--
				if e := history.Get(histdex); e != nil {
					input = e.String()
				}
			}
		case nc.KEY_DOWN:
			if histdex == -1 {
				histdex++
				input = saved
			} else if histdex != 0 {
				histdex++
				e := history.Get(histdex)
				if e != nil {
					input = e.String()
				}
			}
		case nc.KEY_BACKSPACE:
			if pos > 0 {
				input = input[:pos-1] + input[pos:]
				pos = max(0, pos-1)
			}
		case nc.KEY_HOME:
			pos = 0
		case nc.KEY_END:
			pos = len(input)
		case nc.KEY_LEFT:
			pos = max(0, pos-1)
		case nc.KEY_RIGHT:
			pos = min(len(input), pos+1)
		case 547: //Left Ctrl
			ind := strings.LastIndex(input[:pos], " ")
			if ind >= 0 {
				pos = ind
			} else {
				pos = 0
			}
		case 562: //Right Ctrl
			if len(input) == pos {
				continue
			}
			ind := strings.Index(input[pos+1:], " ")
			if ind >= 0 {
				pos = ind + pos + 1
			} else {
				pos = len(input)
			}
		default:
			str := nc.KeyString(k)
			if len(str) > 1 { //Some other special char
				continue
			}
			input += str
			pos = pos + 1
		}
		redraw()
	}
	return nil
}
