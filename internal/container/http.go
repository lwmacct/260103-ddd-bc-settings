package container

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/internal/config"
	ginroutes "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"
)

// AllRoutesParams 注入 Settings 模块路由
type AllRoutesParams struct {
	fx.In

	SettingsRoutes []ginroutes.Route `name:"settings"`
}

// AllRoutes 返回 Settings 模块路由
func AllRoutes(p AllRoutesParams) []ginroutes.Route {
	return p.SettingsRoutes
}

// HTTPModule 提供 HTTP 处理器和路由注册。
var HTTPModule = fx.Module("http",
	fx.Provide(
		gin.Default,
		AllRoutes,
	),
	fx.Invoke(registerRoutes),
	fx.Invoke(startHTTPServer),
)

// registerRoutes 注册 Settings 模块的路由。
func registerRoutes(
	r *gin.Engine,
	allRoutes []ginroutes.Route,
) {
	// 添加自定义 panic recovery 中间件（在 Gin 默认 recovery 之前）
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Handler panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(debug.Stack()))
				c.JSON(500, gin.H{"code": 500, "message": "Internal Server Error", "error": fmt.Sprintf("%v", err)})
				c.Abort()
			}
		}()
		c.Next()
	})

	slog.Info("Registering routes", "count", len(allRoutes))

	// 注册路由到 Gin Engine
	for _, route := range allRoutes {
		// 将 OpenAPI 风格路径参数 {param} 转换为 Gin 风格 :param
		ginPath := ginroutes.ToGinPath(route.Path)

		switch route.Method {
		case ginroutes.GET:
			r.GET(ginPath, route.Handlers...)
		case ginroutes.POST:
			r.POST(ginPath, route.Handlers...)
		case ginroutes.PUT:
			r.PUT(ginPath, route.Handlers...)
		case ginroutes.DELETE:
			r.DELETE(ginPath, route.Handlers...)
		case ginroutes.PATCH:
			r.PATCH(ginPath, route.Handlers...)
		}
	}
}

// startHTTPServer 启动 HTTP 服务器。
func startHTTPServer(lc fx.Lifecycle, r *gin.Engine, cfg *config.Config) {
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("HTTP server starting", "addr", cfg.Server.Addr)

			// 在 goroutine 中启动服务器，避免阻塞 Fx 启动流程
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("HTTP server error", "err", err)
				}
			}()

			slog.Info("HTTP server started successfully", "addr", cfg.Server.Addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("HTTP server shutting down", "addr", cfg.Server.Addr)

			// 优雅关闭，等待现有请求完成（最多 5 秒）
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				slog.Error("HTTP server shutdown error", "err", err)
				return err
			}

			slog.Info("HTTP server stopped successfully")
			return nil
		},
	})
}
