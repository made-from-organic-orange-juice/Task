package wpsapi

import "fmt"

type ErrOutOfRange string

func (e ErrOutOfRange) Error() string {
	return fmt.Sprintf(string(e))
}
