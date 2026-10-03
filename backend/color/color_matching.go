package color

import (
	"errors"
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
	var newH = RemEuclid(int(color.H+180), 360)
	return hsl{
		H: uint16(newH),
		S: color.S,
		L: color.L,
	}
}

func triad(color hsl) (hsl, hsl) {
	var left = RemEuclid(int(color.H-120), 360)
	var right = RemEuclid(int(color.H+120), 360)

	var leftColor = hsl{H: uint16(left), S: color.S, L: color.L}
	var rightColor = hsl{H: uint16(right), S: color.S, L: color.L}
	return leftColor, rightColor
}

func square(color hsl) (hsl, hsl, hsl) {
	var left = RemEuclid(int(color.H-90), 360)
	var middle = RemEuclid(int(color.H+180), 360)
	var right = RemEuclid(int(color.H+90), 360)

	var leftColor = hsl{H: uint16(left), S: color.S, L: color.L}
	var middleColor = hsl{H: uint16(middle), S: color.S, L: color.L}
	var rightColor = hsl{H: uint16(right), S: color.S, L: color.L}
	return leftColor, middleColor, rightColor
}

func analagous(color hsl) (hsl, hsl) {
	var left = RemEuclid(int(color.H-30), 360)
	var right = RemEuclid(int(color.H+30), 360)

	var leftColor = hsl{H: uint16(left), S: color.S, L: color.L}
	var rightColor = hsl{H: uint16(right), S: color.S, L: color.L}
	return leftColor, rightColor
}

func Monochromatic(color hsl) []hsl {
	var monochrome []hsl
	for l := uint16(50); l < color.L; l += 50 {
		monochrome = append(monochrome, hsl{color.H, color.S, l})
	}

	for l := uint16(50); l <= 1000; l += 50 {
		monochrome = append(monochrome, hsl{color.H, color.S, l})
	}
	return monochrome
}

func Gradient(color hsl, color2 hsl, num uint32) ([]hsl, error) {
	if (color.H > 360 && color2.H > 360) || (color.S > 1000 && color2.S > 1000) || (color.L > 1000 && color2.L > 1000) {
		var empty []hsl
		return empty, errors.New("colors not within proper limits")
	}

	var inum = int32(num)
	var gradient = make([]hsl, uint64(num))
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
		hInterval = (360 - hDifference) / inum
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
	var complements []hsl
	for i := range len(nodes) {
		complements = append(complements, nodes[i].Hsl)
	}
	var redGreenBlue rgb
	var err error
	if len(complements) > 0 {
		redGreenBlue, err = complements[len(complements)-1].ToRgb()
		if err != nil {
			return Color{}, err
		}
	} else {
		return Color{}, errors.New("no colors to complement")
	}
	var r = uint32(redGreenBlue.R)
	var g = uint32(redGreenBlue.G)
	var b = uint32(redGreenBlue.B)
	var newRgb rgb
	for i := 0; i < len(complements); i++ {
		newRgb, err = complements[i].ToRgb()
		if err != nil {
			return Color{}, err
		}
		r = (r + uint32(newRgb.R)) / 2
		g = (g + uint32(newRgb.G)) / 2
		b = (b + uint32(newRgb.B)) / 2
	}
	var newColor Color
	newColor, err = rgb{uint8(r), uint8(g), uint8(b)}.Convert()
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
			var hsl1, hsl2 = triad(newPalette[index].Hsl)
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
			var hsl1, hsl2, hsl3 = square(newPalette[index].Hsl)
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
			var hsl1, hsl2 = triad(currentPalette[index].Hsl)
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
			var hsl1, hsl2, hsl3 = square(currentPalette[index].Hsl)
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
