// Package apitest 提供 Settings BC 特定的测试资源工厂函数。
//
// 本包只包含 BC 特定的测试辅助函数，通用测试功能由 apitest 包提供。
//
// 导入共享包：
//
//	import (
//	    apitest "github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
//	    settingsapitest "github.com/lwmacct/260103-ddd-settings-bc/internal/apitest/settings"
//	)
//
// 使用方式：
//
//	// 创建测试客户端（使用 apitest）
//	c := apitest.NewClient("http://localhost:8080")
//
//	// 创建测试资源（使用 settings apitest factory）
//	setting := settingsapitest.CreateTestSetting(t, c, "testprefix")
//	category := settingsapitest.CreateTestSettingCategory(t, c, "testcat")
//
//	// 发送 HTTP 请求（使用 apitest）
//	result, err := apitest.Get[[]SettingDTO](c, "/api/admin/settings", nil)
//
// 运行测试需要设置 API_TEST 环境变量：
//
//	API_TEST=1 go test -v -count=1 ./internal/apitest/settings/...
package apitest
