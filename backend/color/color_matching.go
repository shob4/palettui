package color

import (
	"errors"
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

func Monochromatic(color hsl) ([]hsl, error) {
	// making sure the color is within bounds
	if _, err := NewHsl(color.H, color.S, color.L); err != nil {
		return nil, err
	}
	var monochrome []hsl
	for l := uint16(50); l < color.L; l += 50 {
		var newHsl, err = NewHsl(color.H, color.S, l)
		if err != nil {
			return nil, err
		}
		monochrome = append(monochrome, newHsl)
	}

	for l := uint16(color.L); l <= 1000; l += 50 {
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

func Gradient(color hsl, color2 hsl, num uint32) ([]hsl, error) {
	if num < 2 {
		return nil, errors.New("gradients must have at least 2 colors")
	}
	// making sure the color is within bounds
	if _, err := NewHsl(color.H, color.S, color.L); err != nil {
		return nil, err
	}
	if _, err := NewHsl(color2.H, color2.S, color2.L); err != nil {
		return nil, err
	}

	const greyChroma = 1e-3
	var gradient = make([]hsl, 0, num)
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

	steps := float64(num - 1)
	distanceChroma := end.C - start.C
	distanceLightness := end.L - start.L

	for i := uint32(1); i < num - 1; i++ {
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

func GeneratePalette(num uint) ([]Color, error) {
	if num <= 0 {
		return []Color{}, errors.New("no numbers to generate")
	}

	var newPalette = make([]Color, num)
	var newColor, err = GenerateColor()
	if err != nil {
		return []Color{}, err
	}
	var i uint = 1
	newPalette = append(newPalette, newColor)
	if i < num {
		newColor, err = complement(newPalette[0].Hsl).ToColor()
		if err != nil {
			return []Color{}, err
		}
		newPalette = append(newPalette, newColor)
		i += 1
	}

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
			newColor, err = complement(newPalette[index].Hsl).ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
		case 1:
			newColor, err = GenerateColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
		case 2:
			newColor, err = nColorAverageComplement(newPalette)
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
		case 3:
			var hsl1, hsl2 hsl
			hsl1, hsl2, err = triad(newPalette[index].Hsl)
			if err != nil {
				return []Color{}, err
			}
			newColor, err = hsl1.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			newColor, err = hsl2.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			i += 1
		case 4:
			var hsl1, hsl2, hsl3 hsl
			hsl1, hsl2, hsl3, err = square(newPalette[index].Hsl)
			if err != nil {
				return []Color{}, err
			}
			newColor, err = hsl1.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			newColor, err = hsl2.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			newColor, err = hsl3.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			i += 2
		default:
			return []Color{}, errors.New("generate palette generated an invalid number")
		}
		i += 1
	}
	return newPalette, nil
}

func GeneratePaletteFromBase(currentPalette []Color, num uint64) ([]Color, error) {
	var newPalette = make([]Color, num)
	var inum = int(num)
	var i = 1
	var newColor Color
	var err error
	for i < inum {
		var method int
		if inum-i < 2 {
			method = rand.IntN(3)
		} else if inum-i < 3 {
			method = rand.IntN(4)
		} else {
			method = rand.IntN(5)
		}
		var index = rand.IntN(len(currentPalette))

		switch method {
		case 0:
			newColor, err = complement(currentPalette[index].Hsl).ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
		case 1:
			newColor, err = GenerateColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
		case 2:
			newColor, err = nColorAverageComplement(currentPalette)
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
		case 3:
			var hsl1, hsl2 hsl
			hsl1, hsl2, err = triad(currentPalette[index].Hsl)
			if err != nil {
				return []Color{}, err
			}
			newColor, err = hsl1.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			newColor, err = hsl2.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			i += 1
		case 4:
			var hsl1, hsl2, hsl3 hsl
			hsl1, hsl2, hsl3, err = square(currentPalette[index].Hsl)
			if err != nil {
				return []Color{}, err
			}
			newColor, err = hsl1.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			newColor, err = hsl2.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			newColor, err = hsl3.ToColor()
			if err != nil {
				return []Color{}, err
			}
			newPalette = append(newPalette, newColor)
			i += 2
		default:
			return []Color{}, errors.New("generate palette from base generated an invalid random number")
		}
		i += 1
	}
	return newPalette, nil
}
