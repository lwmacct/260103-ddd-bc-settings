package routes

import (
	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-shared/pkg/platform/http/gin/routes"
)

// RouteMetadata 定义路由的元数据（用于审计、文档生成等）。
type RouteMetadata struct {
	Method    string // HTTP 方法：GET, POST, PUT, DELETE
	Path      string // URL 路径
	Operation string // 操作标识（URN 格式：scope:resource:action）
	Category  string // 审计分类
	Label     string // 审计标签
}

// AllRouteBindings 返回 Settings 模块的所有路由绑定。
// 这是规范的路由注册入口点，提供完整的路由元数据。
func AllRouteBindings(settingHandler *handler.SettingHandler) []routes.Route {
	return Admin(settingHandler)
}

// GetRouteMetadata 返回所有路由的元数据（用于审计日志、权限检查等）。
func GetRouteMetadata() []RouteMetadata {
	return []RouteMetadata{
		// Category Routes
		{Method: "GET", Path: "/api/admin/settings/categories", Operation: "admin:setting:categories:list", Category: "setting", Label: "配置分类列表"},
		{Method: "GET", Path: "/api/admin/settings/categories/:id", Operation: "admin:setting:categories:get", Category: "setting", Label: "配置分类详情"},
		{Method: "POST", Path: "/api/admin/settings/categories", Operation: "admin:setting:categories:create", Category: "setting", Label: "创建配置分类"},
		{Method: "PUT", Path: "/api/admin/settings/categories/:id", Operation: "admin:setting:categories:update", Category: "setting", Label: "更新配置分类"},
		{Method: "DELETE", Path: "/api/admin/settings/categories/:id", Operation: "admin:setting:categories:delete", Category: "setting", Label: "删除配置分类"},
		// Setting Routes
		{Method: "GET", Path: "/api/admin/settings", Operation: "admin:settings:list", Category: "setting", Label: "配置列表"},
		{Method: "GET", Path: "/api/admin/settings/:key", Operation: "admin:settings:get", Category: "setting", Label: "配置详情"},
		{Method: "POST", Path: "/api/admin/settings", Operation: "admin:settings:create", Category: "setting", Label: "创建配置"},
		{Method: "PUT", Path: "/api/admin/settings/:key", Operation: "admin:settings:update", Category: "setting", Label: "更新配置"},
		{Method: "DELETE", Path: "/api/admin/settings/:key", Operation: "admin:settings:delete", Category: "setting", Label: "删除配置"},
		{Method: "POST", Path: "/api/admin/settings/batch", Operation: "admin:settings:batch:update", Category: "setting", Label: "批量更新配置"},
	}
}
