package routes

import (
	"go.uber.org/fx"

	sharedRoutes "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/adapters/gin/handler"
	settingsconfig "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/config"
)

// RoutesModule 导出 Settings 模块的路由配置
var RoutesModule = fx.Module("settings.routes",
	fx.Provide(
		fx.Annotate(
			NewAllRoutes,
			fx.ResultTags(`name:"settings"`),
		),
	),
)

// NewAllRoutes 聚合所有 Settings 路由
func NewAllRoutes(
	h *handler.Handlers,
	cfg *settingsconfig.Config,
) []sharedRoutes.Route {
	var allRoutes []sharedRoutes.Route
	allRoutes = append(allRoutes, Admin(h.Setting, cfg)...)
	allRoutes = append(allRoutes, Public(h.Setting, cfg)...)
	return allRoutes
}
