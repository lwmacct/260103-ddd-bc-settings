package handler

import (
	"go.uber.org/fx"
)

// HandlerModule 提供 Settings 模块的 HTTP Handler。
var HandlerModule = fx.Module("settings.handler",
	fx.Provide(
		NewSettingHandler,
	),
)
