package display

import "github.com/cam72cam/burrow/attached"

type BreakPointView struct {
	*FilePartialView
}

func NewBreakPointView() *BreakPointView {
	return &BreakPointView{NewFilePartialView()}
}

func (bpv *BreakPointView) Show(p *attached.Process) {
	bpv.Empty()
	for addr, bp := range p.BreakPoints() {
		bpv.Printf("%#v Func:%s File:%s:%d", addr, bp.Func, bp.File, bp.Line)
		if bp.InFile() {
			bpv.FileContext(bp.File, bp.Line, 4)
			bpv.Print("")
		}
	}
	bpv.Redraw()
}
