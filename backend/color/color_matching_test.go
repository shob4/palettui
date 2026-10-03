package color

import (
	"fmt"
	"testing"
)

func HslIsEqual(color1, color2 hsl) bool {
	if color1.H != color2.H {
		return false
	}
	if color1.S != color2.S {
		return false
	}
	if color1.L != color2.L {
		return false
	}
	return true
}

func TestRemEuclid(t *testing.T) {
	var tests = []struct {
		input    int
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

func TestComplement(t *testing.T) {
	var tests = []struct {
		input    hsl
		expected hsl
	}{
		{hsl{120, 29, 98}, hsl{300, 29, 98}},
		{hsl{300, 29, 98}, hsl{120, 29, 98}},
		{hsl{180, 29, 98}, hsl{0, 29, 98}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v", tt.input, tt.expected)
		t.Run(testname, func(t *testing.T) {
			ans := complement(tt.input)
			if !HslIsEqual(ans, tt.expected) {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}

func TestTriad(t *testing.T) {
	var tests = []struct {
		input                hsl
		expected1, expected2 hsl
	}{
		{hsl{120, 29, 98}, hsl{0, 29, 98}, hsl{240, 29, 98}},
		{hsl{200, 29, 98}, hsl{80, 29, 98}, hsl{320, 29, 98}},
		{hsl{100, 29, 98}, hsl{340, 29, 98}, hsl{220, 29, 98}},
		{hsl{340, 29, 98}, hsl{220, 29, 98}, hsl{100, 29, 98}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v, %v", tt.input, tt.expected1, tt.expected2)
		t.Run(testname, func(t *testing.T) {
			ans1, ans2 := triad(tt.input)
			if !HslIsEqual(ans1, tt.expected1) {
				t.Errorf("got %v, expected %v", ans1, tt.expected1)
			} 
			if !HslIsEqual(ans2, tt.expected2) {
				t.Errorf("got %v, expected %v", ans2, tt.expected2)
			}
		})
	}
}

func TestSquare(t *testing.T) {
	var tests = []struct {
		input                hsl
		expected1, expected2, expected3 hsl
	}{
		{hsl{120, 29, 98}, hsl{30, 29, 98}, hsl{300, 29, 98}, hsl{210, 29, 98}},
		{hsl{200, 29, 98}, hsl{110, 29, 98}, hsl{20, 29, 98}, hsl{290, 29, 98}},
		{hsl{100, 29, 98}, hsl{10, 29, 98}, hsl{280, 29, 98}, hsl{190, 29, 98}},
		{hsl{340, 29, 98}, hsl{250, 29, 98}, hsl{160, 29, 98}, hsl{70, 29, 98}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v, %v", tt.input, tt.expected1, tt.expected2)
		t.Run(testname, func(t *testing.T) {
			ans1, ans2, ans3 := square(tt.input)
			if !HslIsEqual(ans1, tt.expected1) {
				t.Errorf("got %v, expected %v", ans1, tt.expected1)
			} 
			if !HslIsEqual(ans2, tt.expected2) {
				t.Errorf("got %v, expected %v", ans2, tt.expected2)
			}
			if !HslIsEqual(ans3, tt.expected3) {
				t.Errorf("got %v, expected %v", ans3, tt.expected3)
			} 
		})
	}
}
