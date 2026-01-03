package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// GetCategoryHandler 获取单个配置分类查询处理器
type GetCategoryHandler struct {
	queryRepo setting.SettingCategoryQueryRepository
}

// NewGetCategoryHandler 创建 GetCategoryHandler 实例
func NewGetCategoryHandler(queryRepo setting.SettingCategoryQueryRepository) *GetCategoryHandler {
	return &GetCategoryHandler{queryRepo: queryRepo}
}

// Handle 处理获取配置分类查询
func (h *GetCategoryHandler) Handle(ctx context.Context, query GetCategoryQuery) (*CategoryDTO, error) {
	category, err := h.queryRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return ToCategoryDTO(category), nil
}
