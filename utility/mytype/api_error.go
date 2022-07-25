package mytype

type ApiError struct {
	Status int
	Err    error
	Hint   string
}

func (a ApiError) Error() string {
	if a.Hint != "" {
		return a.Hint
	}

	if a.Err == nil {
		return ""
	}

	return a.Err.Error()
}
