package setting

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
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
// 返回扁平结构的配置列表，每个 item 包含 category 和 group 字段供前端分组
//
// 支持 CategoryKey 过滤：
//   - 为空时返回全量系统设置（用于总配置页）
//   - 指定 Key 时只返回该分类（用于分散页面的懒加载）
//
// 缓存策略：
//   - 先查缓存，命中直接返回
//   - 未命中时查数据库，同步回写缓存
func (h *ListSettingsHandler) Handle(ctx context.Context, query ListSettingsQuery) ([]SettingsItemDTO, error) {
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
func (h *ListSettingsHandler) fetchAndBuild(ctx context.Context, query ListSettingsQuery) ([]SettingsItemDTO, error) {
	// 1. 根据 CategoryKey 决定查询范围
	settings, err := h.fetchSettings(ctx, query.CategoryKey)
	if err != nil {
		return nil, err
	}

	if len(settings) == 0 {
		return []SettingsItemDTO{}, nil
	}

	// 2. 查询所有分类元数据（用于填充 category key）
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	// 3. 构建 CategoryID -> Key 映射
	categoryKeyByID := make(map[uint]string, len(categories))
	for _, cat := range categories {
		categoryKeyByID[cat.ID] = cat.Key
	}

	// 4. 转换为扁平 DTO 列表
	result := make([]SettingsItemDTO, 0, len(settings))
	for _, s := range settings {
		categoryKey := categoryKeyByID[s.CategoryID]
		dto := ToSettingsItemDTO(s, categoryKey)
		if dto != nil {
			result = append(result, *dto)
		}
	}

	// 5. 按 Category Order + Group + Setting Order 排序
	categoryOrderByKey := make(map[string]int, len(categories))
	for _, cat := range categories {
		categoryOrderByKey[cat.Key] = cat.Order
	}

	sort.Slice(result, func(i, j int) bool {
		// 先按 Category Order
		catOrderI := categoryOrderByKey[result[i].Category]
		catOrderJ := categoryOrderByKey[result[j].Category]
		if catOrderI != catOrderJ {
			return catOrderI < catOrderJ
		}
		// 再按 Group 名称
		if result[i].Group != result[j].Group {
			// default 组排最后
			if result[i].Group == "default" {
				return false
			}
			if result[j].Group == "default" {
				return true
			}
			return result[i].Group < result[j].Group
		}
		// 最后按 Setting Order
		return result[i].Order < result[j].Order
	})

	return result, nil
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
