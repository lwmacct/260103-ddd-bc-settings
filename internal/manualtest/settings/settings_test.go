package settings_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	manualtest "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/manualtest"
)

const baseURL = "http://localhost:8080"

// =============================================================================
// 公开 API 测试（无需认证）
// =============================================================================

// TestPublicSettings 测试获取公开配置列表
func TestPublicSettings(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 调用公开 API - 无需认证（返回扁平结构）
	result, err := manualtest.Get[[]map[string]any](c, "/api/public/settings", nil)
	require.NoError(t, err, "获取公开配置列表失败")

	// 公开配置可能为空（取决于 seeder 数据）
	if result == nil || len(*result) == 0 {
		t.Log("公开配置列表为空（没有标记为 public 的设置）")
		return
	}

	t.Logf("找到 %d 个公开配置项", len(*result))

	// 验证扁平结构：每个 item 只包含 key, value, label
	for _, setting := range *result {
		// 公开 API 的 setting 只应包含 key, value, label
		assert.Contains(t, setting, "key", "Setting 应该有 key 字段")
		assert.Contains(t, setting, "value", "Setting 应该有 value 字段")
		assert.Contains(t, setting, "label", "Setting 应该有 label 字段")

		// 公开 API 不应暴露敏感字段
		assert.NotContains(t, setting, "visible_at", "公开 API 不应暴露 visible_at")
		assert.NotContains(t, setting, "configurable_at", "公开 API 不应暴露 configurable_at")
		assert.NotContains(t, setting, "ui_config", "公开 API 不应暴露 ui_config")
		assert.NotContains(t, setting, "validation", "公开 API 不应暴露 validation")
		assert.NotContains(t, setting, "input_type", "公开 API 不应暴露 input_type")

		t.Logf("公开设置: Key=%s, Label=%s", setting["key"], setting["label"])
	}
}

// TestPublicSettingsWithCategoryFilter 测试按分类过滤公开配置
func TestPublicSettingsWithCategoryFilter(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 使用 category_key 参数过滤（扁平结构不包含 category 字段）
	result, err := manualtest.Get[[]map[string]any](c, "/api/public/settings", map[string]string{
		"category_key": "general",
	})
	require.NoError(t, err, "获取 general 分类的公开配置失败")

	if result == nil || len(*result) == 0 {
		t.Log("general 分类没有公开配置")
		return
	}

	// 验证返回的是扁平结构（只有 key, value, label）
	for _, setting := range *result {
		assert.Contains(t, setting, "key", "Setting 应该有 key 字段")
		assert.Contains(t, setting, "value", "Setting 应该有 value 字段")
	}

	t.Logf("general 分类返回 %d 个公开配置项", len(*result))
}

// TestPublicSettingsNoAuth 测试公开 API 无需认证
func TestPublicSettingsNoAuth(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)
	// 明确不设置任何 token

	// 公开 API 应该可以访问
	_, err := manualtest.Get[[]map[string]any](c, "/api/public/settings", nil)
	assert.NoError(t, err, "公开 API 应该无需认证即可访问")
}

// =============================================================================
// 管理员 API 测试（需要认证）
// =============================================================================

// TestGetCategories 测试获取配置分类列表
func TestGetCategories(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 调用 API - 使用 Get 而非 GetList，因为该接口返回简单列表而非分页数据
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings/categories", nil)
	require.NoError(t, err, "获取分类列表失败")

	// 验证结果
	assert.NotEmpty(t, result, "分类列表不应为空")

	t.Logf("找到 %d 个分类", len(*result))

	// 验证至少有默认分类
	hasGeneral := false
	for _, cat := range *result {
		if key, ok := cat["key"].(string); ok && key == "general" {
			hasGeneral = true
			break
		}
	}
	assert.True(t, hasGeneral, "应包含 general 分类")
}

// TestCreateCategory 测试创建配置分类
func TestCreateCategory(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 创建测试分类
	category := manualtest.CreateTestSettingCategory(t, c, "test_cat")

	// 验证创建成功
	assert.NotNil(t, category, "分类不应为 nil")

	if category != nil {
		t.Logf("创建分类成功: Key=%s, Label=%s", category.Key, category.Label)
	}
}

