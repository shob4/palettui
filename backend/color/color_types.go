// Package color provides definitions for color types and algorithms to translate them. May be broken up further
package color

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

type hsb struct {
	H uint16 `json:"h"`
	S uint16 `json:"s"`
	B uint16 `json:"b"`
}

type hex struct {
	H uint32 `json:"h"`
}

type name struct {
	Name string `json:"name"`
}

type Color struct {
	Rgb    rgb  `json:"rgb"`
	Hsl    hsl  `json:"hsl"`
	Hsb    hsb  `json:"hsb"`
	Hex    hex  `json:"hex"`
	Name   name `json:"name"`
	Locked bool `json:"locked"`
}

type Palette struct {
	Colors []Color `json:"color_list"`
	ID     uint64 `json:"id"`
	PaletteName   string `json:"palette_name"`
}
