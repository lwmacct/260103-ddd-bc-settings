package setting

import "errors"

var (
	// ErrDefinitionNotFound 配置定义不存在
	ErrDefinitionNotFound = errors.New("配置定义不存在")

	// ErrDefinitionKeyExists 配置定义键已存在
	ErrDefinitionKeyExists = errors.New("配置定义键已存在")

	// ErrUserSettingNotFound 用户配置不存在
	ErrUserSettingNotFound = errors.New("用户配置不存在")

	// ErrInvalidValueType 无效的值类型
	ErrInvalidValueType = errors.New("无效的值类型")

	// ErrInvalidInputType 无效的控件类型
	ErrInvalidInputType = errors.New("无效的控件类型")

	// ErrInvalidValue 无效的配置值
	ErrInvalidValue = errors.New("无效的配置值")

	// ErrCategoryNotFound 配置分类不存在
	ErrCategoryNotFound = errors.New("配置分类不存在")

	// ErrValidationFailed 验证失败
	ErrValidationFailed = errors.New("验证失败")

	// ErrInvalidValidationRule 无效的验证规则
	ErrInvalidValidationRule = errors.New("无效的验证规则")

	// ErrInvalidScope 无效的配置作用域
	ErrInvalidScope = errors.New("无效的配置作用域")

	// ErrCannotOverrideSystemSetting 系统设置不能被用户覆盖
	ErrCannotOverrideSystemSetting = errors.New("系统设置不能被用户覆盖")

	// ErrInvalidKeyFormat 无效的配置键格式
	ErrInvalidKeyFormat = errors.New("无效的配置键格式")

	// ErrInvalidCategoryID 无效的分类 ID
	ErrInvalidCategoryID = errors.New("无效的分类 ID")
)
