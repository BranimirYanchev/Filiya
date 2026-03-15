package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Marionvd/filia-project-backend/internal/service"

	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

// GetCategories retrieves a paginated list of categories
// @Summary Get all categories
// @Description Returns a paginated list of all categories in the system
// @Tags categories
// @Accept json
// @Produce json
// @Param limit query int false "Number of categories to return (default: 10, max: 200)" default(10)
// @Param offset query int false "Number of categories to skip (default: 0)" default(0)
// @Success 200 {object} map[string]interface{} "List of categories"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/categories [get]
func GetCategories(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit <= 0 || limit > 200 || err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	categories, err := service.GetCategories(limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error|failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": categories,
	})
}

// GetCategory retrieves a specific category by ID
// @Summary Get category by ID
// @Description Returns detailed information about a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{} "Category details"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Router /api/categories/{id} [get]
func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path parameter",
			"data": gin.H{
				"id": c.Param("id"),
			},
		})

		return
	}

	category, err := service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unexisting category",
			"data": gin.H{
				"id": c.Param("id"),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// FindCategoryByName finds a category by name
// @Summary Find category by name
// @Description Searches for a category by its name (authenticated endpoint)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body object true "Category name" example({"name":"Technology"})
// @Success 200 {object} map[string]interface{} "Category details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Category not found"
// @Router /api/categories/name [get]
func FindCategoryByName(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	category, err := service.GetCategoryByName(input.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "unexisting category",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// GetCategoryPosts retrieves all posts in a specific category
// @Summary Get posts by category
// @Description Returns a paginated list of posts belonging to a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param limit query int false "Number of posts to return (default: 10, max: 200)" default(10)
// @Param offset query int false "Number of posts to skip (default: 0)" default(0)
// @Success 200 {object} map[string]interface{} "Category and its posts"
// @Failure 400 {object} map[string]string "Invalid category ID or query parameters"
// @Router /api/categories/{id}/posts [get]
func GetCategoryPosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path param",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit <= 0 || limit > 200 || err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	category := model.Category{
		ID: id,
	}
	err = service.GetCategoryPosts(&category, limit, offset)
	if err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"category": category.Name,
			"posts":    category.Posts,
		},
	})
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Creates a new category with optional parent category (requires moderator:create permission)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.CreateCategoryInput true "Category creation data"
// @Success 201 {object} map[string]interface{} "Created category"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/categories [post]
func CreateCategory(c *gin.Context) {
	var input model.CreateCategoryInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	var category = model.Category{
		Name:     input.Name,
		ParentID: input.ParentID,
	}

	if err := service.SaveCategory(&category); err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": category,
	})
}

// UpdateCategoryHandler updates an existing category
// @Summary Update a category
// @Description Updates category name and/or parent (requires moderator:update permission)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param input body model.UpdateCategoryInput true "Updated category data"
// @Success 200 {object} map[string]interface{} "Updated category"
// @Failure 400 {object} map[string]string "Invalid input or category ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/categories/{id} [put]
func UpdateCategoryHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid path parameter",
		})
		return
	}

	var input model.UpdateCategoryInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	category := model.Category{
		ID: id,
	}
	updatedCategory := model.Category{
		Name:     input.Name,
		ParentID: input.ParentID,
	}

	if err = service.UpdateCategory(&category, &updatedCategory); err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Permanently deletes a category and its related usages (requires moderator:delete permission)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "Deletion confirmation"
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid path parameter",
		})
		return
	}

	category := model.Category{
		ID: id,
	}
	if err := service.DeleteCategory(&category); err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			log.Debugf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "deleted category and it's related usages",
	})
}
