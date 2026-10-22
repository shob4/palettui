package color

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
)

func RemEuclid(a, b int) int {
	var rem = a % b
	if rem < 0 {
		return b + rem
	}
	return rem
}

func RemEuclidFloat(a, b float64) float64 {
	var rem = math.Mod(a, b)
	if rem < 0 {
		return b + rem
	}
	return rem
}

func complement(color hsl) hsl {
	var h = int(color.H)
	var newH = RemEuclid(h+180, 360)
	return hsl{
		H: uint16(newH),
		S: color.S,
		L: color.L,
	}
}

func triad(color hsl) (hsl, hsl, error) {
	var h = int(color.H)
	var left = RemEuclid(h-120, 360)
	var right = RemEuclid(h+120, 360)

	var leftColor, rightColor hsl
	var err error
	leftColor, err = NewHsl(uint16(left), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, err
	}
	rightColor, err = NewHsl(uint16(right), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, err
	}
	return leftColor, rightColor, nil
}

func square(color hsl) (hsl, hsl, hsl, error) {
	var h = int(color.H)
	var left = RemEuclid(h-90, 360)
	var middle = RemEuclid(h+180, 360)
	var right = RemEuclid(h+90, 360)

	var leftColor, middleColor, rightColor hsl
	var err error
	leftColor, err = NewHsl(uint16(left), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, hsl{}, err
	}
	middleColor, err = NewHsl(uint16(middle), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, hsl{}, err
	}
	rightColor, err = NewHsl(uint16(right), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, hsl{}, err
	}
	return leftColor, middleColor, rightColor, nil
}

func analagous(color hsl) (hsl, hsl, error) {
	var h = int(color.H)
	var left = RemEuclid(h-30, 360)
	var right = RemEuclid(h+30, 360)

	var leftColor, rightColor hsl
	var err error
	leftColor, err = NewHsl(uint16(left), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, err
	}
	rightColor, err = NewHsl(uint16(right), color.S, color.L)
	if err != nil {
		return hsl{}, hsl{}, err
	}
	return leftColor, rightColor, nil
}

// returns a monochromatic scale of 10 colors based on the given color
// errors on color failing to convert to color.Color type
func Monochromatic(color hsl) ([]hsl, error) {
	// making sure the color is within bounds
	if _, err := NewHsl(color.H, color.S, color.L); err != nil {
		return nil, err
	}
	var monochrome []hsl
	for l := ( color.L % 100 ); l < color.L; l += 100 {
		var newHsl, err = NewHsl(color.H, color.S, l)
		if err != nil {
			return nil, err
		}
		monochrome = append(monochrome, newHsl)
	}

	for l := uint16(color.L); l <= 1000; l += 100 {
		var newHsl, err = NewHsl(color.H, color.S, l)
		if err != nil {
			return nil, err
		}
		monochrome = append(monochrome, newHsl)
	}
	return monochrome, nil
}

func divRound(n, d int64) int64 {
	if n < 0 {
		return -((-n + d/2) / d)
	}
	return (n + d/2) / d
}

func divRoundFloat(n, d float64) float64 {
	if n < 0 {
		return -((-n + d/2) / d)
	}
	return (n + d/2) / d
}

// returns a gradient of num length with color and color2 as the endpoints
// errors if color fails to convert to color.Color type or num less than 1
func Gradient(color hsl, color2 hsl, num uint32) ([]hsl, error) {
	if num < 1 {
		return nil, errors.New("gradients must have at least 1 step")
	}
	// making sure the color is within bounds
	if _, err := NewHsl(color.H, color.S, color.L); err != nil {
		return nil, err
	}
	if _, err := NewHsl(color2.H, color2.S, color2.L); err != nil {
		return nil, err
	}

	const greyChroma = 1e-3
	var gradient = make([]hsl, 0, num+2)
	var start, end = HslToOklch(color), HslToOklch(color2)
	gradient = append(gradient, color)

	startHue, endHue := start.H, end.H
	if start.C == greyChroma {
		startHue = endHue
	} else if end.C == greyChroma {
		endHue = startHue
	}

	distanceHue := endHue - startHue
	if distanceHue > 180 {
		distanceHue -= 360
	} else if distanceHue < -180 {
		distanceHue += 360
	}

	steps := float64(num + 1)
	distanceChroma := end.C - start.C
	distanceLightness := end.L - start.L

	for i := 1; i <= int(num); i++ {
		t := float64(i) / steps
		gradient = append(gradient, OklchToHsl(oklch{
			L: start.L + distanceLightness*t,
			C: start.C + distanceChroma*t,
			H: RemEuclidFloat(startHue+distanceHue*t, 360),
		}))
	}
	gradient = append(gradient, color2)

	return gradient, nil
}

func nColorAverageComplement(nodes []Color) (Color, error) {
	if len(nodes) == 0 {
		return Color{}, errors.New("No colors to complement")
	}
	var s, l uint32 = 0, 0
	var sinSum, cosSum float64 = 0, 0
	for _, node := range nodes {
		color := complement(node.Hsl)
		radius := float64(color.H) * math.Pi / 180
		sinSum += math.Sin(radius)
		cosSum += math.Cos(radius)
		s += uint32(color.S)
		l += uint32(color.L)
	}
	count := uint32(len(nodes))
	h := math.Atan2(sinSum, cosSum) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	s = s / count
	l = l / count
	var newHsl, err = NewHsl(uint16(math.Round(h))%360, uint16(s), uint16(l))
	if err != nil {
		return Color{}, err
	}
	newColor, err := newHsl.ToColor()
	if err != nil {
		return Color{}, err
	}
	return newColor, nil
}

