package display

import "github.com/cam72cam/burrow/attached"

type ContextView struct {
	*FilePartialView
}

func NewContextView() *ContextView {
	return &ContextView{NewFilePartialView()}
}

func (f *ContextView) LoadContext(p *attached.Process) error {
	pts, err := p.CurrentThreadPoints()
	if err != nil {
		return err
	}
	for id, pt := range pts {
		f.Printf("Thread %d at %#v %s:%d %s\n", id, pt.Addr, pt.File, pt.Line, pt.Func)
		if pt.InFile() {
			f.FileContext(pt.File, pt.Line, 4)
			f.Print("")
		}
	}
	f.Redraw()
	return nil
}
