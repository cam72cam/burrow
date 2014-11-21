package attached

import (
	"errors"
	"fmt"

	"github.com/derekparker/delve/proctl"
)

func Attach(pid int) (*Process, error) {
	dbp, err := proctl.Attach(pid)
	if err != nil {
		return nil, err
	}
	return &Process{dbp: dbp}, nil
}

func Launch(cmd []string) (*Process, error) {
	dbp, err := proctl.Launch(cmd)
	if err != nil {
		return nil, err
	}
	return &Process{dbp: dbp}, nil
}

type Process struct {
	dbp *proctl.DebuggedProcess
}

func (p *Process) Kill() error {
	return p.dbp.Process.Kill()
}

func (p *Process) Funcs() []string {
	res := make([]string, 0)
	for _, f := range p.dbp.GoSymTable.Funcs {
		res = append(res, f.Name)
	}
	return res
}

func (p *Process) Step() error {
	return p.dbp.Next()
}

func (p *Process) Clear(pt Point) (Point, error) {
	if pt.InFile() {
		if err := pt.fromFile(p); err != nil {
			return Point{}, err
		}
	}
	_, err := p.dbp.Clear(pt.Addr)
	return pt, err
}

func (p *Process) Break(pt Point) (Point, error) {
	if pt.InFile() {
		if err := pt.fromFile(p); err != nil {
			return Point{}, err
		}
	}
	_, err := p.dbp.Break(uintptr(pt.Addr))
	return pt, err
}

type Point struct {
	Func string
	File string
	Line int
	Addr uint64
}

func (pt *Point) String() string {
	return fmt.Sprintf("0x%x %s (%s:%d)", pt.Addr, pt.Func, pt.File, pt.Line)
}

var ErrAddrNotFound = errors.New("Error Address not found!")

func (pt *Point) fromAddr(p *Process) error {
	f, l, fn := p.dbp.GoSymTable.PCToLine(pt.Addr)
	if fn != nil {
		pt.File = f
		pt.Line = l
		pt.Func = fn.Name
		return nil
	}
	return ErrAddrNotFound
}

func (pt *Point) fromFile(p *Process) error {
	pc, fn, err := p.dbp.GoSymTable.LineToPC(pt.File, pt.Line)
	if err != nil {
		return err
	}
	pt.Addr = pc
	pt.Func = fn.Name
	return nil
}

func (pos Point) InFile() bool {
	return pos.File != ""
}

func (p *Process) CurrentPoint() (pt Point, err error) {
	regs, err := p.dbp.Registers()
	if err != nil {
		return pt, err
	}
	pt = Point{Addr: regs.PC()}
	return pt, pt.fromAddr(p)
}
