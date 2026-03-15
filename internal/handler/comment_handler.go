package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// GetPostComments retrieves all comments for a specific post
// @Summary Get post comments
// @Description Returns all comments associated with a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} model.Comment "List of comments"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Router /api/posts/{id}/comments [get]
func GetPostComments(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path parameter",
		})
		return
	}
	post := model.Post{
		ID: postId,
	}

	comments, err := service.GetPostComments(&post)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// CreateComment creates a new comment on a post
// @Summary Create a comment
// @Description Creates a new comment on a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Param input body object true "Comment content" example({"content":"Great post!"})
// @Success 201 {object} map[string]interface{} "Created comment"
// @Failure 400 {object} map[string]string "Invalid input or post ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/posts/{id}/comments [post]
func CreateComment(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path parameter",
		})
		return
	}
	var body struct {
		Content string `json:"content"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid data",
		})
		return
	}

	user, _ := c.Get("user")

	comment := model.Comment{
		PostID:   postId,
		Content:  body.Content,
		AuthorID: user.(model.User).ID,
	}

	if err := service.SaveComment(&comment); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": comment,
	})
}

// UpdateComment updates an existing comment
// @Summary Update a comment
// @Description Updates the content of an existing comment (requires user:update permission)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Param input body object true "Updated comment content" example({"content":"Updated comment text"})
// @Success 200 {object} map[string]interface{} "Updated comment"
// @Failure 400 {object} map[string]string "Invalid input or comment ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/comments/{id} [put]
func UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid path parameter",
			"data": gin.H{
				"id": c.Param("id"),
			},
		})

		return
	}

	var body struct {
		Content string `json:"content"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid data",
		})
		return
	}

	comment := model.Comment{ID: id}
	updatedComment := model.Comment{Content: body.Content, UpdatedAt: time.Now()}
	if err := service.UpdateComment(&comment, &updatedComment); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  comment,
	})
}

// DeleteComment deletes a comment
// @Summary Delete a comment
// @Description Permanently deletes a comment (requires user:update permission)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]string "Deletion confirmation"
// @Failure 400 {object} map[string]string "Invalid comment ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid path parameter",
			"data": gin.H{
				"id": c.Param("id"),
			},
		})

		return
	}

	comment := model.Comment{ID: id}
	if err := service.DeleteComment(&comment); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "successfully deleted comment",
	})
}

// LikeComment likes or unlikes a comment
// @Summary Like/unlike a comment
// @Description Toggles the like status for a comment. If already liked, removes the like (requires user:update permission)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]interface{} "Like count"
// @Failure 400 {object} map[string]string "Invalid comment ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/comments/{id}/like [post]
func LikeComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid path parameter",
			"data": gin.H{
				"id": c.Param("id"),
			},
		})

		return
	}

	u, _ := c.Get("user")
	user, _ := u.(model.User)

	comment := model.Comment{ID: id}

	likes, err := service.LikeComment(&comment, &user)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"likes": likes,
		},
	})
}

// GetComment retrieves a specific comment by ID
// @Summary Get comment by ID
// @Description Returns detailed information about a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]interface{} "Comment details"
// @Failure 400 {object} map[string]string "Invalid comment ID"
// @Failure 404 {object} map[string]string "Comment not found"
// @Router /api/comments/{id} [get]
func GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path parameter",
		})
		return
	}

	comment := model.Comment{
		ID: id,
	}

	if err := service.GetComment(&comment); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comment,
	})
}
