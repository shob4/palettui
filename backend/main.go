package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"shob4/palette_site/api"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5500"},
	}))
	router.GET("/palettes", api.GetPalettes)
	router.GET("/palette/:id", api.GetPaletteByName)
	router.GET("/palette/create", api.GetCreatePalette)
	router.POST("/palette", api.PostPalette)
	router.Run("localhost:8080")
}
