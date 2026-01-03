package setting

import (
	"sort"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// SettingsBuilder 构建 Category → Group → Settings 层级结构
// 用于 ListSettingsHandler 和 UserListSettingsHandler 共享构建逻辑
type SettingsBuilder struct {
	categoryByID  map[uint]*setting.SettingCategory
	categoryByKey map[string]*setting.SettingCategory // 按 Key 索引，用于 O(1) 排序
}

// NewSettingsBuilder 创建 Settings 构建器
func NewSettingsBuilder(categories []*setting.SettingCategory) *SettingsBuilder {
	categoryByID := make(map[uint]*setting.SettingCategory, len(categories))
	categoryByKey := make(map[string]*setting.SettingCategory, len(categories))
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
		categoryByKey[cat.Key] = cat
	}
	return &SettingsBuilder{
		categoryByID:  categoryByID,
		categoryByKey: categoryByKey,
	}
}

// SettingMapper 将 Setting 转换为 SettingsItemDTO 的函数类型
// Admin 场景使用 ToSettingsItemDTO，User 场景使用 ToUserSettingsItemDTO
type SettingMapper func(s *setting.Setting, us *setting.UserSetting) *SettingsItemDTO

// Build 构建 Settings 层级结构
// settings: 配置定义列表
// userSettingMap: 用户配置映射（Admin 场景传 nil）
// mapper: 转换函数
func (b *SettingsBuilder) Build(
	settings []*setting.Setting,
	userSettingMap map[string]*setting.UserSetting,
	mapper SettingMapper,
) []SettingsCategoryDTO {
	// 按 CategoryID 分组
	categoryMap := make(map[uint]map[string][]SettingsItemDTO)

	for _, s := range settings {
		categoryID := s.CategoryID
		group := s.Group
		if group == "" {
			group = "default"
		}

		if _, ok := categoryMap[categoryID]; !ok {
			categoryMap[categoryID] = make(map[string][]SettingsItemDTO)
		}

		var us *setting.UserSetting
		if userSettingMap != nil {
			us = userSettingMap[s.Key]
		}
		dto := mapper(s, us)
		if dto != nil {
			categoryMap[categoryID][group] = append(categoryMap[categoryID][group], *dto)
		}
	}

	// 构建响应
	result := make([]SettingsCategoryDTO, 0, len(categoryMap))
	for categoryID, groupMap := range categoryMap {
		cat, ok := b.categoryByID[categoryID]
		if !ok {
			// 跳过未知 category
			continue
		}

		groups := make([]SettingsGroupDTO, 0, len(groupMap))
		for group, settingDTOs := range groupMap {
			// 按 Order 排序设置项
			sort.Slice(settingDTOs, func(i, j int) bool {
				return settingDTOs[i].Order < settingDTOs[j].Order
			})

			// Group 字段直接存储 label（如 "基本设置"）
			groups = append(groups, SettingsGroupDTO{
				Name:     group,
				Settings: settingDTOs,
			})
		}

		// 按分组内最小 Order 排序（Settings 已按 Order 排序，首个即最小）
		sort.Slice(groups, func(i, j int) bool {
			// default 组（无分组）排在最后
			if groups[i].Name == "default" {
				return false
			}
			if groups[j].Name == "default" {
				return true
			}
			// 按组内最小 Order 排序
			minOrderI := groups[i].Settings[0].Order
			minOrderJ := groups[j].Settings[0].Order
			return minOrderI < minOrderJ
		})

		result = append(result, SettingsCategoryDTO{
			Category: cat.Key,
			Label:    cat.Label,
			Icon:     cat.Icon,
			Groups:   groups,
		})
	}

	// 按 Category Order 排序（使用 categoryByKey 实现 O(1) 查找）
	sort.Slice(result, func(i, j int) bool {
		catI := b.categoryByKey[result[i].Category]
		catJ := b.categoryByKey[result[j].Category]
		if catI == nil || catJ == nil {
			return result[i].Category < result[j].Category
		}
		return catI.Order < catJ.Order
	})

	return result
}

// AdminSettingMapper Admin 场景的 Setting 转换器（包含全部字段）
func AdminSettingMapper(s *setting.Setting, _ *setting.UserSetting) *SettingsItemDTO {
	return ToSettingsItemDTO(s)
}

// UserSettingMapper User 场景的 Setting 转换器（省略权限字段，合并用户值）
func UserSettingMapper(s *setting.Setting, us *setting.UserSetting) *SettingsItemDTO {
	return ToUserSettingsItemDTO(s, us)
}
