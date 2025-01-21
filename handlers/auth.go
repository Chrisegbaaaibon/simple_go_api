package handlers

import (
	"net/http"
	"time"

	"authsystems/models"
	"authsystems/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-redis/redis/v8"
)

type AuthHandler struct {
	db *mongo.Database
	redis *redis.client
	config *config.Config
}

func NewAuthHandler (db *mongo.Database, redis *redis.client, cfg *config.Config) *AuthHandler{
	return &AuthHandler{
		db: db,
		redis: redis,
		config: cfg
	}
}

func (ah* AuthHandler) Register (c *gin.context){
	var user &models.user{
		isVerified: false
	}
	if err := c.ShouldBindJSON(&User); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error()})
		return
	}

	existingUser, err := ah.db.Collection("users").FindOne(c, bson.M{"email": user.email}).Decode(&user)
	if existingUser != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email exists"})
		return
	}



	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password."})
		return
	}

	user.Password = hashedPassword

	newUser, err := ah.db.Collection("users").InsertOne(c, user)
	if err != nill {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	otp, err := utils.GenerateOTP()
	if err != nil 
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Otp"})
	return

	err = ah.redis.Set(c, newUser._id, otp, time.Duration(ah.config.OTPExpiraton)*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent," "otp": otp })
}

func (ah *AuthHandler) Login(c *gin.Context){
	var loginData struct {
		Email string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := ah.db.Collection("users").FindOne(c, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil{
		c.JSON(http.StatusUnAuthorized, gin.H{"error": "user not found"})
		return
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password){
		c.JSON(http.StatusUnAuthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.isVerified {
		otp, err := utils.GenerateOTP()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Otp"})
			return
		}

		err = ah.redis.Set(c, newUser._id, otp, time.Duration(ah.config.OTPExpiraton)*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent," "otp": otp })
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

func (ah* AuthHandler) VerifyOtp(c *gin.Context){
	var otpData struct  {
		Email string `json:"email" binding:"required,email"`
		OTP string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&otpData); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stored, err := ah.redis.Get(c, otpData.Email).Result()
	if err == redis.Nil{
		c.JSON(http.StatusUnAuthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}

	token, err := utils.GenerateJWT(otpData.Email, ah.config.JWTSecret)
	if err != nil {
		c.JSON(httpStatusInternalServerError, ginH{"error": "Failed to generate token"})
		return
	}

	ah.redis.Del(c, otpData.Email)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

