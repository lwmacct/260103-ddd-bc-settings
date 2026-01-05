// Package bootstrap 提供 Gin Engine 和 HTTP Server 启动逻辑。
package bootstrap

import (
	"io"
	"log/slog"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/middleware"
)

func init() {
	// 设置 Gin 为 Release 模式，禁用 debug 日志（必须在包初始化时设置）
	gin.SetMode(gin.ReleaseMode)
	// 禁用 Gin 的默认日志输出
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}

// NewEngine 创建 Gin Engine，注册全局中间件
func NewEngine() *gin.Engine {
	engine := gin.New()

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// alphanumhyphen: 字母、数字、连字符、下划线
		if err := v.RegisterValidation("alphanumhyphen", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			for _, r := range value {
				if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
					return false
				}
			}
			return true
		}); err != nil {
			slog.Warn("failed to register alphanumhyphen validation", "err", err)
		}
		// loweralphanumhyphen: 小写字母、数字、连字符、下划线
		if err := v.RegisterValidation("loweralphanumhyphen", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			for _, r := range value {
				if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
					return false
				}
			}
			return true
		}); err != nil {
			slog.Warn("failed to register loweralphanumhyphen validation", "err", err)
		}
	}

	// 全局中间件
	engine.Use(gin.Recovery())         // Panic 恢复
	engine.Use(middleware.Logger())    // 请求日志（基于 slog）
	engine.Use(middleware.RequestID()) // 请求 ID
	engine.Use(middleware.CORS())      // CORS 支持

	return engine
}
