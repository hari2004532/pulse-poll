package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/harih/pulse-poll/backend/config"
	"github.com/harih/pulse-poll/backend/handlers"
	"github.com/harih/pulse-poll/backend/middleware"
	"github.com/harih/pulse-poll/backend/services"
	ws "github.com/harih/pulse-poll/backend/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg := config.Load()

	if cfg.MongoURI == "" {
		log.Fatal("MONGODB_URI is not configured")
	}

	if cfg.RedisAddr == "" {
		log.Fatal("REDIS_ADDR is not configured")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not configured")
	}

	// MongoDB
	mongoClient, err := config.ConnectMongo(cfg)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	db := mongoClient.Database(cfg.MongoDatabase)

	voteCollection := db.Collection("votes")

	_, err = voteCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "pollId", Value: 1},
				{Key: "voterId", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Fatal("failed to create vote unique index:", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println("MongoDB disconnect error:", err)
		}
	}()

	// Redis
	redisService := services.NewRedisService(
		cfg.RedisAddr,
		cfg.RedisPassword,
	)

	redisCtx, redisCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer redisCancel()

	if err := redisService.Client.Ping(redisCtx).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	// Services
	authService := services.NewAuthService(
		db,
		cfg.JWTSecret,
	)

	pollService := services.NewPollService(db.Collection("polls"), redisService)

	voteService := services.NewVoteService(
		db.Collection("votes"),
		db.Collection("polls"),
		redisService,
	)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	pollHandler := handlers.NewPollHandler(pollService)
	voteHandler := handlers.NewVoteHandler(voteService)

	// WebSocket hub
	hub := ws.NewHub(redisService.Client)
	go hub.StartRedisSubscriber(context.Background())

	// Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Public routes
	router.POST("/api/auth/signup", authHandler.Signup)
	router.POST("/api/auth/login", authHandler.Login)

	router.GET("/api/polls/:id", pollHandler.GetPoll)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	protected.POST("/polls", pollHandler.CreatePoll)
	protected.POST("/polls/:id/vote", voteHandler.Vote)

	// WebSocket
	router.GET("/ws/polls/:id", ws.HandleConnection(hub))

	log.Println("PulsePoll backend running on port", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
