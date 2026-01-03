// Package setting 提供配置管理的领域模型。
//
// 核心实体：
//   - Setting：配置定义实体，包含配置的 Schema 和默认值
//   - SettingCategory：配置分类实体，用于组织配置项
//
// Repository 接口：
//   - CommandRepository：配置定义写操作（Create/Update/Delete/BatchUpsert）
//   - QueryRepository：配置定义读操作（FindByKey/FindByCategory/FindVisibleToUser）
//   - SettingCategoryCommandRepository：配置分类写操作
//   - SettingCategoryQueryRepository：配置分类读操作
//
// 领域概念：
//   - Scope：system（系统级，全局唯一）、user（用户级，可覆盖）
//   - ValueType：string、number、boolean、json
//   - InputType：text、email、url、password、select、number、checkbox、textarea
//   - Validation：JSON Logic 格式的自定义校验规则
//   - UIConfig：前端展示配置（hint、options、depends_on）
package setting
