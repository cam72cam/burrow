package display

import (
	"strings"

	"github.com/cam72cam/burrow/completion"
	"github.com/cam72cam/burrow/history"
	nc "github.com/gbin/goncurses"
)

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

type Command struct {
	Name   string
	Params []string
}

func NextCommand(match completion.MatchFunc) *Command {
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
		_, x := prompt.MaxYX()
		prompt.HLine(1, 0, ' ', x)
		switch k {
		case nc.KEY_RETURN:
			e := history.NewEntry(input)
			if e != nil {
				history.Add(*e)
				matches := match(e.Name)
				if len(matches) == 1 {
					return &Command{Name: matches[0].Name, Params: e.Args}
				}
				for _, match := range matches {
					if match.Name == e.Name {
						return &Command{Name: match.Name, Params: e.Args}
					}
				}
			}
			return nil
		case nc.KEY_TAB:
			var name string
			var args []string
			e := history.NewEntry(input) //Easier to use entry than duplicate split logic here
			if e != nil {
				name = e.Name
				args = e.Args
			}
			matches := match(name)
			if len(matches) > 0 {
				if len(matches) == 1 {
					m := matches[0]
					if len(m.Name) != len(name) { //no args yet, autocomplete name
						input = m.Name + " "
						pos = len(input)
					} else if m.Complete != nil { //Complete args
						suggested := m.Complete(args)
						prompt.MovePrint(1, 1, strings.Join(suggested, ", "))
					}
				} else { //Still trying to find command
					suggested := make([]string, len(matches))
					for i, m := range matches {
						suggested[i] = m.Name
					}
					prompt.MovePrint(1, 1, strings.Join(suggested, ", "))
				}
			}
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
