package routes

import (
	"marketplace/handlers"
	"marketplace/middleware"

	_ "marketplace/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	api.POST("/auth/login", handlers.Login)
	api.POST("/auth/register", handlers.Register)

	api.GET("/products", handlers.GetProducts)
	api.GET("/products/:id", handlers.GetProduct)
	api.GET("/categories", handlers.GetCategories)
	api.GET("/categories/:slug", handlers.GetCategory)

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/auth/profile", handlers.GetProfile)
		auth.PUT("/auth/profile", handlers.UpdateProfile)

		auth.POST("/products", middleware.SellerMiddleware(), handlers.CreateProduct)
		auth.PUT("/products/:id", handlers.UpdateProduct)
		auth.DELETE("/products/:id", handlers.DeleteProduct)
		auth.GET("/products/mine", handlers.GetMyProducts)

		auth.GET("/cart", handlers.GetCart)
		auth.POST("/cart", handlers.AddToCart)
		auth.PUT("/cart/:id", handlers.UpdateCartItem)
		auth.DELETE("/cart/:id", handlers.RemoveFromCart)
		auth.DELETE("/cart/clear", handlers.ClearCart)

		auth.POST("/orders", handlers.CreateOrder)
		auth.GET("/orders", handlers.GetMyOrders)
		auth.GET("/orders/:id", handlers.GetOrder)
		auth.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
		auth.GET("/orders/seller", handlers.GetSellerOrders)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/stats", handlers.GetDashboardStats)
		admin.GET("/users", handlers.GetUsers)
		admin.PUT("/users/:id/role", handlers.UpdateUserRole)
		admin.PUT("/users/:id/toggle", handlers.ToggleUserActive)
		admin.POST("/categories", handlers.CreateCategory)
		admin.GET("/products", handlers.GetAllProductsAdmin)
	}
}
