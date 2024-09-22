package routers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"yuka/internal/handlers"
	"yuka/pkg/streaming_connection"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterOptions struct {
	logger *zap.Logger
	db     *gorm.DB
}

type ApiRouterOptions struct {
	RouterOptions
	wsHandler *handlers.WsHandler
	port      int
}
type TunnelRouterOptions struct {
	RouterOptions
	port           int
	connectionPool streaming_connection.StreamingConnectionPool
}

var (
	g errgroup.Group
)

func NewRouterOptions(logger *zap.Logger, db *gorm.DB) RouterOptions {
	return RouterOptions{
		logger: logger,
		db:     db,
	}
}

func Run(ctx context.Context, routerOptions *RouterOptions) error {
	wsHandler := handlers.NewWsHandler(routerOptions.logger, routerOptions.db)
	connectionPool := streaming_connection.NewStreamingConnectionPool(routerOptions.logger)
	apiRouter := setupApiRouter(ctx, &ApiRouterOptions{
		RouterOptions: *routerOptions,
		wsHandler:     &wsHandler,
		port:          8080,
	})
	tunnelRouter := setupTunnelRouter(ctx, &TunnelRouterOptions{
		RouterOptions:  *routerOptions,
		port:           8081,
		connectionPool: *connectionPool,
	})

	tcpTunnel := streaming_connection.NewTcpTunnel(routerOptions.logger, 8085, connectionPool)
	g.Go(func() error {
		return tcpTunnel.Listen(ctx)
	})

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

func setupApiRouter(ctx context.Context, routerOptions *ApiRouterOptions) *http.Server {
	r := gin.New()

	r.Use(ginzap.GinzapWithConfig(routerOptions.logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			// SkipPaths:  []string{"/no_log"},
		},
	))
	r.Use(ginzap.RecoveryWithZap(routerOptions.logger, true))
	// TODO: Add this back in later when we want to support authentication
	// r.Use(gin.BasicAuth(gin.Accounts{
	// 	os.Getenv("HTTP_USERNAME"): os.Getenv("HTTP_PASSWORD"),
	// 	"test":                     "test",
	// 	// Can add more users here if you want
	// }))

	v1 := r.Group("/v1")

	// Users
	userHandler := handlers.NewUserHandler(routerOptions.logger, routerOptions.db)
	v1.POST("/users", createUser(userHandler))
	v1.GET("/users/:id", getUser(userHandler))
	v1.PUT("/users/:id", updateUser(userHandler))
	v1.DELETE("/users/:id", deleteUser(userHandler))

	// Setup websockets
	r.GET("/ws", handleWsConnection(*routerOptions.wsHandler))

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
		Addr:         fmt.Sprintf(":%v", routerOptions.port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

}

func setupTunnelRouter(ctx context.Context, routerOptions *TunnelRouterOptions) *http.Server {
	r := gin.New()

	r.Use(ginzap.GinzapWithConfig(routerOptions.logger,
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
		},
	))
	r.Use(ginzap.RecoveryWithZap(routerOptions.logger, true))

	tunnelHandler := handlers.NewTunnelHandler(routerOptions.logger, routerOptions.db, routerOptions.connectionPool)
	r.Any("/*tunnelPath", tunnelRequest(tunnelHandler))

	return &http.Server{
		Addr:         fmt.Sprintf(":%v", routerOptions.port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
