package app_test

import(
	"github.com/yossefsabry/gotype/internal/app"
	"testing"
)

func TestLayout_Recalculate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		width  int
		height int
		mode   app.Mode
		focus  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var l app.Layout
			l.Recalculate(tt.width, tt.height, tt.mode, tt.focus)
		})
	}
}

