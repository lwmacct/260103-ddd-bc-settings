package setting

// ==================== Category Queries ====================

// GetCategoryQuery 获取单个配置分类查询
type GetCategoryQuery struct {
	ID uint
}

// ListCategoriesQuery 获取配置分类列表查询（管理员端，全量）
type ListCategoriesQuery struct{}

// UserListCategoriesQuery 获取用户可见的分类列表查询
// 只返回包含 scope="user" 设置的分类（用于懒加载场景）
type UserListCategoriesQuery struct {
	UserID uint
}

// ==================== Setting Queries ====================

// GetQuery 获取配置查询
type GetQuery struct {
	Key string
}

// ListQuery 获取配置列表查询
type ListQuery struct {
	CategoryID uint // 可选: 按类别 ID 过滤
}

// ListSettingsQuery 获取配置 Settings 查询（系统配置）
// Schema 返回按 Category → Group → Settings 层级组织的数据
//
// 支持按 CategoryKey 过滤，用于按需加载（懒加载）场景。
type ListSettingsQuery struct {
	CategoryKey string // 可选：按分类 Key 过滤（如 "general"），为空时返回全量
}

// ==================== UserSetting Queries ====================

// UserGetQuery 获取用户配置查询（合并默认值）
type UserGetQuery struct {
	UserID uint
	Key    string
}

// UserListQuery 获取用户配置列表查询（合并默认值）
type UserListQuery struct {
	UserID     uint
	CategoryID uint // 可选: 按类别 ID 过滤
}

// UserListSettingsQuery 获取用户配置 Settings 查询（带合并值）
//
// 支持按 CategoryKey 过滤，用于按需加载（懒加载）场景。
type UserListSettingsQuery struct {
	UserID      uint
	CategoryKey string // 可选：按分类 Key 过滤（如 "profile"），为空时返回全量
}
