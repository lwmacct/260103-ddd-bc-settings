package handler

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app"
)

// Handlers 聚合 Settings 模块的所有 HTTP Handler。
//
// 采用聚合结构体模式，而非 fx.Out 导出多个 Handler：
//   - 新增/删除 Handler 只需修改本结构体（1 处）
//   - fx.Out 模式需修改 struct + 构造函数 + 消费方（3 处）
type Handlers struct {
	Setting *SettingHandler
}

// HandlersParams 定义创建 Handlers 所需的依赖。
type HandlersParams struct {
	fx.In

	SettingUseCases *app.SettingUseCases
}

// NewHandlers 创建 Handlers 聚合。
//
// 接收 UseCases 聚合，传递给各个 Handler 构造函数。
func NewHandlers(p HandlersParams) *Handlers {
	return &Handlers{
		Setting: NewSettingHandler(p.SettingUseCases),
	}
}

// HandlerModule 提供 Settings 模块的 HTTP Handler。
var HandlerModule = fx.Module("settings.handler",
	fx.Provide(
		NewHandlers,
	),
)
