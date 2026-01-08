package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// ListHandler 获取设置列表查询处理器
type ListHandler struct {
	settingQueryRepo setting.QueryRepository
}

// NewListHandler 创建 ListHandler 实例
func NewListHandler(settingQueryRepo setting.QueryRepository) *ListHandler {
	return &ListHandler{
		settingQueryRepo: settingQueryRepo,
	}
}

// Handle 处理获取设置列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]*SettingDTO, error) {
	var settings []*setting.Setting
	var err error

	if query.CategoryID != 0 {
		settings, err = h.settingQueryRepo.FindByCategoryID(ctx, query.CategoryID)
	} else {
		settings, err = h.settingQueryRepo.FindAll(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}

	// 转换为 DTO
	settingResponses := make([]*SettingDTO, 0, len(settings))
	for _, s := range settings {
		settingResponses = append(settingResponses, ToSettingDTO(s))
	}

	return settingResponses, nil
}
