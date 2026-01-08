package config

// Config Settings 模块配置
type Config struct {
	// RedisKeyPrefix Redis 缓存键前缀
	// 格式：{basePrefix}{module}:
	// 例如："app:settings:" 或 "prod:settings:"
	RedisKeyPrefix string `koanf:"redis-key-prefix"`

	// AdminPath Settings 管理员 API 的路径（需要认证）
	AdminPath string `koanf:"admin-path"`

	// PublicPath Settings 公开 API 的路径（无需认证）
	PublicPath string `koanf:"public-path"`
}

// DefaultConfig 返回 Settings 模块的默认配置。
//
// 注意：此函数主要用于单元测试或手动测试。
// 生产环境的默认值应该在 internal/config.DefaultConfig() 中定义，
// 然后通过容器层的 newSettingsConfig() 转换为本模块的 Config。
func DefaultConfig() Config {
	return Config{
		RedisKeyPrefix: "app:settings:",
		AdminPath:      "/api/admin/settings",
		PublicPath:     "/api/public/settings",
	}
}
