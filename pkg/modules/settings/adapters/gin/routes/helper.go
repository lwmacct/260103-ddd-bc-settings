package routes

import "strings"

// buildPath 构建完整路径。
// basePath: 基础路径（如 /api/admin/settings）
// subPath: 子路径（如 /categories/:id）
func buildPath(basePath, subPath string) string {
	basePath = "/" + strings.Trim(basePath, "/")
	if subPath == "" {
		return basePath
	}
	return basePath + subPath
}
