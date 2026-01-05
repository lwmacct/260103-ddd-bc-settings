package setting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// UpdateHandler 更新配置命令处理器
type UpdateHandler struct {
	commandRepo   setting.CommandRepository
	queryRepo     setting.QueryRepository
	settingsCache SettingsCacheService
}

// NewUpdateHandler 创建 UpdateHandler 实例
func NewUpdateHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	settingsCache SettingsCacheService,
) *UpdateHandler {
	return &UpdateHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
		settingsCache: settingsCache,
	}
}

// Handle 处理更新配置命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*SettingDTO, error) {
	// 1. 查询配置定义
	def, err := h.queryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, errors.New("setting not found")
	}

	// 2. 基于 ValueType 的类型校验
	if err := def.ValidateValue(cmd.DefaultValue); err != nil {
		return nil, fmt.Errorf("value validation failed: %w", err)
	}

	// 3. 基于 InputType 的格式校验（email/url/password 等）
	if err := def.ValidateByInputType(cmd.DefaultValue); err != nil {
		return nil, err
	}

	// 4. 自定义 Validation 规则校验（JSON Logic）
	// TODO: 实现完整的验证逻辑（需要 Validator 服务）
	// 当前跳过 JSON Logic 验证

	// 5. 更新字段
	def.DefaultValue = cmd.DefaultValue
	if cmd.Label != "" {
		def.Label = cmd.Label
	}
	if cmd.UIConfig != "" {
		def.UIConfig = cmd.UIConfig
	}
	if cmd.Order != 0 {
		def.Order = cmd.Order
	}

	// 6. 保存更新
	if err := h.commandRepo.Update(ctx, def); err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	// 7. 失效 Settings 缓存
	if h.settingsCache != nil {
		if err := h.settingsCache.DeleteAdminSettingsAll(ctx); err != nil {
			slog.Warn("admin settings cache invalidation failed", "key", cmd.Key, "error", err.Error())
		}
	}

	return ToSettingDTO(def), nil
}
