package setting

import "context"

// SettingChangeNotifier 配置变更通知器接口。
//
// Domain 业务抽象：配置写操作后需要通知外部（如缓存失效）。
// Infrastructure 层实现此接口（如缓存服务）。
//
// 业务语义方法（不含技术词）：
//   - NotifyAllChanged：通知所有配置已变更
//   - NotifyCategoryChanged：通知指定分类的配置已变更
type SettingChangeNotifier interface {
	// NotifyAllChanged 通知所有配置已变更（对应全量失效）
	NotifyAllChanged(ctx context.Context) error

	// NotifyCategoryChanged 通知指定分类的配置已变更（对应分类失效）
	NotifyCategoryChanged(ctx context.Context, categoryKey string) error
}
