package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// CreateCategoryHandler 创建配置分类命令处理器
type CreateCategoryHandler struct {
	commandRepo   setting.SettingCategoryCommandRepository
	queryRepo     setting.SettingCategoryQueryRepository
	settingsCache SettingsCacheService
}

// NewCreateCategoryHandler 创建 CreateCategoryHandler 实例
func NewCreateCategoryHandler(
	commandRepo setting.SettingCategoryCommandRepository,
	queryRepo setting.SettingCategoryQueryRepository,
	settingsCache SettingsCacheService,
) *CreateCategoryHandler {
	return &CreateCategoryHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
		settingsCache: settingsCache,
	}
}

// Handle 处理创建配置分类命令
func (h *CreateCategoryHandler) Handle(ctx context.Context, cmd CreateCategoryCommand) (*CreateCategoryResultDTO, error) {
	// 1. 验证 Key 是否已存在
	exists, err := h.queryRepo.ExistsByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing category: %w", err)
	}
	if exists {
		return nil, errors.New("category key already exists")
	}

	// 2. 设置默认图标
	icon := cmd.Icon
	if icon == "" {
		icon = "mdi-cog"
	}

	// 3. 创建分类实体
	category := &setting.SettingCategory{
		Key:   cmd.Key,
		Label: cmd.Label,
		Icon:  icon,
		Order: cmd.Order,
	}

	// 4. 验证实体
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// 5. 保存分类
	if err := h.commandRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// 6. 失效 Settings 缓存（Category 变更影响所有设置页）
	_ = h.settingsCache.DeleteAll(ctx)
	_ = h.settingsCache.DeleteAllCategories(ctx)

	return &CreateCategoryResultDTO{
		ID: category.ID,
	}, nil
}
