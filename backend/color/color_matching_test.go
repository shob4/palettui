package color

import (
	"fmt"
	"slices"
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

func mustHsl(t *testing.T, h, s, l uint16) hsl {
	t.Helper()
	color, err := NewHsl(h, s, l)
	if err != nil {
		t.Fatalf("bad test fixture: %v", err)
	}
	return color
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
			ans1, ans2, err := triad(tt.input)
			if err != nil {
				t.Errorf("%s", err)
			}
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
		input                           hsl
		expected1, expected2, expected3 hsl
	}{
		{hsl{120, 29, 98}, hsl{30, 29, 98}, hsl{300, 29, 98}, hsl{210, 29, 98}},
		{hsl{200, 29, 98}, hsl{110, 29, 98}, hsl{20, 29, 98}, hsl{290, 29, 98}},
		{hsl{100, 29, 98}, hsl{10, 29, 98}, hsl{280, 29, 98}, hsl{190, 29, 98}},
		{hsl{340, 29, 98}, hsl{250, 29, 98}, hsl{160, 29, 98}, hsl{70, 29, 98}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v, %v, %v", tt.input, tt.expected1, tt.expected2, tt.expected3)
		t.Run(testname, func(t *testing.T) {
			ans1, ans2, ans3, err := square(tt.input)
			if err != nil {
				t.Errorf("%s", err)
			}
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

func TestAnalagous(t *testing.T) {
	var tests = []struct {
		input                hsl
		expected1, expected2 hsl
	}{
		{hsl{120, 29, 98}, hsl{90, 29, 98}, hsl{150, 29, 98}},
		{hsl{200, 29, 98}, hsl{170, 29, 98}, hsl{230, 29, 98}},
		{hsl{100, 29, 98}, hsl{70, 29, 98}, hsl{130, 29, 98}},
		{hsl{340, 29, 98}, hsl{310, 29, 98}, hsl{10, 29, 98}},
		{hsl{20, 29, 98}, hsl{350, 29, 98}, hsl{50, 29, 98}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v, %v", tt.input, tt.expected1, tt.expected2)
		t.Run(testname, func(t *testing.T) {
			ans1, ans2, err := analagous(tt.input)
			if err != nil {
				t.Errorf("%s", err)
			}
			if !HslIsEqual(ans1, tt.expected1) {
				t.Errorf("got %v, expected %v", ans1, tt.expected1)
			}
			if !HslIsEqual(ans2, tt.expected2) {
				t.Errorf("got %v, expected %v", ans2, tt.expected2)
			}
		})
	}
}

func TestNewHsl(t *testing.T) {
	var tests = []struct {
		name    string
		h, s, l uint16
		wantErr bool
	}{
		{"all zero", 0, 0, 0, false},
		{"max bounds", 360, 1000, 1000, false},
		{"hue too large", 361, 0, 0, true},
		{"saturation too large", 0, 1001, 0, true},
		{"lightness too large", 0, 0, 1001, true},
		{"uint16 max", 65535, 65535, 65535, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHsl(tt.h, tt.s, tt.l)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != (hsl{tt.h, tt.s, tt.l}) {
				t.Errorf("got %v", got)
			}
			if err != nil && got != (hsl{}) {
				t.Errorf("expected zero value on error, got %v", got)
			}
		})
	}
}

func TestMonochromatic(t *testing.T) {
	steps := func(h, s uint16, ls ...uint16) []hsl {
		out := make([]hsl, len(ls))
		for i, l := range ls {
			out[i] = hsl{h, s, l}
		}
		return out
	}
	var tests = []struct {
		name     string
		input    hsl
		expected []hsl
	}{
		{"near top", hsl{200, 400, 950}, steps(200, 400,
			50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000)},
		{"off-step lightness", hsl{10, 100, 975}, steps(10, 100,
			50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 975)},
		{"max lightness", hsl{0, 0, 1000}, steps(0, 0,
			50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000)},
		{"zero lightness", hsl{0, 0, 0}, steps(0, 0,
			0, 50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950, 1000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Monochromatic(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tt.expected) {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestMonochromaticProperties(t *testing.T) {
	for l := uint16(0); l <= 1000; l += 25 {
		input := hsl{H: 120, S: 500, L: 1}
		t.Run(fmt.Sprint("L=", l), func(t *testing.T) {
			got, err := Monochromatic(input)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.ContainsFunc(got, func(c hsl) bool { return c.L == l }) {
				t.Errorf("missing input lightness %d", l)
			}
			for index, color := range got {
				if color.H != input.H || color.S != input.S {
					t.Errorf("index %d changed H/S: %v", index, color)
				}
				if color.L > 1000 {
					t.Errorf("index %d out of range: %d", index, color.L)
				}
				if index > 0 && color.L <= got[index-1].L {
					t.Errorf("not strictly ascending at %d", index)
				}
			}
		})
	}
}

func TestMonochromaticErrors(t *testing.T) {
	var tests = []struct {
		name  string
		input hsl
	}{
		{"hue too large", hsl{361, 0, 0}},
		{"saturation too large", hsl{0, 1001, 0}},
		{"lightness too large", hsl{0, 0, 1001}},
		{"would overflow uint16", hsl{0, 0, 65535}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Monochromatic(tt.input)
			if err == nil {
				t.Errorf("expected error, got %v", got)
			}
		})
	}
}

func TestGradient(t *testing.T) {
	var tests = []struct {
		name     string
		a, b     hsl
		num      uint32
		expected []hsl
	}{
		{"linear lightness", hsl{100, 500, 200}, hsl{100, 500, 600}, 5,
			[]hsl{{100, 500, 200}, {100, 500, 300}, {100, 500, 400}, {100, 500, 500}, {100, 500, 600}}},
		{"saturation and lightness together", hsl{50, 0, 1000}, hsl{50, 400, 0}, 3,
			[]hsl{{50, 0, 1000}, {50, 200, 500}, {50, 400, 0}}}, // S from 0 so hue comes from color2 (also 50)
		{"hue wraps upward", hsl{350, 500, 500}, hsl{10, 500, 500}, 3,
			[]hsl{{350, 500, 500}, {0, 500, 500}, {10, 500, 500}}},
		{"hue wraps downward", hsl{10, 500, 500}, hsl{350, 500, 500}, 3,
			[]hsl{{10, 500, 500}, {0, 500, 500}, {350, 500, 500}}},
		{"grey start takes end hue", hsl{0, 0, 500}, hsl{120, 0, 500}, 4,
			[]hsl{{0, 0, 500}, {120, 0, 500}, {120, 0, 500}, {120, 0, 500}}},
		{"single step", hsl{100, 500, 200}, hsl{100, 500, 600}, 3,
			[]hsl{{100, 500, 600}}},
		{"truncated interval stops short", hsl{100, 500, 0}, hsl{100, 500, 1000}, 3,
			[]hsl{{100, 500, 333}, {100, 500, 666}, {100, 500, 999}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gradient(tt.a, tt.b, tt.num)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tt.expected) {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}
