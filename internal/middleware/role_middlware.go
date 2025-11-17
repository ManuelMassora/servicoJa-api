package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func HasRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesInterface, exists := c.Get("userRoles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User roles not found in context"})
			return
		}

		var userRoles []string
		switch v := rolesInterface.(type) {
		case string:
			userRoles = strings.Split(v, ",")
		case []string:
			userRoles = v
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user roles type in context"})
			return
		}

		isAuthorized := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range userRoles {
				if strings.EqualFold(strings.TrimSpace(userRole), requiredRole) {
					isAuthorized = true
					break
				}
			}
			if isAuthorized {
				break
			}
		}

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Calma ai hackerzinho: Previlégios insuficientes"})
			return
		}
		c.Next()
	}
}