package api

func newError(text string, code int) *apiError {
	return &apiError{text, code}
}

func newEmptyError() *apiError {
	return newError("", 0)
}

type apiError struct {
	m string
	c int
}

func (e *apiError) Error() string {
	return e.m
}

func (e *apiError) Code() int {
	return e.c
}
