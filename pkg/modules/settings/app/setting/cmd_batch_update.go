package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// BatchUpdateHandler 批量更新配置命令处理器
type BatchUpdateHandler struct {
	commandRepo   setting.CommandRepository
	queryRepo     setting.QueryRepository
}

// NewBatchUpdateHandler 创建 BatchUpdateHandler 实例
func NewBatchUpdateHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
) *BatchUpdateHandler {
	return &BatchUpdateHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
	}
}

// Handle 处理批量更新配置命令
func (h *BatchUpdateHandler) Handle(ctx context.Context, cmd BatchUpdateCommand) error {
	if len(cmd.Settings) == 0 {
		return nil
	}

	// 1. 提取所有 keys
	keys := make([]string, len(cmd.Settings))
	keyValueMap := make(map[string]any, len(cmd.Settings))
	for i, item := range cmd.Settings {
		keys[i] = item.Key
		keyValueMap[item.Key] = item.Value
	}

	// 2. 批量查询所有现有配置定义
	existingDefs, err := h.queryRepo.FindByKeys(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to find settings: %w", err)
	}

	// 3. 构建 key -> def 映射
	existingMap := make(map[string]*setting.Setting, len(existingDefs))
	for _, d := range existingDefs {
		existingMap[d.Key] = d
	}

	// 4. 验证所有 key 存在并更新值
	settings := make([]*setting.Setting, 0, len(cmd.Settings))
	for _, key := range keys {
		existing, ok := existingMap[key]
		if !ok {
			return fmt.Errorf("setting key %s does not exist", key)
		}
		existing.DefaultValue = keyValueMap[key]
		settings = append(settings, existing)
	}

	// 5. 批量更新（Repository 装饰器会自动失效缓存）
	if err := h.commandRepo.BatchUpsert(ctx, settings); err != nil {
		return fmt.Errorf("failed to batch update settings: %w", err)
	}

	return nil
}
