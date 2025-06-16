package cli_test

import (
	"testing"
	"time"

	"github.com/kanopy-platform/go-library/cli"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestTimeFlag_Time(t *testing.T) {
	expectedTime := time.Date(1970, 12, 25, 15, 30, 45, 0, time.UTC)
	tf := cli.TimeFlag(expectedTime)

	result := tf.Time()
	assert.True(t, result.Equal(expectedTime),
		"Time() = %v, want %v", result, expectedTime)
}

func TestTimeFlag_WithPflag(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		wantErr   bool
		expected  time.Time
	}{
		{
			name:      "valid",
			flagValue: "1970-12-25T15:30:45Z",
			wantErr:   false,
			expected : time.Date(1970, 12, 25, 15, 30, 45, 0, time.UTC),
		},
		{
			name:      "invalid",
			flagValue: "invalid-time",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

			var timeFlag cli.TimeFlag

			fs.VarP(&timeFlag, "time", "t", "time flag for testing")

			args := []string{"--time", tt.flagValue}
			err := fs.Parse(args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, timeFlag.Time())
		})
	}
}

func TestTimeFlag_RoundTrip(t *testing.T) {
	// test that String() and Set() are consistent
	originalTime := time.Date(1970, 12, 25, 15, 30, 45, 123456789, time.UTC)
	tf1 := cli.TimeFlag(originalTime)

	str := tf1.String()

	var tf2 cli.TimeFlag
	err := tf2.Set(str)
	assert.NoError(t, err)

	assert.True(t, tf1.Time().Equal(tf2.Time()),
		"Round trip failed: original = %v, result = %v", tf1.Time(), tf2.Time())
}
