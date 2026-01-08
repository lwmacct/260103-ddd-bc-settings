package persistence

import (
	"context"
	"log/slog"
	"time"

	settingdomain "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// settingCommandWithCacheInvalidation 写操作后失效下游缓存的装饰器。
//
// 简化设计：
//   - 不再缓存 Setting 实体本身（由 Application 层 Settings 缓存覆盖）
//   - 只负责写操作后失效下游缓存（Settings）
//
// 失效策略：
//   - Create/Update/Delete/BatchUpsert 后异步失效 Settings 缓存
type settingCommandWithCacheInvalidation struct {
	delegate       settingdomain.CommandRepository
	changeNotifier settingdomain.SettingChangeNotifier
}

// NewCachedSettingCommandRepository 创建带缓存失效的 Setting 命令仓储。
func NewCachedSettingCommandRepository(
	delegate settingdomain.CommandRepository,
	changeNotifier settingdomain.SettingChangeNotifier,
) settingdomain.CommandRepository {
	return &settingCommandWithCacheInvalidation{
		delegate:       delegate,
		changeNotifier: changeNotifier,
	}
}

// Create 创建配置定义。
func (r *settingCommandWithCacheInvalidation) Create(ctx context.Context, s *settingdomain.Setting) error {
	if err := r.delegate.Create(ctx, s); err != nil {
		return err
	}
	r.invalidateSettingsCacheAsync(s.GetKeyCategory()) //nolint:contextcheck // 故意使用独立 context 进行异步失效
	return nil
}

// Update 更新配置定义。
func (r *settingCommandWithCacheInvalidation) Update(ctx context.Context, s *settingdomain.Setting) error {
	if err := r.delegate.Update(ctx, s); err != nil {
		return err
	}
	r.invalidateSettingsCacheAsync(s.GetKeyCategory()) //nolint:contextcheck // 故意使用独立 context 进行异步失效
	return nil
}

// Delete 删除配置定义。
func (r *settingCommandWithCacheInvalidation) Delete(ctx context.Context, key string) error {
	if err := r.delegate.Delete(ctx, key); err != nil {
		return err
	}
	r.invalidateSettingsCacheAsync("") //nolint:contextcheck // 删除操作失效所有缓存
	return nil
}

// BatchUpsert 批量插入或更新配置定义。
func (r *settingCommandWithCacheInvalidation) BatchUpsert(ctx context.Context, settings []*settingdomain.Setting) error {
	if err := r.delegate.BatchUpsert(ctx, settings); err != nil {
		return err
	}
	r.invalidateSettingsCacheAsync("") //nolint:contextcheck // 批量操作失效所有缓存
	return nil
}

// invalidateSettingsCacheAsync 异步失效 Settings 缓存。
func (r *settingCommandWithCacheInvalidation) invalidateSettingsCacheAsync(categoryKey string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		if categoryKey == "" {
			err = r.changeNotifier.NotifyAllChanged(ctx)
		} else {
			err = r.changeNotifier.NotifyCategoryChanged(ctx, categoryKey)
		}

		if err != nil {
			slog.Warn("failed to notify setting changes", "category", categoryKey, "error", err.Error())
		}
	}()
}

var _ settingdomain.CommandRepository = (*settingCommandWithCacheInvalidation)(nil)