// TestSettingsAPIPerformance 测试 API 响应时间
func TestSettingsAPIPerformance(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 测试列表接口响应时间
	start := time.Now()
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings/categories", nil)
	duration := time.Since(start)

	require.NoError(t, err, "获取分类列表失败")
	assert.NotEmpty(t, result, "应返回分类数据")
	assert.Less(t, duration, 500*time.Millisecond, "列表接口应在 500ms 内响应")

	t.Logf("API 响应时间: %v, 返回 %d 个分类", duration, len(*result))
}

// TestGetSettings 测试获取配置列表（扁平结构）
func TestGetSettings(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 调用 API - 获取扁平结构的配置列表
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置列表失败")

	// 验证结果
	assert.NotEmpty(t, result, "配置列表不应为空")

	t.Logf("找到 %d 个配置项", len(*result))

	// 验证扁平结构：每个 item 包含 category/group 字段供前端分组
	for i, setting := range *result {
		// 检查必要字段
		key, ok := setting["key"].(string)
		assert.True(t, ok, "Setting 应该有 key 字段")
		assert.NotEmpty(t, key, "Setting key 不应为空")

		// 扁平结构中每个 item 包含 category 和 group 字段
		category, ok := setting["category"].(string)
		assert.True(t, ok, "Setting 应该有 category 字段")
		assert.NotEmpty(t, category, "Setting category 不应为空")

		group, ok := setting["group"].(string)
		assert.True(t, ok, "Setting 应该有 group 字段")
		assert.NotEmpty(t, group, "Setting group 不应为空")

		// 只打印前 3 个作为示例
		if i < 3 {
			t.Logf("配置项: Key=%s, Category=%s, Group=%s", key, category, group)
		}
	}
}

// TestGetSettingByKey 测试获取单个配置详情
func TestGetSettingByKey(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 先获取配置列表，找到一个有效的 key（扁平结构）
	settingsList, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置列表失败")
	require.NotEmpty(t, settingsList, "至少需要一个配置")

	// 从扁平列表中获取第一个配置的 key
	firstSetting := (*settingsList)[0]
	testKey, ok := firstSetting["key"].(string)
	require.True(t, ok, "Setting 应该有 key 字段")

	// 调用 API 获取配置详情
	result, err := manualtest.Get[map[string]any](c, "/api/admin/settings/"+testKey, nil)
	require.NoError(t, err, "获取配置详情失败")

	// 验证结果
	assert.NotNil(t, result, "配置不应为 nil")
	assert.Equal(t, testKey, (*result)["key"], "配置 key 应该匹配")
	assert.NotEmpty(t, (*result)["label"], "配置 label 不应为空")

	t.Logf("配置详情: Key=%s, Label=%s, ValueType=%v",
		(*result)["key"], (*result)["label"], (*result)["value_type"])
}

// TestGetCategoryByID 测试获取单个分类详情
func TestGetCategoryByID(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 先获取分类列表，找到一个有效的 ID
	categories, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings/categories", nil)
	require.NoError(t, err, "获取分类列表失败")
	require.NotEmpty(t, categories, "至少需要一个分类")

	// 获取第一个分类的 ID
	firstCategory := (*categories)[0]
	categoryID, ok := firstCategory["id"].(float64)
	require.True(t, ok, "分类 ID 应该是数字")

	// 调用 API 获取分类详情
	result, err := manualtest.Get[map[string]any](c, fmt.Sprintf("/api/admin/settings/categories/%d", int(categoryID)), nil)
	require.NoError(t, err, "获取分类详情失败")

	// 验证结果
	assert.NotNil(t, result, "分类不应为 nil")
	assert.Equal(t, categoryID, (*result)["id"], "分类 ID 应该匹配")
	assert.NotEmpty(t, (*result)["key"], "分类 key 不应为空")
	assert.NotEmpty(t, (*result)["label"], "分类 label 不应为空")

	t.Logf("分类详情: ID=%v, Key=%s, Label=%s",
		(*result)["id"], (*result)["key"], (*result)["label"])
}

