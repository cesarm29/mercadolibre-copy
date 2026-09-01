package handlers

import (
	"math"
	"net/http"
	"strconv"

	"marketplace/database"
	"marketplace/models"

	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary      List products
// @Description  Get paginated list of products with filters
// @Tags         products
// @Produce      json
// @Param        page     query int false "Page number" default(1)
// @Param        limit    query int false "Items per page" default(20)
// @Param        search   query string false "Search term"
// @Param        category query int false "Category ID"
// @Param        min_price query number false "Minimum price"
// @Param        max_price query number false "Maximum price"
// @Param        condition query string false "new or used"
// @Param        sort      query string false "Sort by: price, created_at, sold_count"
// @Param        order     query string false "asc or desc"
// @Success      200  {object}  models.PaginatedResponse
// @Router       /api/products [get]
func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	categoryID := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	condition := c.Query("condition")
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := database.DB.Model(&models.Product{}).Where("status = ?", "active")

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if condition != "" {
		query = query.Where("condition = ?", condition)
	}

	var total int64
	query.Count(&total)

	validSorts := map[string]bool{"price": true, "created_at": true, "sold_count": true}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}
	if order != "asc" {
		order = "desc"
	}

	var products []models.Product
	query.Preload("Category").Preload("Seller").
		Offset((page - 1) * limit).Limit(limit).
		Order(sortBy + " " + order).
		Find(&products)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	})
}

// GetProduct godoc
// @Summary      Get product
// @Description  Get product by ID
// @Tags         products
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  map[string]string
// @Router       /api/products/{id} [get]
func GetProduct(c *gin.Context) {
	var product models.Product
	if err := database.DB.Preload("Category").Preload("Seller.Role").Preload("Reviews.User").
		First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary      Create product
// @Description  Create a new product listing
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Product data"
// @Success      201  {object}  models.Product
// @Failure      400  {object}  map[string]string
// @Router       /api/products [post]
func CreateProduct(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var body struct {
		Title        string   `json:"title" binding:"required"`
		Description  string   `json:"description"`
		Price        float64  `json:"price" binding:"required,gt=0"`
		Stock        int      `json:"stock" binding:"required,gte=0"`
		CategoryID   uint     `json:"category_id" binding:"required"`
		Images       []string `json:"images"`
		Condition    string   `json:"condition" binding:"required,oneof=new used"`
		FreeShipping bool     `json:"free_shipping"`
		Location     string   `json:"location"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		SellerID:     user.ID,
		CategoryID:   body.CategoryID,
		Title:        body.Title,
		Description:  body.Description,
		Price:        body.Price,
		Stock:        body.Stock,
		Images:       body.Images,
		Condition:    body.Condition,
		Status:       "active",
		FreeShipping: body.FreeShipping,
		Location:     body.Location,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	database.DB.Preload("Category").Preload("Seller").First(&product, product.ID)
	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update a product listing
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Param        request body object true "Product data"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]string
// @Router       /api/products/{id} [put]
func UpdateProduct(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var product models.Product
	if err := database.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.SellerID != user.ID && user.Role.Name != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var body struct {
		Title        string   `json:"title"`
		Description  string   `json:"description"`
		Price        float64  `json:"price"`
		Stock        int      `json:"stock"`
		CategoryID   uint     `json:"category_id"`
		Images       []string `json:"images"`
		Condition    string   `json:"condition"`
		Status       string   `json:"status"`
		FreeShipping bool     `json:"free_shipping"`
		Location     string   `json:"location"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if body.Title != "" {
		updates["title"] = body.Title
	}
	if body.Description != "" {
		updates["description"] = body.Description
	}
	if body.Price > 0 {
		updates["price"] = body.Price
	}
	if body.Stock >= 0 {
		updates["stock"] = body.Stock
	}
	if body.CategoryID > 0 {
		updates["category_id"] = body.CategoryID
	}
	if len(body.Images) > 0 {
		updates["images"] = body.Images
	}
	if body.Condition != "" {
		updates["condition"] = body.Condition
	}
	if body.Status != "" {
		updates["status"] = body.Status
	}
	updates["free_shipping"] = body.FreeShipping
	if body.Location != "" {
		updates["location"] = body.Location
	}

	database.DB.Model(&product).Updates(updates)
	database.DB.Preload("Category").Preload("Seller").First(&product, product.ID)
	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Soft delete a product
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var product models.Product
	if err := database.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.SellerID != user.ID && user.Role.Name != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	database.DB.Model(&product).Update("status", "inactive")
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// GetMyProducts godoc
// @Summary      Get my products
// @Description  Get products of current seller
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Product
// @Router       /api/products/mine [get]
func GetMyProducts(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var products []models.Product
	database.DB.Where("seller_id = ?", user.ID).Preload("Category").Find(&products)
	c.JSON(http.StatusOK, products)
}
