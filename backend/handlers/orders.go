package handlers

import (
	"math"
	"net/http"
	"strconv"

	"marketplace/database"
	"marketplace/models"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary      Create order
// @Description  Create an order from cart items
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Order data"
// @Success      201  {object}  models.Order
// @Failure      400  {object}  map[string]string
// @Router       /api/orders [post]
func CreateOrder(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var body struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", user.ID).Preload("Product").Find(&cartItems)

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	total := 0.0
	var orderItems []models.OrderItem

	for _, ci := range cartItems {
		if ci.Product.Stock < ci.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for " + ci.Product.Title})
			return
		}
		total += ci.Product.Price * float64(ci.Quantity)
	}

	order := models.Order{
		UserID:          user.ID,
		Total:           total,
		Status:          "pending",
		ShippingAddress: body.ShippingAddress,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	for _, ci := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
			Price:     ci.Product.Price,
		}
		orderItems = append(orderItems, orderItem)
		database.DB.Create(&orderItem)

		database.DB.Model(&models.Product{}).Where("id = ?", ci.ProductID).
			Updates(map[string]interface{}{
				"stock":      ci.Product.Stock - ci.Quantity,
				"sold_count": ci.Product.SoldCount + ci.Quantity,
			})
	}

	database.DB.Where("user_id = ?", user.ID).Delete(&models.CartItem{})

	database.DB.Preload("Items.Product").Preload("User").First(&order, order.ID)
	c.JSON(http.StatusCreated, order)
}

// GetMyOrders godoc
// @Summary      Get my orders
// @Description  Get orders of current user
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200  {object}  models.PaginatedResponse
// @Router       /api/orders [get]
func GetMyOrders(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	var total int64
	database.DB.Model(&models.Order{}).Where("user_id = ?", user.ID).Count(&total)

	var orders []models.Order
	database.DB.Where("user_id = ?", user.ID).
		Preload("Items.Product").
		Offset((page - 1) * limit).Limit(limit).
		Order("created_at DESC").
		Find(&orders)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       orders,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	})
}

// GetOrder godoc
// @Summary      Get order
// @Description  Get order by ID
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Success      200  {object}  models.Order
// @Failure      404  {object}  map[string]string
// @Router       /api/orders/{id} [get]
func GetOrder(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var order models.Order
	if err := database.DB.Where("id = ? AND (user_id = ? OR EXISTS (SELECT 1 FROM products WHERE seller_id = ? AND id IN (SELECT product_id FROM order_items WHERE order_id = ?)))",
		c.Param("id"), user.ID, user.ID, c.Param("id")).
		Preload("Items.Product").Preload("User").First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus godoc
// @Summary      Update order status
// @Description  Update order status (seller or admin)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Param        request body object true "Status data"
// @Success      200  {object}  models.Order
// @Failure      400  {object}  map[string]string
// @Router       /api/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var order models.Order
	if err := database.DB.First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required,oneof=pending paid shipped delivered cancelled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Role.Name != "admin" {
		isSeller := false
		var items []models.OrderItem
		database.DB.Where("order_id = ?", order.ID).Find(&items)
		for _, item := range items {
			var product models.Product
			database.DB.First(&product, item.ProductID)
			if product.SellerID == user.ID {
				isSeller = true
				break
			}
		}
		if !isSeller {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
			return
		}
	}

	order.Status = body.Status
	database.DB.Save(&order)
	database.DB.Preload("Items.Product").Preload("User").First(&order, order.ID)

	c.JSON(http.StatusOK, order)
}

// GetSellerOrders godoc
// @Summary      Get seller orders
// @Description  Get orders containing seller's products
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Order
// @Router       /api/orders/seller [get]
func GetSellerOrders(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var orders []models.Order
	database.DB.
		Where("id IN (SELECT DISTINCT oi.order_id FROM order_items oi JOIN products p ON oi.product_id = p.id WHERE p.seller_id = ?)", user.ID).
		Preload("Items.Product").Preload("User").
		Order("created_at DESC").
		Find(&orders)

	c.JSON(http.StatusOK, orders)
}
