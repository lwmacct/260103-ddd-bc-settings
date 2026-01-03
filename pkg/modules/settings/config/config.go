package config

// Redis 缓存配置
type Redis struct {
	KeyPrefix string `koanf:"key-prefix"`
}

// Config Settings 模块配置
type Config struct {
	Redis Redis `koanf:"redis"`
}
