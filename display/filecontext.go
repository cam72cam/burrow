package display

import (
	"bufio"
	"os"
)

type FilePartialView struct {
	*Output
}

func NewFilePartialView() *FilePartialView {
	return &FilePartialView{NewOutput()}
}

func (f *FilePartialView) FileContext(path string, line, offset int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		ind := "\t   "
		if i == line {
			ind = "\t=> "
		}
		if i >= line-offset && i <= line+offset {
			f.Print(ind + scanner.Text())
		}
		i++
	}
	return scanner.Err()
}
