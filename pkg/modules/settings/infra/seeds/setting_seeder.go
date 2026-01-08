package seeds

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/infra/persistence"
)

// SettingSeeder 系统设置种子数据
type SettingSeeder struct{}

// Seed 执行系统设置种子数据填充
func (s *SettingSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 1. 查询所有 Category ID
	categoryIDs, err := s.loadCategoryIDs(db)
	if err != nil {
		return fmt.Errorf("load category IDs: %w", err)
	}

	// 2. 检查必需的 Category，缺失时记录警告
	requiredCategories := []string{"general", "security", "email", "oauth", "notification", "backup"}
	for _, key := range requiredCategories {
		if _, ok := categoryIDs[key]; !ok {
			slog.Warn("required category not found, skipping dependent settings", "category", key)
		}
	}

	// 3. 构建配置定义（使用 CategoryID）
	definitions := s.buildDefinitions(categoryIDs)

	// 4. 过滤掉 Category ID 为 0 的定义（Category 不存在）
	filteredDefinitions := s.filterValidDefinitions(definitions)

	// 5. 批量插入/更新
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"input_type", "validation", "ui_config", // hint 已移入 ui_config
			"group", "order", "label",
			"visible_at", "configurable_at", // 新字段
		}), // 仅更新 UI 元数据，保留已有的 DefaultValue（通过 API 修改的值）
	}).Create(&filteredDefinitions)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded setting definitions", "total", len(definitions), "valid", len(filteredDefinitions), "inserted", result.RowsAffected)
	return nil
}

// Name 返回 seeder 名称
func (s *SettingSeeder) Name() string {
	return "SettingSeeder"
}

// loadCategoryIDs 从数据库加载 Category key -> ID 映射
func (s *SettingSeeder) loadCategoryIDs(db *gorm.DB) (map[string]uint, error) {
	var categories []persistence.SettingCategoryModel
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make(map[string]uint, len(categories))
	for _, cat := range categories {
		result[cat.Key] = cat.ID
	}
	return result, nil
}

// filterValidDefinitions 过滤掉 Category ID 为 0 的定义（Category 不存在）
func (s *SettingSeeder) filterValidDefinitions(definitions []persistence.SettingModel) []persistence.SettingModel {
	filtered := make([]persistence.SettingModel, 0, len(definitions))
	for _, def := range definitions {
		if def.CategoryID == 0 {
			slog.Warn("skipping setting with missing category", "key", def.Key)
			continue
		}
		filtered = append(filtered, def)
	}
	return filtered
}

