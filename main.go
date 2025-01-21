package main

import (
	"log"

	"authsystem/config"
	"authsystem/handlers"
	"authsystem/middlewares"
	"authsystem/models"

	"github.com/gin-gonic/gin"
)


func main (){
	cfg, err := config.Load()
	if err != nil{
		log.Fatalf("Failed to load configuration %v", err)
	}

	//db
	db, err := models.InitDatabase(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to db %v", err)
	}

	defer db.client().Disconnect(context.Background())

	redisClient, err := modles.InitRedis(cfg.Redis)
	if err != nil{
		log.Fatalf("Failed to connect to Redis %v", err)
	}

	defer redisClient.close()

	r := gin.Default()

	setupRoutes(r, db, redisClient, cfg)

	if err := r.Run(cfg.ServerAddress); err != nil{
		log.Fatalf("Failed  to start server %v", err)
	}
}

func setupRoutes(r *gin.Engine, db *mongo.Database, redisClient *redis.Client, cfg* config.Config){
	auth := handlers.NewAuthHandler(db, redisClient, cfg)

	r.Post("/register", auth.Register)
	r.Post("/login", auth.Login)
	r.Post("/veryfy-otp", auth.verifyOtp)
	r.Post("/logout", middleware.AuthMiddleware(), auth.Logout)
	r.Get("/protected", middleware.AuthMiddleware(), handlers.ProtedctedHandler)
}