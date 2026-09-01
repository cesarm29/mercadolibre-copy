package serverutil

import (
	"net/http"
	"sync"

	"marketplace/config"
	"marketplace/database"
	"marketplace/routes"

	"github.com/gin-gonic/gin"
)

var (
	once    sync.Once
	handler http.Handler
)

func NewHandler() http.Handler {
	once.Do(func() {
		cfg := config.Load()
		database.Connect(cfg)

		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(gin.Recovery())

		r.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		routes.Setup(r)
		handler = r
	})
	return handler
}