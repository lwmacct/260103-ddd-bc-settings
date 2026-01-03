package setting

import "context"

// ValidationContext 验证上下文，包含验证所需的数据。
type ValidationContext struct {
	// Key 当前验证的设置键
	Key string
	// Value 当前验证的值
	Value any
	// Rule 验证规则（JSON Logic 格式）
	Rule string
	// Message 自定义错误消息
	Message string
	// AllSettings 所有设置的当前值（用于跨字段验证）
	AllSettings map[string]any
}

// ValidationResult 验证结果。
type ValidationResult struct {
	// Valid 是否验证通过
	Valid bool
	// Message 错误消息（验证失败时）
	Message string
}

// Validator 设置值验证器接口。
//
// 验证器负责在设置值保存前执行验证逻辑。
// 实现应支持 JSON Logic 规则格式，也可向后兼容简单规则格式。
type Validator interface {
	// Validate 验证单个设置值。
	//
	// 参数：
	//   - ctx: 上下文
	//   - vctx: 验证上下文，包含设置键、值、规则等
	//
	// 返回：
	//   - ValidationResult: 验证结果
	//   - error: 验证过程中的系统错误（非业务验证错误）
	Validate(ctx context.Context, vctx *ValidationContext) (*ValidationResult, error)

	// ValidateBatch 批量验证多个设置值。
	//
	// 返回 map[key]message，只包含验证失败的设置。
	ValidateBatch(ctx context.Context, items []*ValidationContext) (map[string]string, error)
}
