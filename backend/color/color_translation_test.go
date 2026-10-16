package color

import (
	"testing"
)

func TestHslToOklchRountTrip(t *testing.T) {
	near := func(a, b uint16, tol int) bool {
		d := int(a) - int(b)
		return d >=  -tol && d <= tol
	}
	for h := uint16(0); h < 360; h += 15 {
		for _, s := range []uint16{0, 250, 500, 1000} {
			for _, l := range []uint16{100, 500, 900} {
				input := hsl{h, s, l}
				output := OklchToHsl(HslToOklch(input))
				hueOk := s == 0 || near(output.H, input.H, 2) || near(output.H, input.H+360, 2) || !near(output.S, input.S, 3) || !near(output.L, input.L, 3)
				if !hueOk || !near(output.S, input.S, 3) || !near(output.L, input.L, 3) {
					t.Errorf("%v round-tripped to %v", input, output)
				}
			}
		}
	}
}

func TestOklchToHslGamutMapping(t *testing.T) {
	output := OklchToHsl(oklch{0.9, 0.4, 150})
	if output.H > 360 || output.S > 1000 || output.L > 1000 {
		t.Errorf("out of range: %v", output)
	}
}
