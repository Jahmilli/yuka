package routers

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

var (
	g errgroup.Group
)

func NewRouterOptions(logger *zap.Logger, db *gorm.DB) RouterOptions {
	return RouterOptions{
		Logger: logger,
		Db:     db,
	}
}

func Run(ctx context.Context, routerOptions *RouterOptions) error {
	apiRouter := setupApiRouter(ctx, routerOptions)
	tunnelRouter := setupTunnelRouter(ctx, routerOptions)

	g.Go(func() error {
		return apiRouter.ListenAndServe()
	})

	g.Go(func() error {
		return tunnelRouter.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func setupApiRouter(ctx context.Context, routerOptions *RouterOptions) *http.Server {
	r := gin.New()

	r.Use(ginzap.GinzapWithConfig(routerOptions.Logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			// SkipPaths:  []string{"/no_log"},
		},
	))
	r.Use(ginzap.RecoveryWithZap(routerOptions.Logger, true))
	// TODO: Add this back in later when we want to support authentication
	// r.Use(gin.BasicAuth(gin.Accounts{
	// 	os.Getenv("HTTP_USERNAME"): os.Getenv("HTTP_PASSWORD"),
	// 	"test":                     "test",
	// 	// Can add more users here if you want
	// }))

	v1 := r.Group("/v1")

	// Users
	userHandler := handlers.NewUserHandler(routerOptions.Logger, routerOptions.Db)
	v1.POST("/users", createUser(userHandler))
	v1.GET("/users/:id", getUser(userHandler))
	v1.PUT("/users/:id", updateUser(userHandler))
	v1.DELETE("/users/:id", deleteUser(userHandler))

	// Setup websockets
	streamingHandler := handlers.NewStreamingHandler(routerOptions.Logger, routerOptions.Db)
	r.GET("/ws", initializeStream(streamingHandler))

	// Misc
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	return &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

}

func setupTunnelRouter(ctx context.Context, routerOptions *RouterOptions) *http.Server {
	r := gin.New()

	r.Use(ginzap.GinzapWithConfig(routerOptions.Logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
		},
	))
	r.Use(ginzap.RecoveryWithZap(routerOptions.Logger, true))

	tunnelHandler := handlers.NewTunnelHandler(routerOptions.Logger, routerOptions.Db)
	r.Any("/*tunnelPath", tunnelRequest(tunnelHandler))

	return &http.Server{
		Addr:         ":8081",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// g.Go(func() error {
	// 	return server.ListenAndServe()
	// })

}
