package devices

func newError(t string) *EmptyResults {
	return &EmptyResults{t}
}

type EmptyResults struct {
	s string
}

func (e *EmptyResults) Error() string {
	return e.s
}

func IsEmptyResultErr(e error) bool {
	_, ok := e.(*EmptyResults)
	return ok
}
