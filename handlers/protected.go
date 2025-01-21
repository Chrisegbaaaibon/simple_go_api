package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProtectedHandler(c *gin.Context){
	userEmail, _ := c.Get("user_email")
	c.JSON(http.StatusOK, gin.H{"message": "This is a protected route", "user": userEmail})
}