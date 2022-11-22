package zerolog

import (
	"fmt"

	"github.com/pkg/errors"
)

func ESPackMarshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	e, ok := err.(stackTracer)
	if !ok {
		return nil
	}
	for _, frame := range e.StackTrace() {
		fmt.Printf("%+s:%d\r\n", frame, frame)
	}
	return nil
}
