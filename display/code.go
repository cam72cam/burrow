package display

import (
	"fmt"
	"strings"

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

type Output struct {
	line   int
	scroll int //offset from beginning of buf to top of window
	buf    []string
	Expanse
}

var Curr *Output //TODO better

func NewOutput() *Output {
	Curr = &Output{
		line:    -1,
		scroll:  0,
		buf:     make([]string, 0, 100),
		Expanse: createExpanse(),
	}
	return Curr
}

func (o *Output) Size() int {
	y, _ := o.win.MaxYX()
	return y
}

func (o *Output) Printf(fmts string, args ...interface{}) {
	line := fmt.Sprintf(fmts, args...)
	o.Print(line)
}

func (o *Output) Print(l string) {
	o.line++ //GOTO next line
	blen := len(o.buf)

	if o.line == blen { //Append
		if blen > o.Size() {
			o.Scroll(1)
		}
		o.buf = append(o.buf, l)
	} else if len(o.buf) > o.line-1 { //Insert
		if cap(o.buf) == len(o.buf) { //Resize
			newbuf := make([]string, len(o.buf), cap(o.buf)*2)
			copy(newbuf, o.buf)
			o.buf = newbuf
		}
		o.buf = o.buf[0 : len(o.buf)+1]
		copy(o.buf[o.line+1:], o.buf[o.line:])
		o.buf[o.line] = l
	} else {
		panic("Invalid line index")
	}
	for i := o.line; i <= blen; i++ {
		o.updateLine(i)
	}
}

func (o *Output) updateLine(line int) {
	if o.scroll <= line && o.scroll+o.Size() > line && len(o.buf) >= line { //onscreen and in buf
		offset := line - o.scroll //in screen
		y, _ := o.win.MaxYX()
		o.win.HLine(offset, 0, ' ', y)
		o.win.Move(offset, 0)
		o.win.Printf("%d %s", line, o.buf[line])
		o.win.Refresh()
	}
}

//Find s after line l
func (o *Output) FindNext(s string) int {
	bl := len(o.buf)
	for i := o.scroll; i < bl; i++ {
		if strings.Contains(o.buf[i], s) {
			return i
		}
	}
	return 0
}

func (o *Output) MoveCursor(lines int) {
	newpos := lines + o.line
	newpos = min(max(0, newpos), len(o.buf))
	if newpos > o.scroll+o.Size() {
		o.Scroll(newpos - o.scroll - o.Size())
	} else if newpos < o.scroll {
		o.Scroll(newpos - o.scroll)
	}
	o.line = newpos
	o.win.Move(o.line-o.scroll, 1)
	o.win.Refresh()
}

func (o *Output) GoToLine(line int) {
	o.line += o.scroll
	o.Scroll(line - o.scroll)
}

func (o *Output) Scroll(dist int) {
	blen := len(o.buf)
	if dist > 0 {
		// (dist + scroll or last line) - how far we have scrolled already
		dist = min(dist+o.scroll, blen-1) - o.scroll
		if dist > 0 {
			oldscroll := o.scroll
			o.scroll += dist
			o.win.Scroll(dist)

			for i := oldscroll + o.Size(); i < o.scroll+o.Size(); i++ {
				o.updateLine(i)
			}
		}
	} else if dist < 0 {
		dist = -dist //Easier to play with as -
		dist = min(dist, o.scroll)
		if dist > 0 {
			o.win.Scroll(-dist)
			oldscroll := o.scroll
			o.scroll -= dist

			for i := o.scroll; i <= oldscroll; i++ {
				o.updateLine(i)
			}
		}
	}

	if o.scroll > o.line {
		o.line = o.scroll
	} else if o.scroll+o.Size() < o.line {
		o.line = o.scroll + o.Size()
	}

	o.win.Move(o.line, 0)

	o.win.Refresh()
}

func ShowPoint(pt attached.Point) {
	//TODO actually create a point window
	code.MovePrint(1, 1, pt.String())
}
func HidePoint(pt attached.Point) {
	//TODO
}
