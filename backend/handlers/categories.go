package handlers

import (
	"net/http"

	"marketplace/database"
	"marketplace/models"

	"github.com/gin-gonic/gin"
)

// GetCategories godoc
// @Summary      List categories
// @Description  Get all categories with children
// @Tags         categories
// @Produce      json
// @Success      200  {array}   models.Category
// @Router       /api/categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Where("parent_id IS NULL").Preload("Children").Find(&categories)
	c.JSON(http.StatusOK, categories)
}

// GetCategory godoc
// @Summary      Get category
// @Description  Get category by slug
// @Tags         categories
// @Produce      json
// @Param        slug path string true "Category slug"
// @Success      200  {object}  models.Category
// @Failure      404  {object}  map[string]string
// @Router       /api/categories/{slug} [get]
func GetCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.Where("slug = ?", c.Param("slug")).Preload("Children").First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// CreateCategory godoc
// @Summary      Create category
// @Description  Create a new category (admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Category data"
// @Success      201  {object}  models.Category
// @Failure      400  {object}  map[string]string
// @Router       /api/categories [post]
func CreateCategory(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Slug     string `json:"slug" binding:"required"`
		Image    string `json:"image"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:     body.Name,
		Slug:     body.Slug,
		Image:    body.Image,
		ParentID: body.ParentID,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}
