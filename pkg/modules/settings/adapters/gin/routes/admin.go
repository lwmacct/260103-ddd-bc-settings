package routes

import (
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
)

// Admin Settings 模块管理员路由。
func Admin(
	settingHandler *handler.SettingHandler,
) []routes.Route {
	var allRoutes []routes.Route

	// ==================== 配置分类 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings/categories",
			Handler:   settingHandler.GetCategories,
			Operation: "admin:setting:categories:list",
			Tags:      "Admin - Settings",
			Summary:   "配置分类列表",
		},
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings/categories/:id",
			Handler:   settingHandler.GetCategory,
			Operation: "admin:setting:categories:get",
			Tags:      "Admin - Settings",
			Summary:   "配置分类详情",
		},
		{
			Method:    routes.POST,
			Path:      "/api/admin/settings/categories",
			Handler:   settingHandler.CreateCategory,
			Operation: "admin:setting:categories:create",
			Tags:      "Admin - Settings",
			Summary:   "创建配置分类",
		},
		{
			Method:    routes.PUT,
			Path:      "/api/admin/settings/categories/:id",
			Handler:   settingHandler.UpdateCategory,
			Operation: "admin:setting:categories:update",
			Tags:      "Admin - Settings",
			Summary:   "更新配置分类",
		},
		{
			Method:    routes.DELETE,
			Path:      "/api/admin/settings/categories/:id",
			Handler:   settingHandler.DeleteCategory,
			Operation: "admin:setting:categories:delete",
			Tags:      "Admin - Settings",
			Summary:   "删除配置分类",
		},
	}...)

	// ==================== 系统配置 ====================
	allRoutes = append(allRoutes, []routes.Route{
		{
			Method:    routes.POST,
			Path:      "/api/admin/settings/batch",
			Handler:   settingHandler.BatchUpdateSettings,
			Operation: "admin:settings:batch:update",
			Tags:      "Admin - Settings",
			Summary:   "批量更新配置",
		},
		{
			Method:    routes.POST,
			Path:      "/api/admin/settings",
			Handler:   settingHandler.CreateSetting,
			Operation: "admin:settings:create",
			Tags:      "Admin - Settings",
			Summary:   "创建配置",
		},
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings",
			Handler:   settingHandler.GetSettings,
			Operation: "admin:settings:list",
			Tags:      "Admin - Settings",
			Summary:   "配置列表",
		},
		{
			Method:    routes.GET,
			Path:      "/api/admin/settings/:key",
			Handler:   settingHandler.GetSetting,
			Operation: "admin:settings:get",
			Tags:      "Admin - Settings",
			Summary:   "配置详情",
		},
		{
			Method:    routes.PUT,
			Path:      "/api/admin/settings/:key",
			Handler:   settingHandler.UpdateSetting,
			Operation: "admin:settings:update",
			Tags:      "Admin - Settings",
			Summary:   "更新配置",
		},
		{
			Method:    routes.DELETE,
			Path:      "/api/admin/settings/:key",
			Handler:   settingHandler.DeleteSetting,
			Operation: "admin:settings:delete",
			Tags:      "Admin - Settings",
			Summary:   "删除配置",
		},
	}...)

	return allRoutes
}
