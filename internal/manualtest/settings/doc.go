// Package manualtest 提供 Settings BC 特定的测试资源工厂函数。
//
// 本包只包含 BC 特定的测试辅助函数，通用测试功能由 apitest 包提供。
//
// 导入共享包：
//
//	import (
//	    apitest "github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
//	    manualtest "github.com/lwmacct/260103-ddd-bc-settings/internal/manualtest/settings"
//	)
//
// 使用方式：
//
//	// 创建测试客户端（使用 apitest）
//	c := apitest.NewClient("http://localhost:8080")
//
//	// 创建测试资源（使用 manualtest factory）
//	setting := manualtest.CreateTestSetting(t, c, "testprefix")
//	category := manualtest.CreateTestSettingCategory(t, c, "testcat")
//
//	// 发送 HTTP 请求（使用 apitest）
//	result, err := apitest.Get[[]SettingDTO](c, "/api/admin/settings", nil)
//
// 运行测试需要设置 API_TEST 环境变量：
//
//	API_TEST=1 go test -v -count=1 ./internal/manualtest/settings/...
package manualtest
