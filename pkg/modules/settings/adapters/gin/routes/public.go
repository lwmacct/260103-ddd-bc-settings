package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	settingsconfig "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
)

// Public Settings 模块公开路由（无需认证）。
func Public(settingHandler *handler.SettingHandler, cfg *settingsconfig.Config) []routes.Route {
	base := cfg.PublicPath

	return []routes.Route{
		{
			Handlers:    []gin.HandlerFunc{settingHandler.GetPublicSettings},
			Method:      routes.GET,
			Path:        buildPath(base, ""),
			OperationID: "public:settings:list",
			Tags:        []string{"Public"},
			Summary:     "公开配置列表",
			Audit:       false,
		},
	}
}
