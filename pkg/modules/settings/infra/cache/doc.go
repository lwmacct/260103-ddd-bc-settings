// Package cache 提供 Settings 模块的缓存服务实现。
//
// 本包实现了 Application 层定义的 SettingsCacheService 接口：
//   - 管理员 Settings 缓存（按 category）
//   - 用户 Settings 缓存（按 userID + category）
//   - 分类列表缓存
//   - Category 实体缓存
//
// 技术实现：
//   - Redis 作为缓存存储
//   - RedisJSON（JSONGET/JSONSET）存储 JSON 数据
//   - TTL：30 分钟
//   - 支持 SCAN 批量失效
//
// 缓存策略：
//   - 写操作后异步失效相关缓存
//   - 读操作先查缓存，未命中则查数据库
package cache
