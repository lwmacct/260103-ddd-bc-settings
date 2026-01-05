package config

// Redis 缓存配置
type Redis struct {
	KeyPrefix string `koanf:"key-prefix"`
}

// HTTP 路由配置
type HTTP struct {
	BasePath string `koanf:"base-path"`
}

// Config Settings 模块配置
type Config struct {
	Redis Redis `koanf:"redis"`
	HTTP  HTTP  `koanf:"http"`
}
