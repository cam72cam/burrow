package commands

import (
	"errors"
)

var ExitEOF = errors.New("User requested exit via EOF")
var ExitCMD = errors.New("USer requested exit via Command")
