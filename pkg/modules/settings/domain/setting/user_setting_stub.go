package setting

import "time"

// UserSettingStub 用户配置存根定义。
//
// TODO: user_settings 表由 IAM 模块管理，这是临时存根定义用于编译通过。
// 实际的 UserSetting 实体和仓储将在集成时由 IAM 模块提供。
//
// 当 Settings 模块与 IAM 模块集成后：
// 1. IAM 模块提供 UserSetting 实体和仓储
// 2. Settings 模块的 SettingsBuilder 接收 userValues map[string]*UserSetting
// 3. IAM 模块的 ListSettingsHandler 调用 SettingsBuilder 构建 schema

// UserSetting 用户配置（临时存根）。
//
// 存储用户对配置项的自定义值，覆盖系统默认值。
// Value 字段直接存储原生 JSON 值，类型应与对应配置项的 ValueType 一致。
type UserSetting struct {
	ID         uint      `json:"id"`          // 唯一标识
	UserID     uint      `json:"user_id"`     // 用户 ID，关联 users 表
	SettingKey string    `json:"setting_key"` // 配置键
	Value      any       `json:"value"`       // 用户自定义值（JSONB 原生值）
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

// ValidationContext 和 Validator 定义已移至 validator.go
