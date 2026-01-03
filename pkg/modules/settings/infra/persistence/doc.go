// Package persistence 提供 Setting 领域模型的数据持久化实现。
//
// 本包实现了 Domain 层定义的 Repository 接口：
//   - CommandRepository：配置定义写操作
//   - QueryRepository：配置定义读操作
//   - SettingCategoryCommandRepository：配置分类写操作
//   - SettingCategoryQueryRepository：配置分类读操作
//
// 技术实现：
//   - GORM 作为 ORM 框架
//   - PostgreSQL 作为数据库
//   - 支持缓存装饰器（通过 Cache 服务）
//
// 依赖倒置：本包依赖 Domain 层接口，不引入循环依赖。
package persistence
