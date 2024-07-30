package gvgo

var _ error = Error{}

// Error is the error of gvgo.
type Error struct {
	reason string
}

func (e Error) Error() string {
	return "gvgo: " + e.reason
}

var (
	ErrInvalidKind error = Error{"invalid kind"}
	ErrInvalidGit  error = Error{"invalid git info"}
)
