package display

import "github.com/cam72cam/burrow/attached"

type BreakPointView struct {
	*Output
}

func NewBreakPointView() *BreakPointView {
	return &BreakPointView{NewOutput()}
}

func (bpv *BreakPointView) Show(p *attached.Process) {
	bpv.Empty()
	for addr, bp := range p.BreakPoints() {
		bpv.Printf("0x%x Func:%s File:%s:%d", addr, bp.Func, bp.File, bp.Line)
	}
	bpv.Redraw()
}
