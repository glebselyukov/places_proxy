package errors

import (
	"github.com/pkg/errors"
)

// errors database for places operations
var (
	ErrPlacesInternal = errors.New("internal error")
	ErrPlacesNotFound = errors.New("key not found")
	ErrPlacesExists   = errors.New("key exists")
)
