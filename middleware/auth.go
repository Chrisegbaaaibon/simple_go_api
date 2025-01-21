package middleware

import (
	"authsystems/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.JSON(http.StatusUnAuthorized, gin.H{"error": "Authorization header is required!"})
			c.Abort()
			return 
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer"){
			c.JSON(http.StatusUnAuthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return 
		}

		claims, err := utils.ValidateJWT(parts[1], config.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnAuthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_email", claims.Email)
		c.Next()
	}
}