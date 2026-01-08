package setting

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// uiConfigRaw 内部结构用于解析 UIConfig JSONB
type uiConfigRaw struct {
	Hint      string              `json:"hint"`
	Options   []SelectOptionDTO   `json:"options"`
	DependsOn *DependsOnConfigDTO `json:"depends_on"`
}

// parseUIConfig 解析 UIConfig JSON 字符串。
//
// 如果解析失败，返回空的 UIConfigDTO（不返回错误，由调用方决定如何处理）。
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

// parseValidation 解析 Validation 字符串为 any 类型。
//
// Validation 字段存储 JSON Logic 规则（字符串格式的 JSON）。
// 如果解析失败，返回 nil（不返回错误，由调用方决定如何处理）。
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

// toCategoryMetaDTOs 将 SettingCategory 实体列表转换为 CategoryMetaDTO 列表。
//
// 用于 Settings API 返回分类元数据（不含 ID，仅 Key/Label/Icon/Order）。
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
