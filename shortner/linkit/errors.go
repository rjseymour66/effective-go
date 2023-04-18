package linkit

import "errors"

var (
	// ErrExists when the link was already created
	ErrExists = errors.New("already exists")
	// ErrNotExist when you are looking for a link that does not exist
	ErrNotExist = errors.New("does not exist")
	// ErrInternal when you get a link that we cannot expose
	ErrInternal = errors.New("internal error: please try again later or contact support")
)
