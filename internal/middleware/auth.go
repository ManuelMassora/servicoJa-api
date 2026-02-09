package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth(jwtService *JwtService) gin.HandlerFunc {
	return func(context *gin.Context) {
		const Bearer_schema = "Bearer "
		header := context.GetHeader("Authorization")
		if header == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		if !strings.HasPrefix(header, Bearer_schema) || len(header) < len(Bearer_schema) {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		tokenString := header[len(Bearer_schema):]

		claim, isValid := jwtService.ValidateToken(tokenString)

		if !isValid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		context.Set("userID", claim.Sum)
		context.Set("userRoles", claim.Role)
		context.Next()
	}
}