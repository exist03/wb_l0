package common

import "errors"

var ErrNotFound = errors.New("such id does not exist")
var ErrInvalidID = errors.New("invalid id")
