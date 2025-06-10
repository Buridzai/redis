package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsIface, exists := c.Get("user_perms")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			return
		}
		perms := permsIface.([]string)
		for _, p := range perms {
			if p == requiredPerm {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}
