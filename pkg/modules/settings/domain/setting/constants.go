package setting

// ============================================================================
// Scope 常量
// ============================================================================

const (
	// ScopeSystem 系统设置，全局唯一，管理员直接修改 DefaultValue
	ScopeSystem = "system"
	// ScopeOrg 组织设置，Org 可配置，Team 继承但不可覆盖
	ScopeOrg = "org"
	// ScopeTeam 团队设置，Team 可配置，可继承 Org 设置
	ScopeTeam = "team"
	// ScopeUser 用户设置，DefaultValue 作为初始值，用户可覆盖
	ScopeUser = "user"
)

// ============================================================================
// 配置分类常量
// ============================================================================

const (
	// CategoryGeneral 通用配置（站点名称、Logo 等）
	CategoryGeneral = "general"
	// CategorySecurity 安全配置（密码策略、登录限制等）
	CategorySecurity = "security"
	// CategoryNotification 通知配置（邮件、短信等）
	CategoryNotification = "notification"
	// CategoryBackup 备份配置（备份周期、保留策略等）
	CategoryBackup = "backup"
)

// ============================================================================
// 值类型常量
// ============================================================================

const (
	// ValueTypeString 字符串类型，使用文本输入框
	ValueTypeString = "string"
	// ValueTypeNumber 数值类型，使用数字输入框
	ValueTypeNumber = "number"
	// ValueTypeBoolean 布尔类型，使用开关控件
	ValueTypeBoolean = "boolean"
	// ValueTypeJSON JSON 类型，使用 JSON 编辑器
	ValueTypeJSON = "json"
)

// ============================================================================
// Input Type 常量（从前端 input_types.go 迁移）
// ============================================================================

const (
	// InputTypeText 文本输入框
	InputTypeText = "text"
	// InputTypeEmail 邮箱输入框
	InputTypeEmail = "email"
	// InputTypeURL URL 输入框
	InputTypeURL = "url"
	// InputTypePassword 密码输入框
	InputTypePassword = "password"
	// InputTypeNumber 数字输入框
	InputTypeNumber = "number"
	// InputTypeBoolean 布尔开关
	InputTypeBoolean = "boolean"
	// InputTypeSelect 下拉选择框
	InputTypeSelect = "select"
	// InputTypeTextarea 多行文本框
	InputTypeTextarea = "textarea"
	// InputTypeColor 颜色选择器
	InputTypeColor = "color"
	// InputTypeDate 日期选择器
	InputTypeDate = "date"
	// InputTypeTime 时间选择器
	InputTypeTime = "time"
	// InputTypeDateTime 日期时间选择器
	InputTypeDateTime = "datetime"
	// InputTypeJSONEditor JSON 编辑器
	InputTypeJSONEditor = "json"
)
