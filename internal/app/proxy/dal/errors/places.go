package errors

import (
	"github.com/pkg/errors"
)

// dal errors for places transactions
var (
	ErrPlacesInternal = errors.New("internal error")
	//ErrPlacesNotFound = errors.New("key not found")
	//ErrPlacesExists   = errors.New("key exists")
)
