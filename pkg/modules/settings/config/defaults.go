package config

// DefaultConfig 返回 Settings 模块的默认配置。
func DefaultConfig() Config {
	return Config{
		Redis: Redis{
			KeyPrefix: "app:",
		},
		HTTP: HTTP{
			BasePath: "/api/admin/settings", // 默认基础路径
		},
	}
}
