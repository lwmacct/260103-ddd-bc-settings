package setting

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ListSettingsHandler 获取设置 Settings 查询处理器
type ListSettingsHandler struct {
	settingQueryRepo  setting.QueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
	settingsCache     SettingsCacheService
}

// NewListSettingsHandler 创建 ListSettingsHandler 实例
func NewListSettingsHandler(
	settingQueryRepo setting.QueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
	settingsCache SettingsCacheService,
) *ListSettingsHandler {
	return &ListSettingsHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		settingsCache:     settingsCache,
	}
}

// Handle 处理获取设置 Settings 查询
// 返回按 Category → Group → Settings 层级组织的精简数据
//
// 支持 CategoryKey 过滤：
//   - 为空时返回全量系统设置（用于总配置页）
//   - 指定 Key 时只返回该分类（用于分散页面的懒加载）
//
// 缓存策略：
//   - 先查缓存，命中直接返回
//   - 未命中时查数据库，同步回写缓存
func (h *ListSettingsHandler) Handle(ctx context.Context, query ListSettingsQuery) ([]SettingsCategoryDTO, error) {
	// 1. 查缓存
	if cached, err := h.settingsCache.GetAdminSettings(ctx, query.CategoryKey); err == nil && cached != nil {
		return cached, nil
	}

	// 2. 缓存未命中，执行原有逻辑
	result, err := h.fetchAndBuild(ctx, query)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存（仅非空结果，避免缓存无效数据）
	if len(result) > 0 {
		if err := h.settingsCache.SetAdminSettings(ctx, query.CategoryKey, result); err != nil {
			slog.Warn("failed to cache admin settings", "categoryKey", query.CategoryKey, "error", err.Error())
		}
	}

	return result, nil
}

// fetchAndBuild 从数据库获取数据并构建 Settings
func (h *ListSettingsHandler) fetchAndBuild(ctx context.Context, query ListSettingsQuery) ([]SettingsCategoryDTO, error) {
	// 1. 根据 CategoryKey 决定查询范围
	settings, err := h.fetchSettings(ctx, query.CategoryKey)
	if err != nil {
		return nil, err
	}

	// 2. 查询所有分类元数据
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	// 3. 使用共享构建器
	builder := NewSettingsBuilder(categories)
	return builder.Build(settings, nil, AdminSettingMapper), nil
}

// fetchSettings 根据 CategoryKey 获取设置列表
func (h *ListSettingsHandler) fetchSettings(ctx context.Context, categoryKey string) ([]*setting.Setting, error) {
	// 全量查询
	if categoryKey == "" {
		settings, err := h.settingQueryRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch settings: %w", err)
		}
		return settings, nil
	}

	// 按 Category Key 过滤
	category, err := h.categoryQueryRepo.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by key %q: %w", categoryKey, err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found: %s", categoryKey)
	}

	settings, err := h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}
	return settings, nil
}
