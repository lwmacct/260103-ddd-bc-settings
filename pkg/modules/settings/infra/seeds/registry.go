package seeds

// DefaultSeeders 返回 Settings 模块的默认 Seeder 列表。
//
// 执行顺序：
//  1. SettingCategorySeeder - 配置分类（依赖优先）
//  2. SettingSeeder - 配置定义
func DefaultSeeders() []Seeder {
	return []Seeder{
		&SettingCategorySeeder{},
		&SettingSeeder{},
	}
}
