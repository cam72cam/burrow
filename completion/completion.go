package completion

/// Func of current params
/// Returns possible completions to last param
type CompleteParamsFunc func(curr []string) []string

type Match struct {
	Name     string
	Complete CompleteParamsFunc
}

type MatchFunc func(string) []Match
