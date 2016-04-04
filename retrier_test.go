package consul

import (
	"errors"
	"testing"

	"github.com/nbio/st"
)

func TestDefaultRetrier(t *testing.T) {
	retry := NewRetrier(ConstantBackoff)

	var calls int
	err := retry.Run(func() error {
		calls++
		if calls > 3 {
			return nil
		}
		return errors.New("oops")
	})

	st.Expect(t, err, nil)
	st.Expect(t, calls, 4)
}

func TestDefaultRetrierError(t *testing.T) {
	retry := NewRetrier(ConstantBackoff)

	var oops = errors.New("oops")
	var calls int
	err := retry.Run(func() error {
		calls++
		return oops
	})

	st.Expect(t, err, oops)
	st.Expect(t, calls, RetryTimes+1)
}
