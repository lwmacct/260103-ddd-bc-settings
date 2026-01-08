package setting

import "errors"

var (
	// ErrDefinitionNotFound 配置定义不存在
	ErrDefinitionNotFound = errors.New("setting not found")

	// ErrDefinitionKeyExists 配置定义键已存在
	ErrDefinitionKeyExists = errors.New("setting key already exists")

	// ErrUserSettingNotFound 用户配置不存在（IAM 模块集成时启用）
	// TODO: 与 IAM 模块集成时取消注释
	// ErrUserSettingNotFound = errors.New("user setting not found")

	// ErrInvalidValueType 无效的值类型
	ErrInvalidValueType = errors.New("invalid value type")

	// ErrInvalidInputType 无效的控件类型
	ErrInvalidInputType = errors.New("invalid input type")

	// ErrInvalidValue 无效的配置值
	ErrInvalidValue = errors.New("invalid setting value")

	// ErrCategoryNotFound 配置分类不存在
	ErrCategoryNotFound = errors.New("category not found")

	// ErrValidationFailed 验证失败
	ErrValidationFailed = errors.New("validation failed")

	// ErrInvalidValidationRule 无效的验证规则
	ErrInvalidValidationRule = errors.New("invalid validation rule")

	// ErrInvalidScope 无效的配置作用域（废弃，保留向后兼容）
	ErrInvalidScope = errors.New("invalid scope")

	// ErrInvalidVisibleAt 无效的最小可见级别
	ErrInvalidVisibleAt = errors.New("invalid visible_at level")

	// ErrInvalidConfigurableAt 无效的最大可配置级别
	ErrInvalidConfigurableAt = errors.New("invalid configurable_at level")

	// ErrCannotOverrideSystemSetting 系统设置不能被用户覆盖
	ErrCannotOverrideSystemSetting = errors.New("cannot override system setting")

	// ErrInvalidKeyFormat 无效的配置键格式
	ErrInvalidKeyFormat = errors.New("invalid key format")

	// ErrInvalidCategoryID 无效的分类 ID
	ErrInvalidCategoryID = errors.New("invalid category id")

	// ErrInvalidCategory 无效的分类字符串（空值）
	ErrInvalidCategory = errors.New("invalid category")
)
