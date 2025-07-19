package routes

import (
	"go-readme-stats/app/handler"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine) {
	app.NoRoute(ErrRouter)

	app.GET("/favicon.ico", handler.Favicon)
	app.GET("/favicon.png", handler.Favicon)

	route := app.Group("/api")
	{
		route.GET("/langs", handler.GetLanguageStats)
	}
}

func ErrRouter(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": gin.H{
			"code":    404,
			"message": "Not Found",
			"details": "This API endpoint does not exist.",
		},
	})
}
