package setting

// CategoryMeta Category 元数据。
// 用于定义 Category 的显示属性，供内部使用。
// 数据来源：从数据库 setting_categories 表查询。
type CategoryMeta struct {
	Label string `json:"label"` // 显示名称
	Icon  string `json:"icon"`  // Tab 图标（mdi-xxx）
	Order int    `json:"order"` // 排序权重
}

// SelectOption 下拉/单选/多选选项。
type SelectOption struct {
	Value string `json:"value"`          // 选项值
	Label string `json:"label"`          // 显示文本
	Icon  string `json:"icon,omitempty"` // 可选图标（mdi-xxx）
}

// OptionsConfig 选项配置。
// 用于 select, radio, checkbox 等控件的选项列表。
type OptionsConfig struct {
	Items []SelectOption `json:"items"`
}

// DependsOnConfig 依赖关系配置。
// 用于控制设置项的可用状态，当依赖的设置项满足条件时才可编辑。
type DependsOnConfig struct {
	Key      string `json:"key"`                // 依赖的设置项 key
	Value    any    `json:"value,omitempty"`    // 期望的值
	Operator string `json:"operator,omitempty"` // 比较操作符: eq, ne, gt, lt（默认 eq）
}

// InputType 常量已移至 domain/setting/input_types.go
// 使用 setting.InputTypeText, setting.InputTypeEmail 等
