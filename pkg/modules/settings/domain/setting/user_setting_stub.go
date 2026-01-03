package setting

import "context"

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
type UserSetting struct {
	ID          uint
	UserID      uint
	SettingKey  string
	Value       any
	UpdatedAt   string
}

// ValidationContext 校验上下文（临时存根）。
type ValidationContext struct {
	// TODO: 由 IAM 模块提供完整的实现
}

// Validator 校验器接口（临时存根）。
type Validator interface {
	// TODO: 由 IAM 模块提供完整的实现
	Validate(ctx context.Context, vc ValidationContext, s *Setting, value any) error
}
