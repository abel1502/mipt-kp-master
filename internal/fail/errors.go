package fail

import (
	"errors"
)

var (
	ErrNoSnapshots = new("no snapshots made in this repository")
)

func new(desc string) error {
	return errors.New(desc)
}
