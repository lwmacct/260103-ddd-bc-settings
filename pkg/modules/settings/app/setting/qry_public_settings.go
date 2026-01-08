package setting

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// PublicSettingsHandler 公开设置查询处理器
type PublicSettingsHandler struct {
	settingQueryRepo  setting.QueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
	settingsCache     SettingsCacheService
}

// NewPublicSettingsHandler 创建 PublicSettingsHandler 实例
func NewPublicSettingsHandler(
	settingQueryRepo setting.QueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
	settingsCache SettingsCacheService,
) *PublicSettingsHandler {
	return &PublicSettingsHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		settingsCache:     settingsCache,
	}
}

// Handle 处理公开设置查询
// 返回 VisibleAt="public" 的设置，扁平结构
//
// 缓存策略：
//   - 先查缓存，命中直接返回
//   - 未命中时查数据库，同步回写缓存
func (h *PublicSettingsHandler) Handle(ctx context.Context, query PublicSettingsQuery) ([]PublicSettingItemDTO, error) {
	// 1. 查缓存
	if cached, err := h.settingsCache.GetPublicSettings(ctx, query.CategoryKey); err == nil && cached != nil {
		return cached, nil
	}

	// 2. 缓存未命中，查询数据库
	result, err := h.fetchAndBuild(ctx, query)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存（仅非空结果）
	if len(result) > 0 {
		if err := h.settingsCache.SetPublicSettings(ctx, query.CategoryKey, result); err != nil {
			slog.Warn("failed to cache public settings", "categoryKey", query.CategoryKey, "error", err.Error())
		}
	}

	return result, nil
}

// fetchAndBuild 从数据库获取数据并构建公开设置
func (h *PublicSettingsHandler) fetchAndBuild(ctx context.Context, query PublicSettingsQuery) ([]PublicSettingItemDTO, error) {
	// 1. 查询公开设置
	settings, err := h.settingQueryRepo.FindPublic(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public settings: %w", err)
	}

	// 2. 如果指定了 CategoryKey，过滤设置
	if query.CategoryKey != "" {
		category, catErr := h.categoryQueryRepo.FindByKey(ctx, query.CategoryKey)
		if catErr != nil {
			return nil, fmt.Errorf("failed to find category by key %q: %w", query.CategoryKey, catErr)
		}
		if category == nil {
			return nil, fmt.Errorf("category not found: %s", query.CategoryKey)
		}
		filtered := make([]*setting.Setting, 0)
		for _, s := range settings {
			if s.CategoryID == category.ID {
				filtered = append(filtered, s)
			}
		}
		settings = filtered
	}

	if len(settings) == 0 {
		return []PublicSettingItemDTO{}, nil
	}

	// 3. 转换为扁平 DTO 列表
	result := make([]PublicSettingItemDTO, 0, len(settings))
	for _, s := range settings {
		result = append(result, PublicSettingItemDTO{
			Key:   s.Key,
			Value: s.DefaultValue,
			Label: s.Label,
		})
	}

	// 4. 按 Key 排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})

	return result, nil
}
