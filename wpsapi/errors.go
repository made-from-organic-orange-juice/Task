package wpsapi

import "fmt"

type ErrOutOfRange string
type ErrModuleProcessing string

func (e ErrOutOfRange) Error() string {
	return fmt.Sprintf(string(e))
}

func (e ErrModuleProcessing) Error() string {
	return fmt.Sprintf(string(e))
}
