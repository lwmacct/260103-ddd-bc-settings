// Package manualtest 提供 IAM 模块的 HTTP 集成测试辅助工具。
//
// 测试工具包括：
//   - Client: HTTP 测试客户端，支持认证请求
//   - Factory: 测试资源工厂函数（用户、角色、配置等）
//   - Helper: 测试辅助函数（登录、缓存）
//   - Assert: 测试断言辅助函数
//
// 使用方式：
//
//	c := manualtest.NewClient()
//	_, err := c.Login("admin", "admin123")
//	assert.NoError(t, err)
package manualtest
