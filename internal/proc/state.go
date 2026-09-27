package proc

import "fmt"

type StatusCode uint32

const (
	StatusCodeSuccess StatusCode = 0
)

func (s StatusCode) Failed() bool {
	return s > 0
}

func (s StatusCode) String() string {
	return fmt.Sprintf("%d", s)
}
