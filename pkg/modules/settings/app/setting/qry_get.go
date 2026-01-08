package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// GetHandler 获取设置查询处理器
type GetHandler struct {
	settingQueryRepo setting.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例
func NewGetHandler(settingQueryRepo setting.QueryRepository) *GetHandler {
	return &GetHandler{
		settingQueryRepo: settingQueryRepo,
	}
}

// Handle 处理获取设置查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*SettingDTO, error) {
	settingEntity, err := h.settingQueryRepo.FindByKey(ctx, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}

	if settingEntity == nil {
		return nil, errors.New("setting not found")
	}

	return ToSettingDTO(settingEntity), nil
}
