package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// ListCategoriesHandler 获取配置分类列表查询处理器
type ListCategoriesHandler struct {
	queryRepo setting.SettingCategoryQueryRepository
}

// NewListCategoriesHandler 创建 ListCategoriesHandler 实例
func NewListCategoriesHandler(queryRepo setting.SettingCategoryQueryRepository) *ListCategoriesHandler {
	return &ListCategoriesHandler{queryRepo: queryRepo}
}

// Handle 处理获取配置分类列表查询
func (h *ListCategoriesHandler) Handle(ctx context.Context, _ ListCategoriesQuery) ([]CategoryDTO, error) {
	categories, err := h.queryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return ToCategoryListDTO(categories), nil
}
