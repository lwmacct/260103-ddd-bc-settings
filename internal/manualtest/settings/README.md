# Settings Module Integration Tests

## 测试状态 (2026-01-06 - 扁平结构重构)

### API 响应结构变更

**重要**: Settings API 已从层级结构改为扁平结构

| 版本 | 结构                                       | 说明     |
| ---- | ------------------------------------------ | -------- |
| 旧版 | `Category → Group → Settings[]`            | 层级嵌套 |
| 新版 | `Settings[]` (每个 item 含 category/group) | 扁平结构 |

**扁平结构优势**：

- 前端可自由组织展示方式
- 易于搜索/过滤单个配置项
- 支持分页
- Team 场景可携带继承信息

### ✅ 通过的测试

| 测试名称                       | 端点                                   | 状态 | 说明                  |
| ------------------------------ | -------------------------------------- | ---- | --------------------- |
| TestGetCategories              | GET /api/admin/settings/categories     | PASS | 获取分类列表          |
| TestCreateCategory             | POST /api/admin/settings/categories    | PASS | 创建测试分类          |
| TestSettingsAPIPerformance     | GET /api/admin/settings/categories     | PASS | 响应时间 < 500ms      |
| TestGetSettings                | GET /api/admin/settings                | PASS | **扁平结构** 配置列表 |
| TestGetSettingByKey            | GET /api/admin/settings/:key           | PASS | 获取配置详情          |
| TestGetCategoryByID            | GET /api/admin/settings/categories/:id | PASS | 获取分类详情          |
| TestSettingsFlatStructure      | GET /api/admin/settings                | PASS | **扁平结构验证**      |
| TestPublicSettings             | GET /api/public/settings               | PASS | 公开配置（扁平）      |
| TestPublicSettingsNoAuth       | GET /api/public/settings               | PASS | 无需认证              |
| TestPublicVsAdminAPIComparison | -                                      | PASS | 公开/管理员数量一致性 |

## 测试覆盖的端点

### Settings API

- [x] GET /api/admin/settings (扁平结构) ✅
- [x] GET /api/admin/settings/:key ✅
- [ ] POST /api/admin/settings
- [ ] PUT /api/admin/settings/:key
- [ ] DELETE /api/admin/settings/:key
- [ ] POST /api/admin/settings/batch

### Settings Category API

- [x] GET /api/admin/settings/categories ✅
- [x] GET /api/admin/settings/categories/:id ✅
- [x] POST /api/admin/settings/categories ✅
- [ ] PUT /api/admin/settings/categories/:id
- [ ] DELETE /api/admin/settings/categories/:id

### Public Settings API

- [x] GET /api/public/settings (扁平结构) ✅

## 运行测试

```bash
# 运行所有 Settings 测试
MANUAL=1 go test -v -count=1 ./internal/manualtest/settings/...

# 运行特定测试
MANUAL=1 go test -v -count=1 ./internal/manualtest/settings/... -run TestGetSettings

# 运行扁平结构验证测试
MANUAL=1 go test -v -count=1 ./internal/manualtest/settings/... -run "TestGetSettings|TestSettingsFlatStructure"
```

## API 响应示例

### Admin Settings API (扁平结构)

```json
[
  {
    "key": "general.site_name",
    "category": "general",
    "group": "basic",
    "value": "My Site",
    "default_value": "My Site",
    "is_customized": false,
    "label": "站点名称",
    "visible_at": "user",
    "configurable_at": "admin",
    "ui_config": { ... }
  },
  {
    "key": "security.password_min_length",
    "category": "security",
    "group": "密码策略",
    "value": 8,
    ...
  }
]
```

### Public Settings API (精简扁平结构)

```json
[
  {
    "key": "general.site_name",
    "value": "My Site",
    "label": "站点名称"
  }
]
```

## 性能指标

| API 端点                           | 响应时间 | 状态    |
| ---------------------------------- | -------- | ------- |
| GET /api/admin/settings/categories | < 1ms    | ✅ 优秀 |
| GET /api/admin/settings (扁平)     | < 10ms   | ✅ 良好 |
| GET /api/public/settings           | < 200ms  | ✅ 良好 |
