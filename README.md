# Settings Bounded Context

基于 Go 的配置管理系统，采用 **DDD 四层架构 + CQRS 模式**。使用 Uber Fx 依赖注入，提供完整的配置定义、配置分类、层级结构展示等功能。

## 项目概述

本项目是从 `ddd-old` 项目中独立出来的 Settings 模块，实现了配置管理的完整功能：

- **配置定义**：支持系统级和用户级配置，支持多种数据类型（string、number、boolean、json）
- **配置分类**：支持按分类组织配置，便于管理
- **层级结构**：Category → Group → Setting 三层聚合，适配前端动态渲染
- **缓存机制**：基于 Redis 的缓存策略，提升查询性能
- **权限控制**：基于 URN 风格的 RBAC 权限系统

## 架构概览

```
pkg/modules/settings/
├── domain/setting/           # 领域层 - 实体、Repository 接口
├── app/setting/              # 应用层 - UseCase Handlers、DTO
├── infra/                    # 基础设施层
│   ├── persistence/          # GORM 仓储
│   ├── cache/                # Redis 缓存服务
│   └── seeds/                # 种子数据
└── adapters/gin/             # 适配器层 - HTTP 接口
    ├── handler/              # HTTP Handler
    ├── routes/               # 路由定义
    └── manualtest/           # 集成测试

internal/
├── bootstrap/                # 应用启动引导
└── container/                # 依赖注入容器（Uber Fx）

cmd/server/
├── main.go                   # CLI 入口（server / db migrate / db reset）
└── docs/                     # Swagger 生成物
```

## 模块职责

| 层级                 | 职责     | 说明                                      |
| -------------------- | -------- | ----------------------------------------- |
| `domain/setting`     | 业务核心 | 实体、值对象、Repository 接口、领域错误   |
| `app/setting`        | 用例编排 | Command/Query Handler、DTO 映射、业务流程 |
| `infra`              | 技术实现 | GORM Repository、Redis 缓存、种子数据     |
| `adapters/gin`       | 外部交互 | HTTP Handler、路由、测试                  |
| `internal/container` | 依赖注入 | Fx 模块组装、生命周期管理                 |

## 开发工作流

### 运行项目

```bash
# 开发模式（热重载）
air

# 直接运行
go run cmd/server/main.go

# 构建
go build -o /dev/null ./...
```

### 数据库管理

```bash
# 执行迁移
go run cmd/server/main.go db migrate

# 重置数据库（删表+重建+种子数据）
go run cmd/server/main.go db reset

# 执行种子数据
go run cmd/server/main.go db seed
```

### 测试

```bash
# 单元测试
go test ./...

# API 集成测试（需要服务运行）
API_TEST=1 go test -v -count=1 ./internal/manualtest/...
```

### 代码质量

```bash
# Lint 检查（仅新修改）
golangci-lint run --new

# 完整 Lint
golangci-lint run

# 预提交检查
pre-commit run --all
```

## API 端点

### Setting API

| 方法   | 路径                      | 说明             | 权限                  |
| ------ | ------------------------- | ---------------- | --------------------- |
| GET    | /api/admin/settings       | 配置列表（层级） | admin:settings:list   |
| GET    | /api/admin/settings/{key} | 配置详情         | admin:settings:get    |
| POST   | /api/admin/settings       | 创建配置         | admin:settings:create |
| PUT    | /api/admin/settings/{key} | 更新配置         | admin:settings:update |
| DELETE | /api/admin/settings/{key} | 删除配置         | admin:settings:delete |
| POST   | /api/admin/settings/batch | 批量更新配置     | admin:settings:update |

### Category API

| 方法   | 路径                                | 说明     | 权限                             |
| ------ | ----------------------------------- | -------- | -------------------------------- |
| GET    | /api/admin/settings/categories      | 分类列表 | admin:settings:categories:list   |
| GET    | /api/admin/settings/categories/{id} | 分类详情 | admin:settings:categories:get    |
| POST   | /api/admin/settings/categories      | 创建分类 | admin:settings:categories:create |
| PUT    | /api/admin/settings/categories/{id} | 更新分类 | admin:settings:categories:update |
| DELETE | /api/admin/settings/categories/{id} | 删除分类 | admin:settings:categories:delete |

## 数据库设计

### 核心表

**settings 表**（配置定义）：

- 字段：key, default_value (JSONB), scope, public, category_id, group, order, value_type, label, input_type, validation, ui_config (JSONB)
- 索引：idx_settings_category_sort, idx_settings_scope, idx_settings_visible_to_user

**setting_categories 表**（配置分类）：

- 字段：key, label, icon, sort_order
- 索引：主键 id，唯一键 key

## 技术栈

- **Go 1.25.5**：开发语言
- **Gin**：HTTP 框架
- **GORM**：ORM 框架
- **Uber Fx**：依赖注入
- **Redis**：缓存（RedisJSON）
- **PostgreSQL**：数据库
- **Swagger**：API 文档

## 架构原则

**模块完全自治**：Settings 模块包含完整的四层架构，所有业务相关代码（缓存、种子、测试）都内聚在模块内部。

**依赖方向规则**：

```
Platform (外部提供) → Settings Module → Container (组装)
```

**外部依赖**：

- ✅ Platform 层通过外部依赖引入（DB、Redis、EventBus）
- ✅ Settings 模块只依赖 Platform 接口（依赖倒置）
- ✅ 可独立编译和测试

## License

MIT License
