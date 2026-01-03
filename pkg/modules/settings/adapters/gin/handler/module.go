package handler

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
)

// HandlerModule 提供 Settings 模块的 HTTP Handler。
var HandlerModule = fx.Module("settings.handler",
	fx.Provide(
		NewAllHandlers,
	),
)

// HandlersResult 批量返回 Handler（使用 fx.Out）。
type HandlersResult struct {
	fx.Out
	Setting *SettingHandler
}

// NewAllHandlers 创建所有 Handler（如果需要返回多个）。
// 当前只有一个 SettingHandler，所以直接使用 NewSettingHandler。
//
// 如果将来需要添加其他 Handler（如 UserSettingHandler），
// 可以使用此函数返回 HandlersResult。
func NewAllHandlers(usecases *setting.SettingUseCases) HandlersResult {
	return HandlersResult{
		Setting: NewSettingHandler(
			usecases.Create,
			usecases.Update,
			usecases.Delete,
			usecases.BatchUpdate,
			usecases.Get,
			usecases.List,
			nil, // TODO: usecases.ListSettings
			usecases.CreateCategory,
			usecases.UpdateCategory,
			usecases.DeleteCategory,
			usecases.GetCategory,
			usecases.ListCategories,
		),
	}
}
