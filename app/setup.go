package app

import (
	"github.com/gin-gonic/gin"
)

// Function to setup the app object
func SetupApp() *gin.Engine {

	// Create barebone engine
	app := gin.New()
	// Add default recovery middleware
	app.Use(gin.Recovery())

	// disabling the trusted proxy feature
	app.SetTrustedProxies(nil)

	return app
}
