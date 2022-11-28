//go:generate mockgen -source clock.go -destination clock_mock.go -package clock
package clock

import "time"

type Clock interface {
	Now() time.Time
}

type Real struct{}

func (Real) Now() time.Time {
	return time.Now()
}
