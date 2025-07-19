package routes

import (
	"go-readme-stats/app/handler"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine) {
	app.NoRoute(ErrRouter)

	app.GET("/favicon.ico", handler.Favicon)

	route := app.Group("/api")
	{
		route.GET("/langs", handler.GetLanguageStats)
	}
}

func ErrRouter(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": "this page could not be found",
	})
}
