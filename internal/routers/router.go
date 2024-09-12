package routers

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"yuka/internal/handlers"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterOptions struct {
	Logger *zap.Logger
	Db     *gorm.DB
}

func NewRouterOptions(logger *zap.Logger, db *gorm.DB) RouterOptions {
	return RouterOptions{
		Logger: logger,
		Db:     db,
	}
}

func SetupRouter(routerOptions *RouterOptions) {
	r := gin.New()

	r.Use(ginzap.GinzapWithConfig(routerOptions.Logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			// SkipPaths:  []string{"/no_log"},
		},
	))
	r.Use(ginzap.RecoveryWithZap(routerOptions.Logger, true))
	r.Use(gin.BasicAuth(gin.Accounts{
		os.Getenv("HTTP_USERNAME"): os.Getenv("HTTP_PASSWORD"),
		"jahmilli":                 "Jsudhajekshdu123!sjdi1",
		// Can add more users here if you want
	}))

	v1 := r.Group("/v1")
	if os.Getenv("ENVIRONMENT") == "local" {
		// Hack for now to allow local development as we don't have a proper ingress controller
		v1 = r.Group("/api/v1")
	}

	// Users
	userHandler := handlers.NewUserHandler(routerOptions.Logger, routerOptions.Db)
	// v1.GET("/users", getUsers(userHandler))
	v1.POST("/users", createUser(userHandler))
	v1.GET("/users/:id", getUser(userHandler))
	v1.PUT("/users/:id", updateUser(userHandler))
	v1.DELETE("/users/:id", deleteUser(userHandler))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})
	r.GET("/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})
	r.Run(":8080")
}
