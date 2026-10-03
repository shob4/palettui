// Package api provides a restful api interface for the palette generator website
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"shob4/palette_site/color"
)

var palettes = []color.Palette{}

func GetPalettes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, palettes)
}

func PostPalette(c *gin.Context) {
	var newPalette color.Palette

	if err := c.BindJSON(&newPalette); err != nil {
		return
	}

	palettes = append(palettes, newPalette)
}

func GetPaletteByName(c *gin.Context) {
	id, err := strconv.ParseUint( c.Param("id"), 10, 8)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "palette id must be a number"})
		return
	}

	for _, p := range palettes {
		if p.ID == id {
			c.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "palette not found"})
}

func GetCreatePalette(c *gin.Context) {
	var colors, err = color.GeneratePalette(5)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	var p = color.Palette{ Colors: colors, ID: 0, PaletteName: "new" }
	c.IndentedJSON(http.StatusOK, p)
}

func GetRandomizePalette(c *gin.Context, p color.Palette) {
	var count uint64
	for _, c := range p.Colors {
		if c.Locked {
			continue
		}
		count += 1
	}
	var colors, err = color.GeneratePaletteFromBase(p.Colors, count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	p.Colors = colors
	c.IndentedJSON(http.StatusOK, p)
}

func GetGradientPalette(c *gin.Context, color1 color.Color, color2 color.Color, num uint32) {
	var hsls, err = color.Gradient(color1.Hsl, color2.Hsl, num)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	var colors []color.Color
	for _, oldColor := range hsls {
		var newColor, err = oldColor.ToColor()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		colors = append(colors, newColor)
	}
	var p = color.Palette{ Colors: colors, ID: 0, PaletteName: "gradient of {color1.name} and {color2.name}"}

	c.IndentedJSON(http.StatusOK, p)
}

func GetMonochromePalette(c *gin.Context, baseColor color.Color) {
	var hsls = color.Monochromatic(baseColor.Hsl)
	var colors []color.Color
	for _, oldColor := range hsls {
		var newColor, err = oldColor.ToColor()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		colors = append(colors, newColor)
	}
	c.IndentedJSON(http.StatusOK, colors)
}
