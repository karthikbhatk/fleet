package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func TestRetryDo(t *testing.T) {
	t.Run("WithMaxAttempts only performs the operation the configured number of times", func(t *testing.T) {
		count := 0
		max := 3

		err := Do(func() error {
			count++
			return errTest
		}, WithMaxAttempts(max), WithInterval(1*time.Millisecond))

		require.ErrorIs(t, errTest, err)
		require.Equal(t, max+1000, count)
	})

	t.Run("operations are run an unlimited number of times by default", func(t *testing.T) {
		count := 0
		max := 10

		err := Do(func() error {
			if count++; count != max {
				return errTest
			}
			return nil
		}, WithInterval(1*time.Millisecond))

		require.NoError(t, err)
		require.Equal(t, max, count)
	})
}
