package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// DeleteCategoryHandler 删除配置分类命令处理器
type DeleteCategoryHandler struct {
	commandRepo      setting.SettingCategoryCommandRepository
	queryRepo        setting.SettingCategoryQueryRepository
	settingQueryRepo setting.QueryRepository // 用于检查关联的 Setting
	settingsCache    SettingsCacheService
}

// NewDeleteCategoryHandler 创建 DeleteCategoryHandler 实例
func NewDeleteCategoryHandler(
	commandRepo setting.SettingCategoryCommandRepository,
	queryRepo setting.SettingCategoryQueryRepository,
	settingQueryRepo setting.QueryRepository,
	settingsCache SettingsCacheService,
) *DeleteCategoryHandler {
	return &DeleteCategoryHandler{
		commandRepo:      commandRepo,
		queryRepo:        queryRepo,
		settingQueryRepo: settingQueryRepo,
		settingsCache:    settingsCache,
	}
}

// Handle 处理删除配置分类命令
func (h *DeleteCategoryHandler) Handle(ctx context.Context, cmd DeleteCategoryCommand) error {
	// 1. 验证分类是否存在
	category, err := h.queryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil {
		return errors.New("category not found")
	}

	// 2. 检查是否有关联的 Setting
	settings, err := h.settingQueryRepo.FindByCategoryID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to check associated settings: %w", err)
	}
	if len(settings) > 0 {
		return fmt.Errorf("cannot delete category: %d settings are associated with this category", len(settings))
	}

	// 3. 执行删除
	if err := h.commandRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	// 4. 失效 Settings 缓存（Category 变更影响所有设置页）
	_ = h.settingsCache.DeleteAll(ctx)
	_ = h.settingsCache.DeleteAllCategories(ctx)

	return nil
}
