package setting

import "context"

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 配置定义写操作接口。
type CommandRepository interface {
	// Create 创建配置定义
	Create(ctx context.Context, setting *Setting) error

	// Update 更新配置定义
	Update(ctx context.Context, setting *Setting) error

	// Delete 删除配置定义
	Delete(ctx context.Context, key string) error

	// BatchUpsert 批量插入或更新配置定义
	BatchUpsert(ctx context.Context, settings []*Setting) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 配置定义读操作接口。
type QueryRepository interface {
	// FindByKey 根据 Key 查找配置定义
	FindByKey(ctx context.Context, key string) (*Setting, error)

	// FindByKeys 根据多个 Key 批量查找配置定义
	FindByKeys(ctx context.Context, keys []string) ([]*Setting, error)

	// FindByCategoryID 根据分类 ID 查找配置定义列表
	FindByCategoryID(ctx context.Context, categoryID uint) ([]*Setting, error)

	// FindByScope 根据作用域查找配置定义列表
	//
	// Deprecated: 使用 FindByVisibleAt 代替
	FindByScope(ctx context.Context, scope string) ([]*Setting, error)

	// FindVisibleToUser 查找普通用户可见的配置定义
	// 可见条件：VisibleAt <= user
	FindVisibleToUser(ctx context.Context) ([]*Setting, error)

	// FindAll 查找所有配置定义
	FindAll(ctx context.Context) ([]*Setting, error)

	// ExistsByKey 检查 Key 是否已存在
	ExistsByKey(ctx context.Context, key string) (bool, error)

	// FindByVisibleAt 查询对指定级别可见的设置
	// 返回条件：查询级别的层级 >= visible_at 的层级
	FindByVisibleAt(ctx context.Context, visibleAt ScopeLevel) ([]*Setting, error)

	// FindByConfigurableAt 查询指定级别可配置的设置
	// 返回条件：查询级别的层级 >= configurable_at 的层级
	FindByConfigurableAt(ctx context.Context, configurableAt ScopeLevel) ([]*Setting, error)

	// FindByVisibleAndConfigurable 查询同时满足可见性和可配置性的设置
	FindByVisibleAndConfigurable(ctx context.Context, visibleAt, configurableAt ScopeLevel) ([]*Setting, error)
}

// ============================================================================
// Category Command Repository
// ============================================================================

// SettingCategoryCommandRepository 配置分类写操作接口。
type SettingCategoryCommandRepository interface {
	// Create 创建配置分类
	Create(ctx context.Context, category *SettingCategory) error

	// Update 更新配置分类
	Update(ctx context.Context, category *SettingCategory) error

	// Delete 删除配置分类
	Delete(ctx context.Context, id uint) error
}

// ============================================================================
// Category Query Repository
// ============================================================================

// SettingCategoryQueryRepository 配置分类读操作接口。
type SettingCategoryQueryRepository interface {
	// FindByID 根据 ID 查找配置分类
	FindByID(ctx context.Context, id uint) (*SettingCategory, error)

	// FindByKey 根据 Key 查找配置分类
	FindByKey(ctx context.Context, key string) (*SettingCategory, error)

	// FindAll 查找所有配置分类
	FindAll(ctx context.Context) ([]*SettingCategory, error)

	// ExistsByKey 检查 Key 是否已存在
	ExistsByKey(ctx context.Context, key string) (bool, error)
}
