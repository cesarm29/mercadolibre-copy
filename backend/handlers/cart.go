package handlers

import (
	"net/http"

	"marketplace/database"
	"marketplace/models"

	"github.com/gin-gonic/gin"
)

// GetCart godoc
// @Summary      Get cart
// @Description  Get current user's cart items
// @Tags         cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.CartItem
// @Router       /api/cart [get]
func GetCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var items []models.CartItem
	database.DB.Where("user_id = ?", user.ID).Preload("Product").Preload("Product.Seller").Find(&items)

	total := 0.0
	for _, item := range items {
		total += item.Product.Price * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"count": len(items),
	})
}

// AddToCart godoc
// @Summary      Add to cart
// @Description  Add a product to the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Cart item data"
// @Success      200  {object}  models.CartItem
// @Failure      400  {object}  map[string]string
// @Router       /api/cart [post]
func AddToCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var body struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gte=1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, body.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Stock < body.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	if product.SellerID == user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add your own product to cart"})
		return
	}

	var existing models.CartItem
	if err := database.DB.Where("user_id = ? AND product_id = ?", user.ID, body.ProductID).First(&existing).Error; err == nil {
		existing.Quantity += body.Quantity
		database.DB.Save(&existing)
		database.DB.Preload("Product").First(&existing, existing.ID)
		c.JSON(http.StatusOK, existing)
		return
	}

	item := models.CartItem{
		UserID:    user.ID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
	}
	database.DB.Create(&item)
	database.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusCreated, item)
}

// UpdateCartItem godoc
// @Summary      Update cart item
// @Description  Update quantity of a cart item
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Cart item ID"
// @Param        request body object true "Update data"
// @Success      200  {object}  models.CartItem
// @Failure      400  {object}  map[string]string
// @Router       /api/cart/{id} [put]
func UpdateCartItem(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var item models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	var body struct {
		Quantity int `json:"quantity" binding:"required,gte=1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	database.DB.First(&product, item.ProductID)
	if product.Stock < body.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	item.Quantity = body.Quantity
	database.DB.Save(&item)
	database.DB.Preload("Product").First(&item, item.ID)

	c.JSON(http.StatusOK, item)
}

// RemoveFromCart godoc
// @Summary      Remove from cart
// @Description  Remove an item from the cart
// @Tags         cart
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Cart item ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/cart/{id} [delete]
func RemoveFromCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var item models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	database.DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// ClearCart godoc
// @Summary      Clear cart
// @Description  Remove all items from the cart
// @Tags         cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Router       /api/cart/clear [delete]
func ClearCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	database.DB.Where("user_id = ?", user.ID).Delete(&models.CartItem{})
	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
}
