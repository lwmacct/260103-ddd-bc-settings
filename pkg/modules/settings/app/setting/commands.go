package setting

// ==================== Category Commands ====================

// CreateCategoryCommand 创建配置分类命令
type CreateCategoryCommand struct {
	Key   string
	Label string
	Icon  string
	Order int
}

// UpdateCategoryCommand 更新配置分类命令
type UpdateCategoryCommand struct {
	ID    uint
	Label string
	Icon  string
	Order int
}

// DeleteCategoryCommand 删除配置分类命令
type DeleteCategoryCommand struct {
	ID uint
}

// ==================== Setting Commands ====================

// CreateCommand 创建配置命令
type CreateCommand struct {
	Key          string
	DefaultValue any // JSONB 原生值
	CategoryID   uint
	Group        string
	ValueType    string
	Label        string
	UIConfig     string
	Order        int
}

// UpdateCommand 更新配置命令
type UpdateCommand struct {
	Key          string
	DefaultValue any // JSONB 原生值
	Label        string
	UIConfig     string
	Order        int
}

// DeleteCommand 删除配置命令
type DeleteCommand struct {
	Key string
}

// SettingItemCommand 配置项（用于批量操作）
type SettingItemCommand struct {
	Key   string
	Value any
}

// BatchUpdateCommand 批量更新系统配置命令
type BatchUpdateCommand struct {
	Settings []SettingItemCommand
}

// ==================== UserSetting Commands ====================

// UserSetCommand 设置用户配置命令
type UserSetCommand struct {
	UserID uint
	Key    string
	Value  any // JSONB 原生值
}

// UserBatchSetCommand 批量设置用户配置命令
type UserBatchSetCommand struct {
	UserID   uint
	Settings []SettingItemCommand
}

// UserResetCommand 重置用户配置命令（恢复默认值）
type UserResetCommand struct {
	UserID uint
	Key    string
}

// UserResetAllCommand 重置用户所有配置命令
type UserResetAllCommand struct {
	UserID uint
}
