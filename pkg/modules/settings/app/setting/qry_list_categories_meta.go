package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// ListCategoriesMetaHandler 获取配置分类元信息列表处理器（管理员端）
// 返回全量分类元信息，用于懒加载场景
type ListCategoriesMetaHandler struct {
	categoryQueryRepo setting.SettingCategoryQueryRepository
}

// NewListCategoriesMetaHandler 创建 ListCategoriesMetaHandler 实例
func NewListCategoriesMetaHandler(
	categoryQueryRepo setting.SettingCategoryQueryRepository,
) *ListCategoriesMetaHandler {
	return &ListCategoriesMetaHandler{
		categoryQueryRepo: categoryQueryRepo,
	}
}

// Handle 处理获取分类元信息列表查询
// 返回全量分类元信息（不含 settings 数据）
func (h *ListCategoriesMetaHandler) Handle(ctx context.Context, _ ListCategoriesQuery) ([]CategoryMetaDTO, error) {
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all categories: %w", err)
	}

	return toCategoryMetaDTOs(categories), nil
}
