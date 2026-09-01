package handlers

import (
	"math"
	"net/http"
	"strconv"

	"marketplace/database"
	"marketplace/models"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary      List users
// @Description  Get all users (admin only)
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200  {object}  models.PaginatedResponse
// @Router       /api/admin/users [get]
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}

	query := database.DB.Model(&models.User{})
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var users []models.User
	query.Preload("Role").Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&users)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	})
}

// UpdateUserRole godoc
// @Summary      Update user role
// @Description  Update user role (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Param        request body object true "Role data"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Router       /api/admin/users/{id}/role [put]
func UpdateUserRole(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var body struct {
		RoleName string `json:"role_name" binding:"required,oneof=admin seller buyer"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role models.Role
	if err := database.DB.Where("name = ?", body.RoleName).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	user.RoleID = role.ID
	database.DB.Save(&user)
	database.DB.Preload("Role").First(&user, user.ID)

	c.JSON(http.StatusOK, user)
}

// ToggleUserActive godoc
// @Summary      Toggle user active status
// @Description  Enable/disable user account (admin only)
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Router       /api/admin/users/{id}/toggle [put]
func ToggleUserActive(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.IsActive = !user.IsActive
	database.DB.Save(&user)
	database.DB.Preload("Role").First(&user, user.ID)

	c.JSON(http.StatusOK, user)
}

// GetDashboardStats godoc
// @Summary      Dashboard stats
// @Description  Get admin dashboard statistics
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]int64
// @Router       /api/admin/stats [get]
func GetDashboardStats(c *gin.Context) {
	var userCount, productCount, orderCount int64
	var totalRevenue float64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Product{}).Where("status = ?", "active").Count(&productCount)
	database.DB.Model(&models.Order{}).Count(&orderCount)
	database.DB.Model(&models.Order{}).Select("COALESCE(SUM(total), 0)").Where("status != ?", "cancelled").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"users":        userCount,
		"products":     productCount,
		"orders":       orderCount,
		"total_revenue": totalRevenue,
	})
}

// GetAllProducts admin godoc
// @Summary      List all products
// @Description  Get all products including inactive (admin only)
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200  {object}  models.PaginatedResponse
// @Router       /api/admin/products [get]
func GetAllProductsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}

	var total int64
	database.DB.Model(&models.Product{}).Count(&total)

	var products []models.Product
	database.DB.Preload("Category").Preload("Seller").
		Offset((page - 1) * limit).Limit(limit).
		Order("created_at DESC").
		Find(&products)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	})
}
