package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/adapters/gin/handler"
	settingsconfig "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/config"
)

// Admin Settings 模块管理员路由。
func Admin(
	settingHandler *handler.SettingHandler,
	cfg *settingsconfig.Config,
) []routes.Route {
	base := cfg.API.AdminPath

	var allRoutes []routes.Route

	// ==================== 配置分类 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Handlers:    []gin.HandlerFunc{settingHandler.GetCategories},
			Method:      routes.GET,
			Path:        buildPath(base, "/categories"),
			OperationID: "admin:setting:categories:list",
			Tags:        []string{"admin-settings"},
			Summary:     "配置分类列表",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.GetCategory},
			Method:      routes.GET,
			Path:        buildPath(base, "/categories/{id}"),
			OperationID: "admin:setting:categories:get",
			Tags:        []string{"admin-settings"},
			Summary:     "配置分类详情",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.CreateCategory},
			Method:      routes.POST,
			Path:        buildPath(base, "/categories"),
			OperationID: "admin:setting:categories:create",
			Tags:        []string{"admin-settings"},
			Summary:     "创建配置分类",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.UpdateCategory},
			Method:      routes.PUT,
			Path:        buildPath(base, "/categories/{id}"),
			OperationID: "admin:setting:categories:update",
			Tags:        []string{"admin-settings"},
			Summary:     "更新配置分类",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.DeleteCategory},
			Method:      routes.DELETE,
			Path:        buildPath(base, "/categories/{id}"),
			OperationID: "admin:setting:categories:delete",
			Tags:        []string{"admin-settings"},
			Summary:     "删除配置分类",
			Audit:       true,
		},
	}...)

	// ==================== 系统配置 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Handlers:    []gin.HandlerFunc{settingHandler.BatchUpdateSettings},
			Method:      routes.POST,
			Path:        buildPath(base, "/batch"),
			OperationID: "admin:settings:batch:update",
			Tags:        []string{"admin-settings"},
			Summary:     "批量更新配置",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.CreateSetting},
			Method:      routes.POST,
			Path:        buildPath(base, ""),
			OperationID: "admin:settings:create",
			Tags:        []string{"admin-settings"},
			Summary:     "创建配置",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.GetSettings},
			Method:      routes.GET,
			Path:        buildPath(base, ""),
			OperationID: "admin:settings:list",
			Tags:        []string{"admin-settings"},
			Summary:     "配置列表",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.GetSetting},
			Method:      routes.GET,
			Path:        buildPath(base, "/{key}"),
			OperationID: "admin:settings:get",
			Tags:        []string{"admin-settings"},
			Summary:     "配置详情",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.UpdateSetting},
			Method:      routes.PUT,
			Path:        buildPath(base, "/{key}"),
			OperationID: "admin:settings:update",
			Tags:        []string{"admin-settings"},
			Summary:     "更新配置",
			Audit:       true,
		},
		{
			Handlers:    []gin.HandlerFunc{settingHandler.DeleteSetting},
			Method:      routes.DELETE,
			Path:        buildPath(base, "/{key}"),
			OperationID: "admin:settings:delete",
			Tags:        []string{"admin-settings"},
			Summary:     "删除配置",
			Audit:       true,
		},
	}...)

	return allRoutes
}
