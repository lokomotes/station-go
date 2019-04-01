package types

import (
	"fmt"
)

// NotFoundErr is an error that resource was not found.
type NotFoundErr struct {
	What string
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprintf("Resource not found: %s", e.What)
}

// ImageNotFoundErr is an error that image was not found.
type ImageNotFoundErr struct {
	What string
}

func (e *ImageNotFoundErr) Error() string {
	return fmt.Sprintf("Image not found: %s", e.What)
}

// UnexpectedErr is an error that this version of `Station` does not implement handler for `What`
type UnexpectedErr struct {
	What string
}

func (e *UnexpectedErr) Error() string {
	return fmt.Sprintf("")
}
