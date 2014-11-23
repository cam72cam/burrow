package display

import (
	"github.com/cam72cam/burrow/attached"
	nc "github.com/gbin/goncurses"
)

type Expanse struct {
	win *nc.Window
}

func createExpanse() Expanse {
	return Expanse{code} //TODO manage separately
}

func (e *Expanse) Close() {
}

func DisplayBreakpoints(p *attached.Process) {

}