// returns a random color.
// will fail if the color fails to convert to the color.Color type
func GenerateColor() (Color, error) {
	var h = rand.IntN(361)
	var s = rand.IntN(1001)
	var l = rand.IntN(1001)
	var newColor, err = hsl{uint16(h), uint16(s), uint16(l)}.ToColor()
	if err != nil {
		return Color{}, err
	}
	return newColor, nil
}

// returns a new palette of colors, only takes a length
// errors when a generated color fails to convert to the color.Color type
// or num is less than 1
func GeneratePalette(num uint) ([]Color, error) {
	if num < 1 {
		return nil, errors.New("no numbers to generate")
	}

	var newPalette = make([]Color, 0, num)
	var newColor, err = GenerateColor()
	if err != nil {
		return nil, err
	}
	var i uint = 1
	newPalette = append(newPalette, newColor)

	for i < num {
		var method int
		if num-i < 2 {
			method = rand.IntN(3)
		} else if num-i < 3 {
			method = rand.IntN(4)
		} else {
			method = rand.IntN(5)
		}

		var index = rand.IntN(len(newPalette))
		switch method {
		case 0:
			newPalette, err = addDistinct(newPalette, complement(newPalette[index].Hsl))
			if err != nil {
				return nil, err
			}
		case 1:
			newColor, err = GenerateColor()
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, newColor.Hsl)
			if err != nil {
				return nil, err
			}
		case 2:
			newColor, err = nColorAverageComplement(newPalette)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, newColor.Hsl)
			if err != nil {
				return nil, err
			}
		case 3:
			var hsl1, hsl2 hsl
			hsl1, hsl2, err = triad(newPalette[index].Hsl)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl1)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl2)
			if err != nil {
				return nil, err
			}
			i += 1
		case 4:
			var hsl1, hsl2, hsl3 hsl
			hsl1, hsl2, hsl3, err = square(newPalette[index].Hsl)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl1)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl2)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl3)
			if err != nil {
				return nil, err
			}
			i += 2
		default:
			return nil, errors.New("generate palette generated an invalid number")
		}
		i += 1
	}
	return newPalette, nil
}

// takes a base palette and returns a new palette. Any unlocked colors are
// replaced with new ones.
// errors when a new color fails to convert to the color.Color type or
// when there is only one color in the palette
func GeneratePaletteFromBase(currentPalette []Color) ([]Color, error) {
	var inum = len(currentPalette)
	if inum < 2 {
		return nil, errors.New("palette must be at least 2")
	}
	var newPalette = make([]Color, 0, inum)
	for _, color := range currentPalette {
		if color.Locked {
			newPalette = append(newPalette, color)
		}
	}
	locked := len(newPalette)
	basePalette := make([]Color, locked)
	copy(basePalette, newPalette)

	var newColor Color
	var err error
	for len(newPalette) < inum {
		remaining := inum - len(newPalette)
		var method int
		if remaining < 2 {
			method = rand.IntN(3)
		} else if remaining < 3 {
			method = rand.IntN(4)
		} else {
			method = rand.IntN(5)
		}
		var base = basePalette[rand.IntN(locked)]

		switch method {
		case 0:
			newPalette, err = addDistinct(newPalette, complement(base.Hsl))
			if err != nil {
				return nil, err
			}
		case 1:
			newColor, err = GenerateColor()
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, newColor.Hsl)
			if err != nil {
				return nil, err
			}
		case 2:
			newColor, err = nColorAverageComplement(currentPalette)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, newColor.Hsl)
			if err != nil {
				return nil, err
			}
		case 3:
			var hsl1, hsl2 hsl
			hsl1, hsl2, err = triad(base.Hsl)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl1)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl2)
			if err != nil {
				return nil, err
			}
		case 4:
			var hsl1, hsl2, hsl3 hsl
			hsl1, hsl2, hsl3, err = square(base.Hsl)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl1)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl2)
			if err != nil {
				return nil, err
			}
			newPalette, err = addDistinct(newPalette, hsl3)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("generate palette from base generated an invalid random number")
		}
	}
	result := make([]Color, inum)
	next := locked
	for index, color := range currentPalette {
		if color.Locked {
			result[index] = color
			continue
		}
		result[index] = newPalette[next]
		next++
	}
	return result, nil
}

func clampInt(v, lo, hi int) int { return max(lo, min(hi, v)) }

func jitter(c hsl) hsl {
	h := (int(c.H) + rand.IntN(31) - 15 + 360) % 360
	s := clampInt(int(c.S)+rand.IntN(301)-150, 0, 1000)
	l := clampInt(int(c.L)+rand.IntN(341)-170, 0, 1000)
	return hsl{H: uint16(h), S: uint16(s), L: uint16(l)}
}

const minDistance = 0.06

func getMinDistance(palette []Color, c Color) float64 {
	var distance float64 = math.Inf(1)
	for _, p := range palette {
		distance = math.Min(distance, OkDistance(p.Hsl, c.Hsl))
	}
	return distance
}

func addDistinct(palette []Color, base hsl) ([]Color, error) {
	var furthestcolor Color
	var furthestDistance = math.Inf(1)
	for range 100 {
		c, err := jitter(base).ToColor()
		if err != nil {
			return nil, err
		}
		newDistance := getMinDistance(palette, c)
		if newDistance >= minDistance {
			return append(palette, c), nil
		}
		if newDistance > furthestDistance {
			furthestcolor = c
			furthestDistance = newDistance
		}
	}
	return append(palette, furthestcolor), nil
}
