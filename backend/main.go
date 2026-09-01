package main

import (
	"log"
	"net/http"

	"marketplace/config"
	"marketplace/database"
	"marketplace/serverutil"

	"github.com/gin-gonic/gin"
)

// @title           Marketplace API
// @version         1.0
// @description     Shopping cart REST API like MercadoLibre
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	database.Connect(cfg)

	_ = gin.Mode()
	handler := serverutil.NewHandler()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}