package handler

import (
	"net/http"
	"strconv"

	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// SearchUsersHandler searches for users
// @Summary Search users
// @Description Searches for users by email, full name, or bio
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of results to return (default: 20, max: 100)" default(20)
// @Param offset query int false "Number of results to skip (default: 0)" default(0)
// @Success 200 {array} model.User "List of matching users"
// @Failure 400 {object} map[string]string "Invalid query or parameters"
// @Router /api/search/users [get]
func SearchUsersHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter 'q' is required",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := service.SearchUsers(query, limit, offset)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// SearchPostsHandler searches for posts
// @Summary Search posts
// @Description Searches for posts by title or content
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of results to return (default: 20, max: 100)" default(20)
// @Param offset query int false "Number of results to skip (default: 0)" default(0)
// @Success 200 {array} model.Post "List of matching posts"
// @Failure 400 {object} map[string]string "Invalid query or parameters"
// @Router /api/search/posts [get]
func SearchPostsHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter 'q' is required",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	var user *model.User
	if u, exists := c.Get("user"); exists {
		userObj, ok := u.(model.User)
		if ok {
			user = &userObj
		}
	}

	posts, err := service.SearchPosts(query, limit, offset, user)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

// SearchCategoriesHandler searches for categories
// @Summary Search categories
// @Description Searches for categories by name
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of results to return (default: 20, max: 100)" default(20)
// @Param offset query int false "Number of results to skip (default: 0)" default(0)
// @Success 200 {array} model.Category "List of matching categories"
// @Failure 400 {object} map[string]string "Invalid query or parameters"
// @Router /api/search/categories [get]
func SearchCategoriesHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter 'q' is required",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	categories, err := service.SearchCategories(query, limit, offset)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}
