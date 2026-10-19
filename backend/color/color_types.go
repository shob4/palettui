// Package color provides definitions for color types and algorithms to translate them. May be broken up further
package color

import "errors"

type rgb struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

type hsl struct {
	H uint16 `json:"h"`
	S uint16 `json:"s"`
	L uint16 `json:"l"`
}

type oklch struct {
	L float64
	C float64
	H float64
}

func NewHsl(h, s, l uint16) (hsl, error) {
	if h > 360 {
		return hsl{}, errors.New("Invalid hue")
	}
	if s > 1000 {
		return hsl{}, errors.New("Invalid saturation")
	}
	if l > 1000 {
		return hsl{}, errors.New("Invalid luminance: too large")
	}

	return hsl{H: h, S: s, L: l}, nil
}

type hsb struct {
	H uint16 `json:"h"`
	S uint16 `json:"s"`
	B uint16 `json:"b"`
}

type Hex struct {
	H uint32 `json:"h"`
}

type name struct {
	Name string `json:"name"`
}

type Color struct {
	Rgb    rgb  `json:"rgb"`
	Hsl    hsl  `json:"hsl"`
	Hsb    hsb  `json:"hsb"`
	Hex    Hex  `json:"hex"`
	Name   name `json:"name"`
	Locked bool `json:"locked"`
}

// type Kind string
//
// const (
// 	KindPalette Kind = "palette"
// 	KindGradient Kind = "gradient"
// 	KindMonochrome Kind = "monochrome"
// )
//
type Palette struct {
	Colors      []Color `json:"color_list"`
	ID          uint64  `json:"id"`
	PaletteName string  `json:"palette_name"`
	// Kind        Kind  `json:"kind"`
}
