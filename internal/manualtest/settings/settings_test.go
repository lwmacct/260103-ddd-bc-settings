package settings_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	manualtest "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/adapters/gin/manualtest"
)

// TestGetCategories 测试获取配置分类列表
func TestGetCategories(t *testing.T) {
	c := manualtest.NewClient("http://localhost:8080")

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
	c := manualtest.NewClient("http://localhost:8080")

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
	c := manualtest.NewClient("http://localhost:8080")

	// 测试列表接口响应时间
	start := time.Now()
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings/categories", nil)
	duration := time.Since(start)

	require.NoError(t, err, "获取分类列表失败")
	assert.NotEmpty(t, result, "应返回分类数据")
	assert.Less(t, duration, 500*time.Millisecond, "列表接口应在 500ms 内响应")

	t.Logf("API 响应时间: %v, 返回 %d 个分类", duration, len(*result))
}

// TestGetSettings 测试获取配置列表（层级结构）
func TestGetSettings(t *testing.T) {
	c := manualtest.NewClient("http://localhost:8080")

	// 调用 API - 获取层级结构的配置列表
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置列表失败")

	// 验证结果
	assert.NotEmpty(t, result, "配置列表不应为空")

	t.Logf("找到 %d 个分类配置", len(*result))

	// 验证层级结构：应该有 Category → Groups → Settings 的三层结构
	for _, category := range *result {
		// 检查 Category 字段
		categoryKey, ok := category["category"].(string)
		assert.True(t, ok, "Category 应该有 key 字段")
		assert.NotEmpty(t, categoryKey, "Category key 不应为空")

		// 检查 Groups 字段
		groups, ok := category["groups"].([]any)
		assert.True(t, ok, "Category 应该有 groups 数组")
		assert.NotEmpty(t, groups, "Groups 不应为空")

		// 检查第一个 Group 的结构
		if len(groups) > 0 {
			firstGroup := groups[0].(map[string]any)
			assert.Contains(t, firstGroup, "name", "Group 应该有 name 字段")
			assert.Contains(t, firstGroup, "settings", "Group 应该有 settings 数组")

			t.Logf("Category: %s, 第一个 Group: %v", categoryKey, firstGroup["name"])
		}
	}
}

// TestGetSettingByKey 测试获取单个配置详情
func TestGetSettingByKey(t *testing.T) {
	c := manualtest.NewClient("http://localhost:8080")

	// 先获取配置列表，找到一个有效的 key
	settingsList, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置列表失败")
	require.NotEmpty(t, settingsList, "至少需要一个配置")

	// 从第一个分类的第一个分组中获取第一个配置的 key
	firstCategory := (*settingsList)[0]
	groups, ok := firstCategory["groups"].([]any)
	require.True(t, ok, "Category 应该有 groups 数组")
	require.NotEmpty(t, groups, "Groups 不应为空")

	firstGroup := groups[0].(map[string]any)
	settings, ok := firstGroup["settings"].([]any)
	require.True(t, ok, "Group 应该有 settings 数组")
	require.NotEmpty(t, settings, "Settings 不应为空")

	firstSetting := settings[0].(map[string]any)
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
	c := manualtest.NewClient("http://localhost:8080")

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

// TestSettingsSchemaStructure 测试配置层级结构的正确性
func TestSettingsSchemaStructure(t *testing.T) {
	c := manualtest.NewClient("http://localhost:8080")

	// 获取完整的配置 schema
	result, err := manualtest.Get[[]map[string]any](c, "/api/admin/settings", nil)
	require.NoError(t, err, "获取配置 schema 失败")

	require.NotEmpty(t, result, "Schema 不应为空")

	totalCategories := len(*result)
	totalGroups := 0
	totalSettings := 0

	// 遍历层级结构统计
	for _, category := range *result {
		groups, ok := category["groups"].([]any)
		require.True(t, ok, "Category 应该有 groups 数组")

		totalGroups += len(groups)

		for _, group := range groups {
			groupMap, ok := group.(map[string]any)
			require.True(t, ok, "Group 应该是 map")

			settings, ok := groupMap["settings"].([]any)
			require.True(t, ok, "Group 应该有 settings 数组")

			totalSettings += len(settings)
		}
	}

	t.Logf("层级结构统计: %d 个分类, %d 个分组, %d 个配置项",
		totalCategories, totalGroups, totalSettings)

	// 验证数据合理性
	assert.Greater(t, totalCategories, 0, "至少应该有一个分类")
	assert.Greater(t, totalGroups, 0, "至少应该有一个分组")
	assert.Greater(t, totalSettings, 0, "至少应该有一个配置项")
}
