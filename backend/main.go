package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"shob4/palette_site/api"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5500"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}))
	router.GET("/palettes", api.GetPalettes)
	router.GET("/palette/:id", api.GetPaletteByName)
	router.GET("/palette/create", api.GetCreatePalette)
	router.POST("/palette/create-from-base", api.GetRandomizePalette)
	router.POST("/palette/gradient", api.GetGradientPalette)
	router.POST("/palette/create-monochrome", api.GetMonochromePalette)
	router.POST("/palette", api.PostPalette)
	router.POST("/palette/get-color", api.PostColor)
	router.Run("localhost:8080")
}
