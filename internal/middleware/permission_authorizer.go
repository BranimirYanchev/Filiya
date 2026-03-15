package middleware

import (
	"net/http"

	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func AuthorizePermissions(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userObj, _ := c.Get("user")
		userPermissions := userObj.(model.User).Role.Permissions

		for i, perm := range userPermissions {
			if perm.Name != permission && i == len(userPermissions)-1 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"authorized": false,
				})
				c.Abort()

				return
			} else if perm.Name == permission {
				break
			}
		}
	}
}
