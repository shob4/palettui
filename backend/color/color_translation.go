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
	if h < 0 {
		h += 360.0
	}

	var s float64
	if delta == 0.0 {
		s = 0.0
	} else {
		lNormalized := l / 1000.0
		deltaNormalized := delta / 1000.0
		s = deltaNormalized / (1.0 - math.Abs(2.0*lNormalized-1.0) ) * 1000.0
	}

	var hsl = hsl{
		H: uint16(h),
		S: uint16(s),
		L: uint16(l),
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

	var hex = Hex{
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

func (h Hex) ToRgb() rgb {
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

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func threeNodeDistance(rgb1 rgb, rgb2 rgb) uint32 {
	var r = Abs(int32(rgb1.R) - int32(rgb2.R))
	var g = Abs(int32(rgb1.G) - int32(rgb2.G))
	var b = Abs(int32(rgb1.B) - int32(rgb2.B))

	var distance = r + g + b
	return uint32(distance)
}

// ----- hsl to oklch -----

func srgbToLinear(color float64) float64 {
	if color <= 0.04046 {
		return color / 12.92
	}
	return math.Pow((color+0.055)/1.055, 2.4)
}

func linearToSrgb(color float64) float64 {
	if color <= 0.0031308 {
		return 12.92 * color
	}
	return 1.055*math.Pow(color, 1/2.4) - 0.055
}

func clamp01(x float64) float64 {
	return math.Max(0, math.Min(1, x))
}

func hslToSrgb(color hsl) (r, g, b float64) {
	s := float64(color.S) / 1000
	l := float64(color.L) / 1000
	chroma := (1 - math.Abs(2*l-1)) * s
	hp := math.Mod(float64(color.H), 360) / 60
	x := chroma * (1 - math.Abs(math.Mod(hp, 2) - 1))

	var r1, g1, b1 float64
	switch {
	case hp < 1:
		r1, g1, b1 = chroma, x, 0
	case hp < 2:
		r1, g1, b1 = x, chroma, 0
	case hp < 3:
		r1, g1, b1 = 0, chroma, x
	case hp < 4:
		r1, g1, b1 = 0, x, chroma
	case hp < 5:
		r1, g1, b1 = x, 0, chroma
	default:
		r1, g1, b1 = chroma, 0, x
	}
	m := l - chroma/2
	return r1 + m, g1 + m, b1 + m
}

func srgbToHsl(r, g, b float64) hsl {
	r, g, b = clamp01(r), clamp01(g), clamp01(b)
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min
	l := (max + min) / 2

	var h, s float64
	if delta > 0 {
		s = delta / (1 - math.Abs(2*l - 1))
		switch max {
		case r:
			h = math.Mod((g-b)/delta, 6)
		case g:
			h = (b - r)/delta + 2
		default:
			h = (r-g)/delta + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return hsl{
		H: uint16(math.Round(h)) % 360,
		S: uint16(math.Round(clamp01(s) * 1000)),
		L: uint16(math.Round(clamp01(l) * 1000)),
	}
}

func linearToOklab(r, g, b float64) (L, a, bb float64) {
	l := math.Cbrt(0.4122214708*r + 0.5363325363*g + 0.0514459929*b)
	m := math.Cbrt(0.2119034982*r + 0.6806995451*g + 0.1073969566*b)
	s := math.Cbrt(0.0883024619*r + 0.2817188376*g + 0.6299787005*b)

	L = 0.2104542553*l + 0.7936177850*m - 0.0040720468*s
	a = 1.9779984951*l - 2.4285922050*m + 0.4505937099*s
	bb = 0.0259040371*l + 0.7827717662*m - 0.8086757660*s
	return
}

func oklabToLinear(L, a, b float64) (r, g, bl float64) {
	l := L + 0.3963377774*a + 0.2158037573*b
	m := L - 0.1055613458*a - 0.0638541728*b
	s := L - 0.0894841775*a - 1.2914855480*b
	l, m, s = l*l*l, m*m*m, s*s*s

	r = 4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g = -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	bl = -0.0041960863*l - 0.7034186147*m + 1.7076147010*s
	return
}

func HslToOklch(color hsl) oklch {
	r, g, b := hslToSrgb(color)
	L, a, bb := linearToOklab(srgbToLinear(r), srgbToLinear(g), srgbToLinear(b))

	chroma := math.Hypot(a, bb)
	h := math.Atan2(bb, a) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	if chroma < 1e-4 {
		h = 0
	}
	return oklch{L: L, C: chroma, H: h}
}

func oklchToLinear(color oklch) (r, g, b float64) {
	rad := color.H * math.Pi / 180
	return oklabToLinear(color.L, color.C*math.Cos(rad), color.C*math.Sin(rad))
}

func inGamut(color oklch) bool {
	const eps = 1e-4
	r, g, b := oklchToLinear(color)
	in := func(x float64) bool {
		return x >= -eps && x <= 1 +eps
	}
	return in(r) && in(g) && in(b)
}

func OklchToHsl(color oklch) hsl {
	if !inGamut(color) {
		lo, hi := 0.0, color.C
		for range 20 {
			mid := (lo + hi) / 2
			if inGamut(oklch{color.L, mid, color.H}) {
				lo = mid
			} else {
				hi = mid
			}
		}
		color.C = lo
	}
	r, g, b := oklchToLinear(color)
	return srgbToHsl(linearToSrgb(clamp01(r)), linearToSrgb(clamp01(g)), linearToSrgb(clamp01(b)))
}

func toLab(c hsl) (L, a, b float64) {
	o := HslToOklch(c)
	rad := o.H * math.Pi / 180
	return o.L, o.C * math.Cos(rad), o.C * math.Sin(rad)
}

func OkDistance(color1, color2 hsl) float64 {
	L1, a1, b1 := toLab(color1)
	L2, a2, b2 := toLab(color2)
	dL, da, db := L1-L2, a1-a2, b1-b2
	return math.Sqrt(dL * dL + da*da + db*db)
}
