package main

import (
	"7hunt-be-rest-api/auth"
	"7hunt-be-rest-api/internal/core/services"
	"7hunt-be-rest-api/internal/handlers"
	middleware "7hunt-be-rest-api/internal/middlewares"
	"7hunt-be-rest-api/internal/repositories"
	"7hunt-be-rest-api/utils"
	"7hunt-be-rest-api/validator"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; assuming environment variables are set")
	}

	os.Setenv("TZ", "Asia/Bangkok")
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("user_management")
	logger := utils.NewLogger("user-service")
	passwordHasher := utils.NewBcryptHasher(10)
	authManager := auth.NewJWTManager(os.Getenv("JWT_SECRET_KEY"))
	validator := validator.NewValidator()

	userRepo := repositories.NewUserRepository(db, "users")
	userService := services.NewUserService(userRepo, logger, passwordHasher, authManager)
	userHandler := handlers.NewUserHandler(userService, validator)

	router := gin.Default()
	router.Use(middleware.LogginMiddleware())

	auth := router.Group("/auth")
	auth.POST("/register", userHandler.RegisterUser)
	auth.POST("/login", userHandler.LoginUser)

	api := router.Group("/api")
	users := api.Group("/users")

	protectedUserRoute := users.Group("/")
	protectedUserRoute.Use(middleware.AuthMiddleware(authManager))
	{
		protectedUserRoute.GET("", userHandler.GetUsers)
		protectedUserRoute.GET("/:userId", userHandler.GetUserByID)
		protectedUserRoute.PUT("/:userId", userHandler.UpdateUser)
		protectedUserRoute.DELETE("/:userId", userHandler.DeleteUser)
	}

	taskCtx, cancelTasks := context.WithCancel(context.Background())
	defer cancelTasks()

	go func(ctx context.Context) {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		log.Println("Background task: User counter started")
		for {
			select {
			case <-ticker.C:
				queryCtx, queryCancel := context.WithTimeout(ctx, 10*time.Second)
				count, err := userService.Count(queryCtx)
				if err != nil {
					log.Printf("Error counting users: %v", err)
				} else {
					log.Printf("Current user count: %d", count)
				}
				queryCancel()
			case <-ctx.Done():
				log.Println("Background task: User counter stopped")
				return
			}
		}
	}(taskCtx)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	cancelTasks()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Closing MongoDB connection...")
	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}

	log.Println("Server exiting")
}
