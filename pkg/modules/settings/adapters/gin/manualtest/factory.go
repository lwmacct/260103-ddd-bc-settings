package manualtest

import (
	"testing"

	"github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
)

// CreateTestSetting 创建测试配置，自动注册 Cleanup。
//
// 使用示例：
//
//	setting := manualtest.CreateTestSetting(t, c, "testprefix")
//	// 使用 setting...
//	// 测试结束后自动清理
func CreateTestSetting(t *testing.T, c *Client, prefix string) *setting.SettingDTO {
	result, _ := CreateTestSettingWithCleanupControl(t, c, prefix)
	return result
}

// CreateTestSettingWithCleanupControl 创建测试配置，返回清理控制函数。
//
// 使用场景：测试本身包含删除操作时。
//
// 使用示例：
//
//	setting, markDeleted := manualtest.CreateTestSettingWithCleanupControl(t, c, "testprefix")
//	// ... 测试删除逻辑
//	err := c.Delete("/api/admin/settings/" + setting.Key)
//	// ... 验证删除成功
//	markDeleted() // 标记已删除，跳过 Cleanup
func CreateTestSettingWithCleanupControl(t *testing.T, c *Client, prefix string) (*setting.SettingDTO, func()) {
	key := prefix + "_" + randomString(8)

	createReq := map[string]any{
		"key":           key,
		"default_value": "test_value",
		"category_id":   1, // 假设分类 ID 为 1
		"group":         "test",
		"value_type":    "string",
		"label":         "Test Setting",
		"order":         100,
		"input_type":    "text",
	}

	result, err := Post[setting.SettingDTO](c, "/api/admin/settings", createReq)
	if err != nil {
		t.Fatalf("failed to create test setting: %v", err)
	}

	// 注册清理
	cleanup := func() {
		_ = c.Delete("/api/admin/settings/" + key)
	}

	markDeleted := func() {
		// 将 cleanup 替换为空函数
		cleanup = func() {}
	}

	t.Cleanup(func() {
		cleanup()
	})

	return result, markDeleted
}

// CreateTestSettingCategory 创建测试配置分类，自动注册 Cleanup。
func CreateTestSettingCategory(t *testing.T, c *Client, prefix string) *setting.CategoryDTO {
	result, _ := CreateTestSettingCategoryWithCleanupControl(t, c, prefix)
	return result
}

// CreateTestSettingCategoryWithCleanupControl 创建测试配置分类，返回清理控制函数。
func CreateTestSettingCategoryWithCleanupControl(t *testing.T, c *Client, prefix string) (*setting.CategoryDTO, func()) {
	key := prefix + "_" + randomString(8)

	createReq := map[string]any{
		"key":   key,
		"label": "Test Category",
		"icon":  "mdi-test",
		"order": 100,
	}

	result, err := Post[setting.CategoryDTO](c, "/api/admin/settings/categories", createReq)
	if err != nil {
		t.Fatalf("failed to create test category: %v", err)
	}

	// 注册清理
	cleanup := func() {
		_ = c.Delete("/api/admin/settings/categories/" + key)
	}

	markDeleted := func() {
		cleanup = func() {}
	}

	t.Cleanup(func() {
		cleanup()
	})

	return result, markDeleted
}

// randomString 生成随机字符串用于测试数据。
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
