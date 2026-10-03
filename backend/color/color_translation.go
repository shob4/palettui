package color

import (
	"errors"
	"math"
)

func (r rgb) Convert() (Color, error) {
	var red = math.Round((float64(r.R) / 255.0) * 1000.0)
	var green = math.Round((float64(r.G) / 255.0) * 1000.0)
	var blue = math.Round((float64(r.B) / 255.0) * 1000.0)
	var bigger = max(int32(red), int32(blue))
	var smaller = min(int32(red), int32(blue))
	var cMax = float64(max(bigger, int32(green)))
	var cMin = float64(min(smaller, int32(green)))
	var delta = cMax - cMin

	var l = (cMax + cMin) / 2.0

	var h float64
	if delta == 0.0 {
		h = 0.0
	} else if cMax == red {
		h = 60.0 * math.Mod((green-blue)/delta, 6.0)
	} else if cMax == green {
		h = 60.0 * ((blue-red)/delta + 2.0)
	} else if cMax == blue {
		h = 60.0 * ((red-green)/delta + 4.0)
	} else {
		return Color{}, errors.New("c_max does not match r, g, or b")
	}

	var s float64
	if delta == 0.0 {
		s = 0.0
	} else {
		s = delta / math.Abs(1.0-(2.0*(1/1000.0)-1.0))
	}

	var hsl = hsl{
		H: uint16(h),
		S: uint16(s),
		L: uint16(l),
	}

	if int32(math.Round(cMax)) == int32(math.Round(red)) {
		h = 60.0 * math.Mod((green-blue)/delta, 6.0)
	} else if int32(math.Round(cMax)) == int32(math.Round(green)) {
		h = 60.0 * ((blue-red)/delta + 2.0)
	} else if int32(math.Round(cMax)) == int32(math.Round(blue)) {
		h = 60.0 * ((red-green)/delta + 4.0)
	} else {
		h = 0.0
	}

	if cMax == 0.0 {
		s = 0.0
	} else {
		s = (delta / cMax) * 1000.0
	}
	var b = cMax

	var hsb = hsb{
		H: uint16(h),
		S: uint16(s),
		B: uint16(b),
	}

	var hex = hex{
		H: uint32(r.R)<<16 | uint32(r.G)<<8 | uint32(r.B),
	}

	var n string
	var minDistance uint32 = math.MaxUint32
	for key, value := range NAMED_COLORS {
		if value == r {
			n = key
			break
		} else {
			var newDistance = threeNodeDistance(r, value)
			if newDistance < minDistance {
				minDistance = newDistance
				n = key
			}
		}
	}

	var name = name{Name: n}

	return Color{
		Rgb:    r,
		Hsl:    hsl,
		Hsb:    hsb,
		Hex:    hex,
		Name:   name,
		Locked: false,
	}, nil
}

func (hsl hsl) ToRgb() (rgb, error) {
	if hsl.H > 360 {
		return rgb{}, errors.New("hue too high")
	}
	if hsl.S > 1000 {
		return rgb{}, errors.New("saturation too high")
	}
	if hsl.L > 1000 {
		return rgb{}, errors.New("lightness too high")
	}

	var region = hsl.H / 60
	var h = float64(hsl.H)
	var s = float64(hsl.S) / 1000.0
	var l = float64(hsl.L) / 1000.0

	var c = (1.0 - math.Abs(2.0*l-1.0)) * s
	var x = c * (1.0 - math.Abs(math.Mod(h/60.0, 2.0)-1.0))
	var m = l - c/2.0

	var r, g, b float64
	switch region {
	case 0:
		r, g, b = c, x, 0.0
	case 1:
		r, g, b = x, c, 0.0
	case 2:
		r, g, b = 0.0, c, x
	case 3:
		r, g, b = 0.0, x, c
	case 4:
		r, g, b = x, 0.0, c
	default:
		r, g, b = c, 0.0, x
	}

	var red, green, blue = uint8(math.Round((r + m) * 255.0)), uint8(math.Round((g + m) * 255.0)), uint8(math.Round((b + m) * 255.0))
	return rgb{
		R: red,
		G: green,
		B: blue,
	}, nil
}

func (hsl hsl) ToColor() (Color, error) {
	var newRgb rgb
	var newColor Color
	var err error
	newRgb, err = hsl.ToRgb()
	if err != nil {
		return Color{}, err
	}
	newColor, err = newRgb.Convert()
	if err != nil {
		return Color{}, err
	}
	return newColor, nil
}

func (hsb hsb) ToRgb() (rgb, error) {
	if hsb.H > 360 {
		return rgb{}, errors.New("hue too high")
	}
	if hsb.S > 1000 {
		return rgb{}, errors.New("saturation too high")
	}
	if hsb.B > 1000 {
		return rgb{}, errors.New("brightness too high")
	}

	var region = hsb.H / 60
	var hue = float64(hsb.H)
	var saturation = float64(hsb.S) / 1000.0
	var brightness = float64(hsb.B) / 1000.0

	var c = brightness * saturation
	var x = c * (1.0 - math.Abs(math.Mod(hue/60.0, 2.0)-1.0))
	var m = brightness - c

	var r, g, b float64
	switch region {
	case 0:
		r, g, b = c, x, 0.0
	case 1:
		r, g, b = x, c, 0.0
	case 2:
		r, g, b = 0.0, c, x
	case 3:
		r, g, b = 0.0, x, c
	case 4:
		r, g, b = x, 0.0, c
	default:
		r, g, b = c, 0.0, x
	}

	var red, green, blue = uint8(math.Round((r + m) * 255.0)), uint8(math.Round((g + m) * 255.0)), uint8(math.Round((b + m) * 255.0))
	return rgb{
		R: red,
		G: green,
		B: blue,
	}, nil
}

func (hsb hsb) ToColor() (Color, error) {
	var newRgb rgb
	var newColor Color
	var err error
	newRgb, err = hsb.ToRgb()
	if err != nil {
		return Color{}, err
	}
	newColor, err = newRgb.Convert()
	if err != nil {
		return Color{}, err
	}
	return newColor, nil
}

func (h hex) ToRgb() rgb {
	var r = uint8((h.H >> 16) & 0xff)
	var g = uint8((h.H >> 8) & 0xff)
	var b = uint8(h.H & 0xff)

	return rgb{
		R: r,
		G: g,
		B: b,
	}
}

func (n name) ToRgb() rgb {
	return NAMED_COLORS[n.Name]
}

func threeNodeDistance(rgb1 rgb, rgb2 rgb) uint32 {
	var r = (int32(rgb1.R) - int32(rgb2.R))
	var g = (int32(rgb1.G) - int32(rgb2.G))
	var b = (int32(rgb1.B) - int32(rgb2.B))

	var distance = r + g + b
	return uint32(distance)
}
