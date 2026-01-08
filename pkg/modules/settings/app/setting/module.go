package setting

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/infra/persistence"
)

// SettingUseCases 聚合 Setting 模块的所有 UseCase Handler。
//
// 按功能领域分组：
//   - Setting: 配置定义的 CRUD 操作
//   - Category: 配置分类的 CRUD 操作
type SettingUseCases struct {
	// Setting Command Handlers
	Create      *CreateHandler
	Update      *UpdateHandler
	Delete      *DeleteHandler
	BatchUpdate *BatchUpdateHandler

	// Setting Query Handlers
	Get            *GetHandler
	List           *ListHandler
	ListSettings   *ListSettingsHandler
	PublicSettings *PublicSettingsHandler

	// Category Command Handlers
	CreateCategory *CreateCategoryHandler
	UpdateCategory *UpdateCategoryHandler
	DeleteCategory *DeleteCategoryHandler

	// Category Query Handlers
	GetCategory    *GetCategoryHandler
	ListCategories *ListCategoriesHandler
}

// UseCaseModule 提供 Settings 模块的 UseCase 层。
//
// 注册所有 Command 和 Query Handler。
var UseCaseModule = fx.Module("settings.usecase",
	fx.Provide(
		newSettingUseCases,
	),
)

// newSettingUseCases 创建 Setting UseCase 聚合。
func newSettingUseCases(repos persistence.SettingRepositories, settingsCache SettingsCacheService) *SettingUseCases {
	// Command Handlers（缓存失效由 Repository 装饰器自动处理）
	createHandler := NewCreateHandler(repos.Command, repos.Query)
	updateHandler := NewUpdateHandler(repos.Command, repos.Query)
	deleteHandler := NewDeleteHandler(repos.Command, repos.Query)
	batchUpdateHandler := NewBatchUpdateHandler(repos.Command, repos.Query)

	// Query Handlers（需要缓存服务进行查询）
	getHandler := NewGetHandler(repos.Query)
	listHandler := NewListHandler(repos.Query)

	// Category Handlers（暂时保留缓存失效逻辑，待后续优化）
	createCategoryHandler := NewCreateCategoryHandler(repos.CategoryCommand, repos.CategoryQuery, settingsCache)
	updateCategoryHandler := NewUpdateCategoryHandler(repos.CategoryCommand, repos.CategoryQuery, settingsCache)
	deleteCategoryHandler := NewDeleteCategoryHandler(repos.CategoryCommand, repos.CategoryQuery, repos.Query, settingsCache)
	getCategoryHandler := NewGetCategoryHandler(repos.CategoryQuery)
	listCategoriesHandler := NewListCategoriesHandler(repos.CategoryQuery)

	return &SettingUseCases{
		Create:         createHandler,
		Update:         updateHandler,
		Delete:         deleteHandler,
		BatchUpdate:    batchUpdateHandler,
		Get:            getHandler,
		List:           listHandler,
		ListSettings:   NewListSettingsHandler(repos.Query, repos.CategoryQuery, settingsCache),
		PublicSettings: NewPublicSettingsHandler(repos.Query, repos.CategoryQuery, settingsCache),
		CreateCategory: createCategoryHandler,
		UpdateCategory: updateCategoryHandler,
		DeleteCategory: deleteCategoryHandler,
		GetCategory:    getCategoryHandler,
		ListCategories: listCategoriesHandler,
	}
}
