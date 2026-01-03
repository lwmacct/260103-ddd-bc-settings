package setting

import "context"

// CacheInvalidator 缓存失效器接口（Domain 层定义）。
//
// 用于 Infrastructure 层在写操作后通知缓存失效。
// 这是一个简化接口，只包含失效方法，不涉及具体的缓存实现。
type CacheInvalidator interface {
	// DeleteAll 删除所有缓存
	DeleteAll(ctx context.Context) error

	// DeleteByCategoryKey 删除指定 category 的所有缓存
	DeleteByCategoryKey(ctx context.Context, categoryKey string) error
}
