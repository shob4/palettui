// Package api provides a restful api interface for the palette generator website
package api

import (
	"net/http"
	"shob4/palette_site/color"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 8)
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
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	var p = color.Palette{Colors: colors, ID: 0, PaletteName: "new"}
	c.IndentedJSON(http.StatusOK, p)
}

func GetRandomizePalette(c *gin.Context) {
	var req struct {
		ColorList []color.Color `json:"color_list"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if len(req.ColorList) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "color_list must be more than 2"})
		return
	}
	colors, err := color.GeneratePaletteFromBase(req.ColorList)
	var p = color.Palette{Colors: colors, ID: 0, PaletteName: "new"}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, p)
}

func GetGradientPalette(c *gin.Context) {
	var req struct {
		StringList []string `json:"color_list"`
		Num       uint32        `json:"num"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	if len(req.StringList) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "two colors required"})
	}
	if req.Num < 1 || req.Num > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "num must be between 1 and 100"})
	}

	var ColorList []color.Color
	for _, str := range req.StringList {
		str = strings.TrimPrefix(str, "#")
		if len(str) != 6 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "hex must be 6 characters"})
			return
		}
		n, err := strconv.ParseUint(str, 16, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hex: " + str})
			return
		}
		newColor, err := color.Hex{H: uint32(n)}.ToRgb().Convert()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ColorList = append(ColorList, newColor)
	}
	var hsls, err = color.Gradient(ColorList[0].Hsl, ColorList[1].Hsl, req.Num)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	var colors []color.Color
	for _, oldColor := range hsls {
		var newColor, err = oldColor.ToColor()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		colors = append(colors, newColor)
	}
	var p = color.Palette{Colors: colors, ID: 0, PaletteName: "gradient of {color1.name} and {color2.name}"}

	c.IndentedJSON(http.StatusOK, p)
}

func GetMonochromePalette(c *gin.Context, baseColor color.Color) {
	var hsls, err = color.Monochromatic(baseColor.Hsl)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	var colors []color.Color
	for _, oldColor := range hsls {
		var newColor, err = oldColor.ToColor()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		colors = append(colors, newColor)
	}
	c.IndentedJSON(http.StatusOK, colors)
}
