// Package vercel starts a web server to generate an SVG for the user's GitHub language statistics.

package api

import (
	"net/http"

	"go-readme-stats/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	app *gin.Engine
)

func init() {
	godotenv.Load()

	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	routes.Register(app)
}

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// log.Println("GITHUB_TOKEN set:", os.Getenv("GITHUB_TOKEN") != "")
	app.ServeHTTP(w, r)
}
