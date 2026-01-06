package setting

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ==================== Category Mappers ====================

// ToCategoryDTO 将 SettingCategory 实体转换为 CategoryDTO
func ToCategoryDTO(c *setting.SettingCategory) *CategoryDTO {
	if c == nil {
		return nil
	}

	return &CategoryDTO{
		ID:        c.ID,
		Key:       c.Key,
		Label:     c.Label,
		Icon:      c.Icon,
		Order:     c.Order,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToCategoryListDTO 将 SettingCategory 实体列表转换为 CategoryDTO 列表
func ToCategoryListDTO(categories []*setting.SettingCategory) []CategoryDTO {
	if len(categories) == 0 {
		return []CategoryDTO{}
	}

	dtos := make([]CategoryDTO, 0, len(categories))
	for _, c := range categories {
		if dto := ToCategoryDTO(c); dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	return dtos
}

// ==================== Setting Mappers ====================

// uiConfigRaw 内部结构用于解析 UIConfig JSONB
type uiConfigRaw struct {
	Hint      string              `json:"hint"`
	Options   []SelectOptionDTO   `json:"options"`
	DependsOn *DependsOnConfigDTO `json:"depends_on"`
}

// parseUIConfig 解析 UIConfig JSON 字符串
func parseUIConfig(jsonStr string) UIConfigDTO {
	if jsonStr == "" || jsonStr == "{}" {
		return UIConfigDTO{}
	}

	var raw uiConfigRaw
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return UIConfigDTO{}
	}

	return UIConfigDTO(raw)
}

// parseValidation 解析 Validation 字符串为 any 类型
func parseValidation(validation string) any {
	if validation == "" {
		return nil
	}
	var result any
	if err := json.Unmarshal([]byte(validation), &result); err != nil {
		return nil
	}
	return result
}

// ToSettingDTO 将 Setting 实体转换为 SettingDTO
func ToSettingDTO(s *setting.Setting) *SettingDTO {
	if s == nil {
		return nil
	}

	// 设置默认 InputType
	inputType := s.InputType
	if inputType == "" {
		inputType = "text"
	}

	return &SettingDTO{
		ID:             s.ID,
		Key:            s.Key,
		DefaultValue:   s.DefaultValue,
		VisibleAt:      s.VisibleAt,
		ConfigurableAt: s.ConfigurableAt,
		CategoryID:     s.CategoryID,
		Group:          s.Group,
		ValueType:      s.ValueType,
		Label:          s.Label,
		Order:          s.Order,
		InputType:      inputType,
		Validation:     parseValidation(s.Validation),
		UIConfig:       parseUIConfig(s.UIConfig),
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

// ToSettingListDTO 将 Setting 实体列表转换为 SettingDTO 列表
func ToSettingListDTO(settings []*setting.Setting) []SettingDTO {
	if len(settings) == 0 {
		return []SettingDTO{}
	}

	dtos := make([]SettingDTO, 0, len(settings))
	for _, s := range settings {
		if dto := ToSettingDTO(s); dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	return dtos
}

// ToSettingsItemDTO 将 Setting 转换为 SettingsItemDTO（Admin 场景，包含全部字段）
//
// categoryKey: 分类 Key（从 Category 实体获取）
func ToSettingsItemDTO(s *setting.Setting, categoryKey string) *SettingsItemDTO {
	if s == nil {
		return nil
	}

	// 设置默认 InputType
	inputType := s.InputType
	if inputType == "" {
		inputType = "text"
	}

	// 设置默认 Group
	group := s.Group
	if group == "" {
		group = "default"
	}

	return &SettingsItemDTO{
		Key:            s.Key,
		Category:       categoryKey,
		Group:          group,
		Value:          s.DefaultValue,
		DefaultValue:   s.DefaultValue,
		IsCustomized:   false,
		VisibleAt:      s.VisibleAt,
		ConfigurableAt: s.ConfigurableAt,
		ValueType:      s.ValueType,
		Label:          s.Label,
		Order:          s.Order,
		InputType:      inputType,
		Validation:     parseValidation(s.Validation),
		UIConfig:       parseUIConfig(s.UIConfig),
	}
}

// ==================== UserSetting Mappers ====================

// ToUserSettingDTO 将 Setting 定义和可选的 UserSetting 合并为 UserSettingDTO
func ToUserSettingDTO(s *setting.Setting, us *setting.UserSetting) *UserSettingDTO {
	if s == nil {
		return nil
	}

	// 设置默认 InputType
	inputType := s.InputType
	if inputType == "" {
		inputType = "text"
	}

	dto := &UserSettingDTO{
		Key:          s.Key,
		Value:        s.DefaultValue, // 默认使用系统默认值
		DefaultValue: s.DefaultValue,
		IsCustomized: false,
		CategoryID:   s.CategoryID,
		Group:        s.Group,
		ValueType:    s.ValueType,
		Label:        s.Label,
		Order:        s.Order,
		InputType:    inputType,
		Validation:   parseValidation(s.Validation),
		UIConfig:     parseUIConfig(s.UIConfig),
	}

	// 如果有用户自定义值，使用用户值
	if us != nil {
		dto.Value = us.Value
		dto.IsCustomized = true
	}

	return dto
}

// ToUserSettingsItemDTO 将 Setting 定义和可选的 UserSetting 合并为 SettingsItemDTO
//
// User 场景现在也返回 VisibleAt 和 ConfigurableAt 字段：
//   - VisibleAt: 前端根据此字段判断可见性
//   - ConfigurableAt: 前端根据此字段判断可编辑性
//
// categoryKey: 分类 Key（从 Category 实体获取）
func ToUserSettingsItemDTO(s *setting.Setting, us *setting.UserSetting, categoryKey string) *SettingsItemDTO {
	if s == nil {
		return nil
	}

	// 设置默认 InputType
	inputType := s.InputType
	if inputType == "" {
		inputType = "text"
	}

	// 设置默认 Group
	group := s.Group
	if group == "" {
		group = "default"
	}

	dto := &SettingsItemDTO{
		Key:            s.Key,
		Category:       categoryKey,
		Group:          group,
		Value:          s.DefaultValue, // 默认使用系统默认值
		DefaultValue:   s.DefaultValue,
		IsCustomized:   false,
		VisibleAt:      s.VisibleAt,      // 返回 VisibleAt，前端判断可见性
		ConfigurableAt: s.ConfigurableAt, // 返回 ConfigurableAt，前端判断可编辑性
		ValueType:      s.ValueType,
		Label:          s.Label,
		Order:          s.Order,
		InputType:      inputType,
		Validation:     parseValidation(s.Validation),
		UIConfig:       parseUIConfig(s.UIConfig),
	}

	// 如果有用户自定义值，使用用户值（仅 VisibleAt=user 时才有用户值）
	if us != nil {
		dto.Value = us.Value
		dto.IsCustomized = true
	}

	return dto
}

// toCategoryMetaDTOs 将 SettingCategory 实体列表转换为 CategoryMetaDTO 列表。
func toCategoryMetaDTOs(categories []*setting.SettingCategory) []CategoryMetaDTO {
	result := make([]CategoryMetaDTO, 0, len(categories))
	for _, cat := range categories {
		result = append(result, CategoryMetaDTO{
			Category: cat.Key,
			Label:    cat.Label,
			Icon:     cat.Icon,
			Order:    cat.Order,
		})
	}
	return result
}
