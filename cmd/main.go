package main

import (
	"context"
	"fmt"
	_ "go-sosmed/docs"
	"go-sosmed/internal/comment"
	"go-sosmed/internal/follow"
	"go-sosmed/internal/like"
	"go-sosmed/internal/post"
	"go-sosmed/internal/report"
	"go-sosmed/internal/session"
	"go-sosmed/internal/user"
	cloudinaryPkg "go-sosmed/pkg/cloudinary"
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GO-SOSMED API
// @version 1.0
// @description API for Golang Social Media Application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@go-sosmed.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and session token.

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
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

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
		&post.Post{},
		&like.Like{},
		&follow.Follow{},
		&comment.Comment{},
		&report.Report{},
		&session.Session{},
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// === Health Check ===
	r.GET("/health", func(c *gin.Context) {
    db := config.GetDB()

    sqlDB, err := db.DB()
    if err != nil {
        c.JSON(500, gin.H{"status": "error"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    if err := sqlDB.PingContext(ctx); err != nil {
        c.JSON(503, gin.H{"status": "error"})
        return
    }

    c.JSON(200, gin.H{"status": "ok"})
})

	r.Use(middlewares.GinErrorHandler())

	//seeder
	user.SeedAdminUser()

	var cloudinaryService *cloudinaryPkg.Service
	if cfg.CloudinaryCloudName != "" && cfg.CloudinaryAPIKey != "" && cfg.CloudinaryAPISecret != "" {
		cld, cldErr := cloudinaryPkg.NewService(
			cfg.CloudinaryCloudName,
			cfg.CloudinaryAPIKey,
			cfg.CloudinaryAPISecret,
			cfg.CloudinaryFolder,
		)
		if cldErr != nil {
			log.Printf("Warning: Cloudinary not configured: %v", cldErr)
		} else {
			cloudinaryService = cld
		}
	} else {
		log.Println("Cloudinary not configured, uploads will fail")
	}

	userRepo := user.NewRepository(db)
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo, cfg)
	userService := user.NewService(userRepo, cfg, sessionService)
	userController := user.NewController(userService, cfg, cloudinaryService)
	user.SetupRoute(r, userController, cfg, sessionService, cloudinaryService)

	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo, cloudinaryService)
	postController := post.NewController(postService)
	post.SetupPostRoute(r, postController, cfg, sessionService, cloudinaryService)

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)
	likeController := like.NewController(likeService)
	like.SetupLikeRoute(r, likeController, cfg, sessionService)

	followRepo := follow.NewRepository(db)
	followService := follow.NewService(followRepo)
	followController := follow.NewController(followService)
	follow.SetupFollowRoute(r, followController, cfg, sessionService)

	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo, postRepo)
	commentController := comment.NewController(commentService)
	comment.SetupCommentRoute(r, commentController, cfg, sessionService)

	reportRepo := report.NewRepository(db)
	reportService := report.NewService(reportRepo)
	reportController := report.NewController(reportService)
	report.SetupRoute(r, reportController, cfg, sessionService)

	// 404 Not Found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	// clean.CleanupUnusedUploads(db, "./uploads")

	// === Start Server ===
	log.Printf("Server running on port %s", cfg.Port)
	log.Printf("Local: http://localhost:%s", cfg.Port)
	log.Printf("Environment: %s", cfg.NodeEnv)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
