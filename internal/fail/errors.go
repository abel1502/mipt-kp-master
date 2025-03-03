package fail

import (
	"errors"
)

var (
// TODO: Err... = new("...")
)

func new(desc string) error {
	return errors.New(desc)
}
