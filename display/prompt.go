package display

import (
	"strconv"
	"strings"

	"github.com/cam72cam/burrow/attached"
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

func commonPrefix(s []string) string {
	if len(s) == 0 {
		return ""
	}
	res := s[0]
	for _, str := range s {
		for i, r := range str {
			if i >= len(res) {
				break
			} else if uint8(r) != res[i] { //May the gods of unicode forgive me
				res = res[:i]
			}
		}
	}
	return res
}

func NextCommand(p *attached.Process, match completion.MatchFunc) *Command {
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
		switch k {
		case nc.KEY_RETURN:
			prompt.HLine(0, 0, ' ', x)
			prompt.Refresh()
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
				if _, err := strconv.Atoi(e.Name); err == nil {
					return &Command{Name: e.Name}
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
			prompt.HLine(1, 0, ' ', x)
			if len(matches) > 0 {
				suggested := make([]string, len(matches))
				for i, m := range matches {
					suggested[i] = m.Name
				}
				shortest := commonPrefix(suggested)
				if len(matches) == 1 || pos > 0 && input[pos-1] == uint8(' ') || len(args) > 0 { //Singe or cmd + " "
					m := matches[0]
					if len(matches) > 1 {
						for _, mm := range matches {
							if mm.Name == e.Name {
								m = mm
								break
							}
						}
					}
					if len(m.Name) != len(name) { //no args yet, autocomplete name
						input = m.Name + " "
						pos = len(input)
					} else if m.Complete != nil { //Complete args
						suggested := m.Complete(p, args)
						if len(suggested) > 0 {
							input = m.Name + " " + strings.Join(e.Args[:len(e.Args)-1], " ") + commonPrefix(suggested)
							pos = len(input)

							if len(suggested) > 1 {
								sugstr := strings.Join(suggested, ", ")
								if len(sugstr) > x {
									sugstr = sugstr[:x]
								}
								prompt.MovePrint(1, 0, sugstr)
							}
						}
					}
				} else {
					prompt.MovePrint(1, 1, strings.Join(suggested, ", "))
					input = shortest
					pos = len(input)
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

func SearchInput() string {
	prompt.Move(0, 0)
	prompt.Print("/")
	sstr := ""
	Echo(true)
	defer Echo(false)
	for {
		k := prompt.GetChar()
		switch k {
		case nc.KEY_RETURN:
			_, x := prompt.MaxYX()
			prompt.HLine(0, 0, ' ', x)
			prompt.Refresh()
			return sstr
		case nc.KEY_BACKSPACE:
			if len(sstr) > 0 {
				sstr = sstr[0 : len(sstr)-1]
			}
		default:
			ks := nc.KeyString(k)
			if len(ks) == 1 {
				sstr += ks
			}
		}
	}
}
