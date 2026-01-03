package setting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// CreateHandler 创建配置命令处理器
type CreateHandler struct {
	commandRepo   setting.CommandRepository
	queryRepo     setting.QueryRepository
	settingsCache SettingsCacheService
}

// NewCreateHandler 创建 CreateHandler 实例
func NewCreateHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	settingsCache SettingsCacheService,
) *CreateHandler {
	return &CreateHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
		settingsCache: settingsCache,
	}
}

// Handle 处理创建配置命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*CreateResultDTO, error) {
	// 1. 验证 Key 是否已存在
	exists, err := h.queryRepo.ExistsByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing setting: %w", err)
	}
	if exists {
		return nil, errors.New("setting key already exists")
	}

	// 2. 设置默认值类型
	valueType := cmd.ValueType
	if valueType == "" {
		valueType = setting.ValueTypeString
	}

	// 3. 创建配置定义实体
	s := &setting.Setting{
		Key:          cmd.Key,
		DefaultValue: cmd.DefaultValue,
		CategoryID:   cmd.CategoryID,
		Group:        cmd.Group,
		ValueType:    valueType,
		Label:        cmd.Label,
		UIConfig:     cmd.UIConfig,
		Order:        cmd.Order,
	}

	// 4. 保存配置定义
	if err := h.commandRepo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("failed to create setting: %w", err)
	}

	// 5. 失效 Settings 缓存
	if h.settingsCache != nil {
		if err := h.settingsCache.DeleteAdminSettingsAll(ctx); err != nil {
			slog.Warn("admin settings cache invalidation failed", "key", cmd.Key, "error", err.Error())
		}
	}

	return &CreateResultDTO{
		ID: s.ID,
	}, nil
}
