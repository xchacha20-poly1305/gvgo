package gvgo

var _ error = versionError{}

type versionError struct {
	message string
}

func (v versionError) Error() string {
	return "gvgo: " + v.message
}

var (
	ErrMissMain error = versionError{"main part is empty"}
	ErrMainLong error = versionError{"main part is too long"}
)
