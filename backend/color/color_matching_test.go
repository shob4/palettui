package color

import (
	"fmt"
	"testing"
)

func TestRemEuclid(t *testing.T) {
	var tests = []struct {
		input int
		expected int
	}{
		{180, 180},
		{275, 275},
		{360, 0},
		{0, 0},
		{365, 5},
		{540, 180},
		{720, 0},
		{-90, 270},
		{-365, 355},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d, %d", tt.input, tt.expected)
		t.Run(testname, func(t *testing.T) {
			ans := RemEuclid(tt.input, 360)
			if ans != tt.expected {
				t.Errorf("got %d, expected %d", ans, tt.expected)
			}
		})
	}
}
