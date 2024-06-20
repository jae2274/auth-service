package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationMarshal(t *testing.T) {

	t.Run("2 hours should be marshaled to 2h0m0s", func(t *testing.T) {
		d := Duration(2 * time.Hour)

		bytes, err := d.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, "2h0m0s", string(bytes))
	})

	t.Run("2h should be unmarshaled to 2 hours", func(t *testing.T) {
		d := Duration(0)
		err := d.UnmarshalJSON([]byte("2h"))
		require.NoError(t, err)
		require.Equal(t, Duration(2*time.Hour), d)
	})

	t.Run("2h0m0s should be unmarshaled to 2 hours", func(t *testing.T) {
		d := Duration(0)
		err := d.UnmarshalJSON([]byte("2h0m0s"))
		require.NoError(t, err)
		require.Equal(t, Duration(2*time.Hour), d)
	})

	t.Run("2hours 30minutes 5seconds should be marshaled to 2h30m5s", func(t *testing.T) {
		d := Duration(2*time.Hour + 30*time.Minute + 5*time.Second)

		bytes, err := d.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, "2h30m5s", string(bytes))
	})

	t.Run("2h30m5s should be unmarshaled to 2hours 30minutes 5seconds", func(t *testing.T) {
		d := Duration(0)
		err := d.UnmarshalJSON([]byte("2h30m5s"))
		require.NoError(t, err)
		require.Equal(t, Duration(2*time.Hour+30*time.Minute+5*time.Second), d)
	})
}