// TestSettingsFlatStructure 测试配置扁平结构的正确性
func TestSettingsFlatStructure(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 获取配置列表（扁平结构）
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置列表失败")
	require.NotEmpty(t, result, "配置列表不应为空")

	totalSettings := len(*result)

	// 统计分类和分组
	categories := make(map[string]bool)
	groups := make(map[string]bool)

	for _, setting := range *result {
		if cat, ok := setting["category"].(string); ok {
			categories[cat] = true
		}
		if grp, ok := setting["group"].(string); ok {
			groups[grp] = true
		}
	}

	t.Logf("扁平结构统计: %d 个分类, %d 个分组, %d 个配置项",
		len(categories), len(groups), totalSettings)

	// 验证数据合理性
	assert.NotEmpty(t, categories, "至少应该有一个分类")
	assert.NotEmpty(t, groups, "至少应该有一个分组")
	assert.Greater(t, totalSettings, 0, "至少应该有一个配置项")
}

// =============================================================================
// 公开 API 与管理员 API 对比测试
// =============================================================================

// TestPublicVsAdminAPIComparison 测试公开 API 和管理员 API 的差异
func TestPublicVsAdminAPIComparison(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 获取管理员 API 的配置列表（扁平结构）
	adminResult, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取管理员配置列表失败")

	// 获取公开 API 的配置列表（扁平结构）
	publicResult, err := manualtest.Get[[]map[string]any](c, "/api/public/settings", nil)
	require.NoError(t, err, "获取公开配置列表失败")

	// 统计管理员 API 中的公开配置数量
	adminPublicCount := 0
	for _, setting := range *adminResult {
		if visibleAt, ok := setting["visible_at"].(string); ok && visibleAt == "public" {
			adminPublicCount++
		}
	}

	// 统计公开 API 返回的配置数量
	publicCount := 0
	if publicResult != nil {
		publicCount = len(*publicResult)
	}

	t.Logf("管理员 API 中公开配置数量: %d", adminPublicCount)
	t.Logf("公开 API 返回的配置数量: %d", publicCount)

	// 公开 API 返回的数量应该等于管理员 API 中 visible_at=public 的数量
	assert.Equal(t, adminPublicCount, publicCount, "公开 API 返回的配置数量应该与管理员 API 中的公开配置数量一致")
}

// =============================================================================
// 性能测试
// =============================================================================

// TestPublicAPIPerformance 测试公开 API 响应时间
func TestPublicAPIPerformance(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 测试公开 API 响应时间
	start := time.Now()
	_, err := manualtest.Get[[]map[string]any](c, "/api/public/settings", nil)
	duration := time.Since(start)

	require.NoError(t, err, "获取公开配置列表失败")
	assert.Less(t, duration, 200*time.Millisecond, "公开 API 应在 200ms 内响应（可能命中缓存）")

	t.Logf("公开 API 响应时间: %v", duration)
}

// TestPublicAPICaching 测试公开 API 缓存效果
func TestPublicAPICaching(t *testing.T) {
	manualtest.SkipIfNotManual(t)
	c := manualtest.NewClient(baseURL)

	// 第一次请求（可能未命中缓存）
	start1 := time.Now()
	_, err := manualtest.Get[[]map[string]any](c, "/api/public/settings", nil)
	duration1 := time.Since(start1)
	require.NoError(t, err, "第一次请求失败")

	// 第二次请求（应该命中缓存）
	start2 := time.Now()
	_, err = manualtest.Get[[]map[string]any](c, "/api/public/settings", nil)
	duration2 := time.Since(start2)
	require.NoError(t, err, "第二次请求失败")

	t.Logf("第一次请求: %v, 第二次请求: %v", duration1, duration2)

	// 第二次请求应该更快（或至少不慢太多）
	// 注意：这个断言可能因网络延迟等因素而不稳定
	if duration1 > 50*time.Millisecond {
		assert.Less(t, duration2, duration1, "第二次请求应该比第一次快（缓存命中）")
	}
}
