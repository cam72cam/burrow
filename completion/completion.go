package completion

import "github.com/cam72cam/burrow/attached"

/// Func of current params
/// Returns possible completions to last param
type CompleteParamsFunc func(p *attached.Process, curr []string) []string

type Match struct {
	Name     string
	Complete CompleteParamsFunc
}

type MatchFunc func(string) []Match
