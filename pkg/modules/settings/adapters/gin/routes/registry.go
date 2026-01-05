package routes

import (
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	settingsconfig "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
	sharedRoutes "github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"
)

// AllRouteBindings 返回 Settings 模块的所有路由绑定。
func AllRouteBindings(
	settingHandler *handler.SettingHandler,
	cfg *settingsconfig.Config,
) []sharedRoutes.Route {
	return Admin(settingHandler, cfg)
}
