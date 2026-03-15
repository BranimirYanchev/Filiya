package handler

import (
	"strconv"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/helper"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RetrieveRoles retrieves a paginated list of roles
// @Summary Get all roles
// @Description Returns a paginated list of all roles in the system (requires moderator:read permission)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of roles to return (default: 20)" default(20)
// @Param offset query int false "Number of roles to skip (default: 0)" default(0)
// @Success 200 {object} map[string]interface{} "List of roles"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/roles [get]
func RetrieveRoles(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit <= 0 || err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0")) // default: 0
	if err != nil || offset < 0 {
		offset = 0
	}

	roles, err := service.GetRoles(limit, offset)
	if err != nil {
		code, msg := respondError(err)

		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": roles,
	})
}

// RetrieveUserRole retrieves the role of the currently authenticated user
// @Summary Get current user role
// @Description Returns the role information for the currently authenticated user
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User role information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/roles/me [get]
func RetrieveUserRole(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	c.JSON(200, gin.H{
		"data": user.Role,
	})
}

// CreateRole creates a new role with permissions
// @Summary Create a new role
// @Description Creates a new role with specified permissions (requires moderator:update permission)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body object true "Role creation data" example({"name":"Editor","permissions":["user:read","user:update"]})
// @Success 201 {object} map[string]interface{} "Created role"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/roles [post]
func CreateRole(c *gin.Context) {
	var roleInput struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}

	if err := c.BindJSON(&roleInput); err != nil || strings.TrimSpace(roleInput.Name) == "" || len(roleInput.Name) == 0 {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	role := model.Role{
		Name: roleInput.Name,
	}

	if err := service.CreateRole(&role, roleInput.Permissions...); err != nil {
		code, msg := respondError(err)

		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": role,
	})
}

// DeleteRole deletes a role
// @Summary Delete a role
// @Description Permanently deletes a role from the system (requires admin:access permission)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Deletion confirmation"
// @Failure 400 {object} map[string]string "Invalid role ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/roles/{id} [delete]
func DeleteRole(c *gin.Context) {
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

	role := model.Role{
		ID: id,
	}
	if err := service.DeleteRole(&role); err != nil {
		code, msg := respondError(err)

		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"role": role,
		},
	})
}

// ChangeUserRole
// @Summary      Change a user's role
// @Description  Changes the role of a target user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id         path      int  true  "User ID"
// @Param        input      body      map[string]interface{} true "Role update payload"
// @Success      200  {object}  map[string]interface{}  "Role changed successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid request or role does not exist"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router      /users/:id/role [put]

func ChangeUserRole(c *gin.Context) {
	idParam := c.Param("id")
	id, cErr := strconv.Atoi(idParam)
	if cErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request|invalid path parameter",
		})

		return
	}

	var target model.User
	if err := database.DbConnection.Where("id=?", id).First(&target).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "target does not exist",
			"data": gin.H{
				"id": idParam,
			},
		})

		return
	}

	var input struct {
		RoleName string `json:"role_name"`
	}

	if err := c.BindJSON(&input); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request|bad request body",
		})
		return
	}

	ok, role := helper.DoesRoleExist(input.RoleName)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "role does not exist",
			"data": gin.H{
				"role_name": input.RoleName,
			},
		})
		return
	}

	if err := database.DbConnection.Model(&target).Set("Role", role).Preload("Role").Error; err != nil {
		log.Error("error while updating target role")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while updating user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data": gin.H{
			"info":         "role changed",
			"target_email": target.Email,
			"target_role":  target.Role.Name,
		},
	})
}
