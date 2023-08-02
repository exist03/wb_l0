package utils

import (
	"errors"
	"time"
)

func DoWithTries(fn func() error, attempts int, delay time.Duration) error {
	for attempts > 0 {
		if err := fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return errors.New("0 connection attempts left: the database is not connected")
}
