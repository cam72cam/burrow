package display

import (
	"bufio"
	"os"
)

type FileView struct {
	*Output
}

func NewFileView() *FileView {
	return &FileView{NewOutput()}
}

func (f *FileView) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	f.line = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f.buf = append(f.buf, scanner.Text())
	}
	f.Redraw()
	return scanner.Err()
}
