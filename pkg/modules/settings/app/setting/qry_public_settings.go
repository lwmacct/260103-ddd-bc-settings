package setting

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
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
// 返回 VisibleAt="public" 的设置，按 Category → Group → Settings 层级组织
//
// 缓存策略：
//   - 先查缓存，命中直接返回
//   - 未命中时查数据库，同步回写缓存
func (h *PublicSettingsHandler) Handle(ctx context.Context, query PublicSettingsQuery) ([]PublicSettingsCategoryDTO, error) {
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
func (h *PublicSettingsHandler) fetchAndBuild(ctx context.Context, query PublicSettingsQuery) ([]PublicSettingsCategoryDTO, error) {
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

	// 3. 查询所有分类元数据
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	// 4. 构建公开设置响应
	return h.buildPublicSettings(settings, categories), nil
}

// buildPublicSettings 构建公开设置响应（精简版）
func (h *PublicSettingsHandler) buildPublicSettings(
	settings []*setting.Setting,
	categories []*setting.SettingCategory,
) []PublicSettingsCategoryDTO {
	// 构建分类映射
	categoryByID := make(map[uint]*setting.SettingCategory)
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
	}

	// 按分类和分组聚合
	type groupKey struct {
		categoryID uint
		group      string
	}
	groupMap := make(map[groupKey][]PublicSettingItemDTO)
	categoryIDs := make(map[uint]bool)

	for _, s := range settings {
		key := groupKey{categoryID: s.CategoryID, group: s.Group}
		groupMap[key] = append(groupMap[key], PublicSettingItemDTO{
			Key:   s.Key,
			Value: s.DefaultValue,
			Label: s.Label,
		})
		categoryIDs[s.CategoryID] = true
	}

	// 构建结果
	var result []PublicSettingsCategoryDTO
	for catID := range categoryIDs {
		cat, ok := categoryByID[catID]
		if !ok {
			continue
		}

		// 收集该分类下的所有分组
		groupNames := make(map[string]bool)
		for key := range groupMap {
			if key.categoryID == catID {
				groupNames[key.group] = true
			}
		}

		// 构建分组列表
		var groups []PublicSettingsGroupDTO
		for groupName := range groupNames {
			key := groupKey{categoryID: catID, group: groupName}
			items := groupMap[key]
			// 按 key 排序
			sort.Slice(items, func(i, j int) bool {
				return items[i].Key < items[j].Key
			})
			groups = append(groups, PublicSettingsGroupDTO{
				Name:     groupName,
				Settings: items,
			})
		}

		// 按分组名排序
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Name < groups[j].Name
		})

		result = append(result, PublicSettingsCategoryDTO{
			Category: cat.Key,
			Label:    cat.Label,
			Groups:   groups,
		})
	}

	// 按分类 key 排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Category < result[j].Category
	})

	return result
}
