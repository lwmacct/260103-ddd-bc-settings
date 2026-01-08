package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	settingsconfig "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/config"
)

// Public Settings 模块公开路由（无需认证）。
func Public(settingHandler *handler.SettingHandler, cfg *settingsconfig.Config) []routes.Route {
	// 从 admin 基础路径派生 public 路径
	// 例如：/api/admin/settings → /api/public/settings
	base := strings.Replace(cfg.HTTP.BasePath, "/admin/", "/public/", 1)

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
