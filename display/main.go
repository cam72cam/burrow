package display

import (
	nc "github.com/gbin/goncurses"
)

func Echo(state bool) {
	nc.Echo(state)
}

var s *nc.Window
var code *nc.Window
var prompt *nc.Window

func Init() (func(), error) {
	var err error
	s, err = nc.Init()
	if err != nil {
		return nil, err
	}

	Echo(false) //Echo user input off
	nc.CBreak(true)
	nc.StartColor() //Start using colors

	//nc.InitPair(1, nc.C_RED, nc.C_GREEN)
	nc.InitPair(2, nc.C_GREEN, nc.C_BLACK)
	//s.SetBackground(nc.ColorPair(1))

	s.Keypad(true)
	s.ScrollOk(true)
	y, x := s.MaxYX()

	code = s.Sub(y-2, x, 0, 0)
	code.ScrollOk(true)
	code.Touch()

	prompt = s.Sub(2, x, y-2, 0)
	prompt.Keypad(true)
	prompt.Touch()
	prompt.Refresh()

	code.Refresh()
	code.SetBackground(nc.ColorPair(2))

	//Cleanup
	return nc.End, nil
}

type Input nc.Key

func (i Input) String() string {
	return nc.KeyString(nc.Key(i))
}

func NextInput() Input {
	k := prompt.GetChar()
	if k == 0 { //input timeout
		return 0
	}

	return Input(k)
}