// buildDefinitions 构建配置定义列表
//
// 排序规则：
//  1. 分组通过 Order 范围控制排序（第1组: 1-99, 第2组: 100-199, 第3组: 200-299）
//  2. 同组内按 Order 排序（10, 20, 30...）
//  3. Group 字段直接存储显示标签（如 "基本设置"）
//
// Order 值规划示例：
//   - 基本设置组: 10, 20, 30    → min=10  → 排第1
//   - 本地化组:   100, 110      → min=100 → 排第2
//   - 外观组:     200           → min=200 → 排第3
func (s *SettingSeeder) buildDefinitions(categoryIDs map[string]uint) []persistence.SettingModel {
	return []persistence.SettingModel{
		// ==================== General 常规设置 ====================
		// 基本设置 分组 (Order 1-99)
		{
			Key:            "general.site_name",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["general"],
			Group:          "基本设置",
			VisibleAt:      "public",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "站点名称",
			Order:          10,
			InputType:      "text",
			Validation:     `{"min_length":1,"max_length":100}`,
			UIConfig:       datatypes.JSON(`{"hint":"显示在浏览器标签和页面标题中"}`),
		},
		{
			Key:            "general.site_url",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["general"],
			Group:          "基本设置",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "站点 URL",
			Order:          20,
			InputType:      "url",
			UIConfig:       datatypes.JSON(`{"hint":"站点完整 URL，如 https://example.com"}`),
		},
		{
			Key:            "general.admin_email",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["general"],
			Group:          "基本设置",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "管理员邮箱",
			Order:          30,
			InputType:      "email",
			UIConfig:       datatypes.JSON(`{"hint":"用于接收系统通知和报警邮件"}`),
		},
		// 本地化 分组 (Order 100-199)
		{
			Key:            "general.timezone",
			DefaultValue:   datatypes.JSON(`"Asia/Shanghai"`),
			CategoryID:     categoryIDs["general"],
			Group:          "本地化",
			VisibleAt:      "user",
			ConfigurableAt: "team", // Team 可为用户设置默认
			ValueType:      "string",
			Label:          "时区",
			Order:          100,
			InputType:      "select",
			Validation:     `{"enum":["Asia/Shanghai","Asia/Tokyo","America/New_York","Europe/London","UTC"]}`,
			UIConfig:       datatypes.JSON(`{"options":[{"value":"Asia/Shanghai","label":"中国标准时间 (UTC+8)"},{"value":"Asia/Tokyo","label":"日本标准时间 (UTC+9)"},{"value":"America/New_York","label":"美国东部时间 (UTC-5)"},{"value":"Europe/London","label":"格林威治时间 (UTC+0)"},{"value":"UTC","label":"协调世界时 (UTC)"}]}`),
		},
		{
			Key:            "general.language",
			DefaultValue:   datatypes.JSON(`"zh-CN"`),
			CategoryID:     categoryIDs["general"],
			Group:          "本地化",
			VisibleAt:      "user",
			ConfigurableAt: "user", // 只有用户可配置
			ValueType:      "string",
			Label:          "语言",
			Order:          110,
			InputType:      "select",
			Validation:     `{"enum":["zh-CN","zh-TW","en-US","ja-JP"]}`,
			UIConfig:       datatypes.JSON(`{"options":[{"value":"zh-CN","label":"简体中文"},{"value":"zh-TW","label":"繁體中文"},{"value":"en-US","label":"English (US)"},{"value":"ja-JP","label":"日本語"}]}`),
		},
		// 外观 分组 (Order 200-299)
		{
			Key:            "general.theme",
			DefaultValue:   datatypes.JSON(`"system"`),
			CategoryID:     categoryIDs["general"],
			Group:          "外观",
			VisibleAt:      "user",
			ConfigurableAt: "team", // Team 可为用户设置默认
			ValueType:      "string",
			Label:          "默认主题",
			Order:          200,
			InputType:      "select",
			Validation:     `{"enum":["light","dark","system"]}`,
			UIConfig:       datatypes.JSON(`{"hint":"新用户默认使用的主题","options":[{"value":"light","label":"浅色模式"},{"value":"dark","label":"深色模式"},{"value":"system","label":"跟随系统"}]}`),
		},

		// ==================== Security 安全设置（系统级，管理员专用）====================
		// 密码策略 分组 (Order 1-99)
		{
			Key:            "security.password_min_length",
			DefaultValue:   datatypes.JSON(`8`),
			CategoryID:     categoryIDs["security"],
			Group:          "密码策略",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "number",
			Label:          "密码最小长度",
			Order:          10,
			InputType:      "number",
			Validation:     `{"min":6,"max":32}`,
			UIConfig:       datatypes.JSON(`{"hint":"建议至少 8 位"}`),
		},
		{
			Key:            "security.max_login_attempts",
			DefaultValue:   datatypes.JSON(`5`),
			CategoryID:     categoryIDs["security"],
			Group:          "密码策略",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "number",
			Label:          "最大登录尝试次数",
			Order:          20,
			InputType:      "number",
			Validation:     `{"min":3,"max":10}`,
			UIConfig:       datatypes.JSON(`{"hint":"超过后账户将被临时锁定"}`),
		},
		// 会话管理 分组 (Order 100-199)
		{
			Key:            "security.session_timeout",
			DefaultValue:   datatypes.JSON(`30`),
			CategoryID:     categoryIDs["security"],
			Group:          "会话管理",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "number",
			Label:          "会话超时时间（分钟）",
			Order:          100,
			InputType:      "number",
			Validation:     `{"min":5,"max":1440}`,
			UIConfig:       datatypes.JSON(`{"hint":"用户无操作后自动登出的时间"}`),
		},
		// 高级设置 分组 (Order 200-299)
		{
			Key:            "security.enable_twofa",
			DefaultValue:   datatypes.JSON(`false`),
			CategoryID:     categoryIDs["security"],
			Group:          "高级设置",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "boolean",
			Label:          "强制启用两步验证",
			Order:          200,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"启用后所有用户必须配置两步验证才能登录"}`),
		},

		// ==================== Notification 通知设置 ====================
		// 无分组，通过 Order 控制排序：系统设置 (10) 优先于用户设置 (100+)
		{
			Key:            "notification.enable_notifications",
			DefaultValue:   datatypes.JSON(`true`),
			CategoryID:     categoryIDs["notification"],
			Group:          "",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "boolean",
			Label:          "启用系统通知",
			Order:          10,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"关闭后所有通知渠道将停止发送"}`),
		},
		{
			Key:            "notification.enable_email",
			DefaultValue:   datatypes.JSON(`true`),
			CategoryID:     categoryIDs["notification"],
			Group:          "",
			VisibleAt:      "user",
			ConfigurableAt: "user",
			ValueType:      "boolean",
			Label:          "启用邮件通知",
			Order:          100,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"通过邮件发送系统通知","depends_on":{"key":"notification.enable_notifications","value":true}}`),
		},
		{
			Key:            "notification.enable_sms",
			DefaultValue:   datatypes.JSON(`false`),
			CategoryID:     categoryIDs["notification"],
			Group:          "",
			VisibleAt:      "user",
			ConfigurableAt: "user",
			ValueType:      "boolean",
			Label:          "启用短信通知",
			Order:          110,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"通过短信发送重要通知（需配置短信服务商）","depends_on":{"key":"notification.enable_notifications","value":true}}`),
		},

		// ==================== Backup 备份设置（系统级，管理员专用）====================
		// 基本设置 分组 (Order 1-99)
		{
			Key:            "backup.enable_backup",
			DefaultValue:   datatypes.JSON(`false`),
			CategoryID:     categoryIDs["backup"],
			Group:          "基本设置",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "boolean",
			Label:          "启用自动备份",
			Order:          10,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"开启数据自动备份功能"}`),
		},
		// 备份计划 分组 (Order 100-199)
		{
			Key:            "backup.backup_frequency",
			DefaultValue:   datatypes.JSON(`24`),
			CategoryID:     categoryIDs["backup"],
			Group:          "备份计划",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "number",
			Label:          "备份频率（小时）",
			Order:          100,
			InputType:      "number",
			Validation:     `{"min":1,"max":168}`,
			UIConfig:       datatypes.JSON(`{"hint":"每隔多少小时执行一次备份","depends_on":{"key":"backup.enable_backup","value":true}}`),
		},
		{
			Key:            "backup.retention_days",
			DefaultValue:   datatypes.JSON(`30`),
			CategoryID:     categoryIDs["backup"],
			Group:          "备份计划",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "number",
			Label:          "备份保留天数",
			Order:          110,
			InputType:      "number",
			Validation:     `{"min":7,"max":365}`,
			UIConfig:       datatypes.JSON(`{"hint":"超过保留期的备份将被自动删除","depends_on":{"key":"backup.enable_backup","value":true}}`),
		},

		// ==================== Email 邮件服务配置（系统级，管理员专用）====================
		// 基本设置 分组 (Order 1-99)
		{
			Key:            "email.enabled",
			DefaultValue:   datatypes.JSON(`false`),
			CategoryID:     categoryIDs["email"],
			Group:          "基本设置",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "boolean",
			Label:          "启用邮件服务",
			Order:          10,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"启用后系统可发送邮件通知、验证码等"}`),
		},
		// SMTP 服务器 分组 (Order 100-199)
		{
			Key:            "email.smtp_host",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["email"],
			Group:          "SMTP 服务器",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "SMTP 服务器",
			Order:          100,
			InputType:      "text",
			UIConfig:       datatypes.JSON(`{"hint":"如 smtp.gmail.com 或 smtp.qq.com","depends_on":{"key":"email.enabled","value":true}}`),
		},
		{
			Key:            "email.smtp_port",
			DefaultValue:   datatypes.JSON(`587`),
			CategoryID:     categoryIDs["email"],
			Group:          "SMTP 服务器",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "number",
			Label:          "SMTP 端口",
			Order:          110,
			InputType:      "number",
			Validation:     `{"min":1,"max":65535}`,
			UIConfig:       datatypes.JSON(`{"hint":"常用端口：25(无加密)、465(SSL)、587(TLS)","depends_on":{"key":"email.enabled","value":true}}`),
		},
		{
			Key:            "email.smtp_encryption",
			DefaultValue:   datatypes.JSON(`"tls"`),
			CategoryID:     categoryIDs["email"],
			Group:          "SMTP 服务器",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "加密方式",
			Order:          120,
			InputType:      "select",
			Validation:     `{"enum":["none","ssl","tls"]}`,
			UIConfig:       datatypes.JSON(`{"options":[{"value":"none","label":"无加密"},{"value":"ssl","label":"SSL/TLS"},{"value":"tls","label":"STARTTLS"}],"depends_on":{"key":"email.enabled","value":true}}`),
		},
		{
			Key:            "email.smtp_username",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["email"],
			Group:          "SMTP 服务器",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "SMTP 用户名",
			Order:          130,
			InputType:      "text",
			UIConfig:       datatypes.JSON(`{"hint":"通常为邮箱地址","depends_on":{"key":"email.enabled","value":true}}`),
		},
		{
			Key:            "email.smtp_password",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["email"],
			Group:          "SMTP 服务器",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "SMTP 密码",
			Order:          140,
			InputType:      "password",
			UIConfig:       datatypes.JSON(`{"hint":"部分服务商需使用应用专用密码","depends_on":{"key":"email.enabled","value":true}}`),
		},
		// 发件人 分组 (Order 200-299)
		{
			Key:            "email.from_address",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["email"],
			Group:          "发件人",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "发件人地址",
			Order:          200,
			InputType:      "email",
			UIConfig:       datatypes.JSON(`{"hint":"系统发送邮件时使用的邮箱地址","depends_on":{"key":"email.enabled","value":true}}`),
		},
		{
			Key:            "email.from_name",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["email"],
			Group:          "发件人",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "发件人名称",
			Order:          210,
			InputType:      "text",
			UIConfig:       datatypes.JSON(`{"hint":"显示在收件人邮箱中的发件人名称","depends_on":{"key":"email.enabled","value":true}}`),
		},

		// ==================== OAuth 第三方登录配置（系统级，管理员专用）====================
		// GitHub 分组 (Order 1-99)
		{
			Key:            "oauth.github_enabled",
			DefaultValue:   datatypes.JSON(`false`),
			CategoryID:     categoryIDs["oauth"],
			Group:          "GitHub",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "boolean",
			Label:          "启用 GitHub 登录",
			Order:          10,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"允许用户使用 GitHub 账号登录"}`),
		},
		{
			Key:            "oauth.github_client_id",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["oauth"],
			Group:          "GitHub",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "GitHub Client ID",
			Order:          20,
			InputType:      "text",
			UIConfig:       datatypes.JSON(`{"hint":"在 GitHub Developer Settings 中创建 OAuth App 获取","depends_on":{"key":"oauth.github_enabled","value":true}}`),
		},
		{
			Key:            "oauth.github_client_secret",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["oauth"],
			Group:          "GitHub",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "GitHub Client Secret",
			Order:          30,
			InputType:      "password",
			UIConfig:       datatypes.JSON(`{"hint":"请妥善保管，不要泄露","depends_on":{"key":"oauth.github_enabled","value":true}}`),
		},
		// Google 分组 (Order 100-199)
		{
			Key:            "oauth.google_enabled",
			DefaultValue:   datatypes.JSON(`false`),
			CategoryID:     categoryIDs["oauth"],
			Group:          "Google",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "boolean",
			Label:          "启用 Google 登录",
			Order:          100,
			InputType:      "switch",
			UIConfig:       datatypes.JSON(`{"hint":"允许用户使用 Google 账号登录"}`),
		},
		{
			Key:            "oauth.google_client_id",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["oauth"],
			Group:          "Google",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "Google Client ID",
			Order:          110,
			InputType:      "text",
			UIConfig:       datatypes.JSON(`{"hint":"在 Google Cloud Console 中创建 OAuth 凭据获取","depends_on":{"key":"oauth.google_enabled","value":true}}`),
		},
		{
			Key:            "oauth.google_client_secret",
			DefaultValue:   datatypes.JSON(`""`),
			CategoryID:     categoryIDs["oauth"],
			Group:          "Google",
			VisibleAt:      "system",
			ConfigurableAt: "system",
			ValueType:      "string",
			Label:          "Google Client Secret",
			Order:          120,
			InputType:      "password",
			UIConfig:       datatypes.JSON(`{"hint":"请妥善保管，不要泄露","depends_on":{"key":"oauth.google_enabled","value":true}}`),
		},
	}
}
