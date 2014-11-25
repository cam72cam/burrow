package display

import "github.com/cam72cam/burrow/attached"

type ContextView struct {
	*Output
}

func NewContextView() *ContextView {
	return &ContextView{NewOutput()}
}

func (f *ContextView) LoadContext(p *attached.Process) error {
	pts, err := p.CurrentThreadPoints()
	if err != nil {
		return err
	}
	for id, pt := range pts {
		f.Printf("Thread %d at %#v %s:%d %s\n", id, pt.Addr, pt.File, pt.Line, pt.Func)
	}
	f.Redraw()
	return nil
}
