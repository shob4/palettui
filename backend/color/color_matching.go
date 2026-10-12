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
			return []hsl{}, err
		}
		monochrome = append(monochrome, newHsl)
	}

	for l := uint16(color.L); l <= 1000; l += 50 {
		var newHsl, err = NewHsl(color.H, color.S, l)
		if err != nil {
			return []hsl{}, err
		}
		monochrome = append(monochrome, newHsl)
	}
	return monochrome, nil
}

func Gradient(color hsl, color2 hsl, num uint32) ([]hsl, error) {
	if num == 0 {
		return []hsl{}, errors.New("cannot create a 0 color gradient")
	}
	// making sure the color is within bounds
	if _, err := NewHsl(color.H, color.S, color.L); err != nil {
		return []hsl{}, err
	}
	if _, err := NewHsl(color2.H, color2.S, color2.L); err != nil {
		return []hsl{}, err
	}

	var inum = int32(num)
	var gradient = make([]hsl, 0, num)
	var hue int32

	if color.S == 0 {
		hue = int32(color2.H)
	} else {
		hue = int32(color.H)
	}

	var h = hue
	var s = int32(color.S)
	var l = int32(color.L)

	var endH = int32(color2.H)
	var sInterval = (int32(color2.S) - s) / inum
	var lInterval = (int32(color2.L) - l) / inum
	var hDifference = endH - h

	var hInterval int32
	if hDifference > 180 {
		hInterval = -(360 - hDifference) / inum
	} else if hDifference < -180 {
		hInterval = (360 + hDifference) / inum
	} else {
		hInterval = hDifference / inum
	}

	for range inum {
		h = int32(RemEuclid(int(h+hInterval), 360))
		s = s + sInterval
		l = l + lInterval

		gradient = append(gradient, hsl{uint16(h), uint16(s), uint16(l)})
	}

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
