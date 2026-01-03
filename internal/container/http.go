package container

import (
	"github.com/gin-gonic/gin"
	ginroutes "github.com/lwmacct/260101-go-pkg-gin/pkg/routes"
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	settingsroutes "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/routes"
)

// HTTPModule 提供 HTTP 处理器和路由注册。
var HTTPModule = fx.Module("http",
	fx.Provide(
		gin.Default,
	),
	fx.Invoke(registerRoutes),
)

// registerRoutes 注册 Settings 模块的路由。
func registerRoutes(r *gin.Engine, settingHandler *handler.SettingHandler) {
	// 获取所有 Settings 路由
	allRoutes := settingsroutes.Admin(settingHandler)

	// 注册路由到 Gin Engine
	for _, route := range allRoutes {
		switch route.Method {
		case ginroutes.GET:
			r.GET(route.Path, route.Handler)
		case ginroutes.POST:
			r.POST(route.Path, route.Handler)
		case ginroutes.PUT:
			r.PUT(route.Path, route.Handler)
		case ginroutes.DELETE:
			r.DELETE(route.Path, route.Handler)
		case ginroutes.PATCH:
			r.PATCH(route.Path, route.Handler)
		}
	}
}
