package errors

import (
	"github.com/pkg/errors"
)

// errors database for places operations
var (
	ErrPlacesNotFound = errors.New("key not found")
	ErrPlacesExists   = errors.New("key exists")
)
