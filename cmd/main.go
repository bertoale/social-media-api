package main

import (
	"fmt"
	"go-sosmed/internal/blog"
	"go-sosmed/internal/comment"
	"go-sosmed/internal/follow"
	"go-sosmed/internal/like"
	"go-sosmed/internal/user"
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %d - %s %s %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.StatusCode,
			param.Method,
			param.Path,
			param.Latency,
		)
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.CorsOrigin},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	r.Use(middlewares.GinErrorHandler())

	// === Static Files untuk Upload ===
	r.Static("/uploads", "./uploads")

	// === Database ===
	if err := config.Connect(cfg); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// === Migrate Database ===
	db := config.GetDB()
	tables := []interface{}{
		&user.User{},
		&blog.Blog{},
		&like.Like{},
		&follow.Follow{},
		&comment.Comment{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("✅ Migrasi database berhasil.")

	// === Home Route ===
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Welcome to the REST API",
			"version":   "1.0.0",
			"timestamp": time.Now(),
		})
	})

	//seeder
	user.SeedAdminUser()
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg)
	userController := user.NewController(userService, cfg)
	user.SetupRoute(r, userController, cfg)

	blogRepo := blog.NewRepository(db)
	blogService := blog.NewService(blogRepo)
	blogController := blog.NewController(blogService)
	blog.SetupBlogRoute(r, blogController, cfg)

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)
	likeController := like.NewController(likeService)
	like.SetupLikeRoute(r, likeController, cfg)

	followRepo := follow.NewRepository(db)
	followService := follow.NewService(followRepo)
	followController := follow.NewController(followService)
	follow.SetupFollowRoute(r, followController, cfg)

	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentController := comment.NewController(commentService)
	comment.SetupCommentRoute(r, commentController, cfg)

	// 404 Not Found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	// === Start Server ===
	log.Printf("Server running on port %s", cfg.Port)
	log.Printf("Local: http://localhost:%s", cfg.Port)
	log.Printf("Environment: %s", cfg.NodeEnv)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
