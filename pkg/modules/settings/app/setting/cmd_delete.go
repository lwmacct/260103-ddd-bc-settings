package setting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// DeleteHandler 删除配置命令处理器
type DeleteHandler struct {
	commandRepo   setting.CommandRepository
	queryRepo     setting.QueryRepository
	settingsCache SettingsCacheService
}

// NewDeleteHandler 创建 DeleteHandler 实例
func NewDeleteHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	settingsCache SettingsCacheService,
) *DeleteHandler {
	return &DeleteHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
		settingsCache: settingsCache,
	}
}

// Handle 处理删除配置命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 查询配置定义
	def, err := h.queryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return errors.New("setting not found")
	}

	// 2. 删除配置定义
	if err := h.commandRepo.Delete(ctx, cmd.Key); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	// 3. 失效 Settings 缓存
	if h.settingsCache != nil {
		if err := h.settingsCache.DeleteAdminSettingsAll(ctx); err != nil {
			slog.Warn("admin settings cache invalidation failed", "key", cmd.Key, "error", err.Error())
		}
	}

	return nil
}
