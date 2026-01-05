package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Setting.Validate 测试
// =============================================================================

func TestSetting_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setting *Setting
		wantErr error
	}{
		{
			name: "valid setting - system level",
			setting: &Setting{
				Key:            "general.site_name",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "My Site",
			},
			wantErr: nil,
		},
		{
			name: "valid setting - org visible, org configurable",
			setting: &Setting{
				Key:            "org.theme",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelOrg),
				ConfigurableAt: string(ScopeLevelOrg),
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "dark",
			},
			wantErr: nil,
		},
		{
			name: "valid setting - user visible, team configurable",
			setting: &Setting{
				Key:            "user.notification",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelUser),
				ConfigurableAt: string(ScopeLevelTeam),
				ValueType:      ValueTypeBoolean,
				InputType:      InputTypeSwitch,
				DefaultValue:   true,
			},
			wantErr: nil,
		},
		{
			name: "empty key",
			setting: &Setting{
				Key:            "",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "test",
			},
			wantErr: ErrInvalidValue,
		},
		{
			name: "invalid key format - no dot",
			setting: &Setting{
				Key:            "invalidkey",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "test",
			},
			wantErr: ErrInvalidKeyFormat,
		},
		{
			name: "zero category id",
			setting: &Setting{
				Key:            "general.site_name",
				CategoryID:     0,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "test",
			},
			wantErr: ErrCategoryNotFound,
		},
		{
			name: "invalid visible at",
			setting: &Setting{
				Key:            "general.site_name",
				CategoryID:     1,
				VisibleAt:      "invalid",
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "test",
			},
			wantErr: ErrInvalidVisibleAt,
		},
		{
			name: "invalid configurable at",
			setting: &Setting{
				Key:            "general.site_name",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: "invalid",
				ValueType:      ValueTypeString,
				InputType:      InputTypeText,
				DefaultValue:   "test",
			},
			wantErr: ErrInvalidConfigurableAt,
		},
		{
			name: "invalid value type",
			setting: &Setting{
				Key:            "general.site_name",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      "invalid",
				InputType:      InputTypeText,
				DefaultValue:   "test",
			},
			wantErr: ErrInvalidValueType,
		},
		{
			name: "invalid input type",
			setting: &Setting{
				Key:            "general.site_name",
				CategoryID:     1,
				VisibleAt:      string(ScopeLevelSystem),
				ConfigurableAt: string(ScopeLevelSystem),
				ValueType:      ValueTypeString,
				InputType:      "invalid_input",
				DefaultValue:   "test",
			},
			wantErr: ErrInvalidInputType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setting.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// Setting.ValidateValue 测试
// =============================================================================

func TestSetting_ValidateValue(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     any
		wantErr   bool
	}{
		// nil 值
		{"nil value is always valid", ValueTypeString, nil, false},

		// string 类型
		{"string - valid", ValueTypeString, "hello", false},
		{"string - invalid number", ValueTypeString, 123, true},
		{"string - invalid bool", ValueTypeString, true, true},

		// number 类型
		{"number - valid int", ValueTypeNumber, 123, false},
		{"number - valid int64", ValueTypeNumber, int64(123), false},
		{"number - valid float64", ValueTypeNumber, 123.45, false},
		{"number - invalid string", ValueTypeNumber, "123", true},
		{"number - invalid bool", ValueTypeNumber, true, true},

		// boolean 类型
		{"boolean - valid true", ValueTypeBoolean, true, false},
		{"boolean - valid false", ValueTypeBoolean, false, false},
		{"boolean - invalid string", ValueTypeBoolean, "true", true},
		{"boolean - invalid number", ValueTypeBoolean, 1, true},

		// json 类型
		{"json - valid map", ValueTypeJSON, map[string]any{"key": "value"}, false},
		{"json - valid slice", ValueTypeJSON, []any{1, 2, 3}, false},
		{"json - invalid string", ValueTypeJSON, `{"key": "value"}`, true},
		{"json - invalid number", ValueTypeJSON, 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ValueType: tt.valueType}
			err := s.ValidateValue(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// Setting.CoerceValue 测试
// =============================================================================

func TestSetting_CoerceValue(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		input     any
		want      any
		wantErr   bool
	}{
		// nil 值
		{"nil returns nil", ValueTypeString, nil, nil, false},

		// 已匹配的类型直接返回
		{"string already matched", ValueTypeString, "hello", "hello", false},
		{"number already matched", ValueTypeNumber, 123.0, 123.0, false},
		{"boolean already matched", ValueTypeBoolean, true, true, false},

		// string 转换
		{"int to string", ValueTypeString, 123, "123", false},
		{"float to string", ValueTypeString, 123.45, "123.45", false},
		{"bool to string", ValueTypeString, true, "true", false},

		// number 转换
		{"string to number", ValueTypeNumber, "123.45", 123.45, false},
		{"int already valid for number", ValueTypeNumber, 100, 100, false}, // int 已是有效 number 类型，直接返回

		// boolean 转换
		{"string true to bool", ValueTypeBoolean, "true", true, false},
		{"string false to bool", ValueTypeBoolean, "false", false, false},
		{"int 1 to bool", ValueTypeBoolean, 1, true, false},
		{"int 0 to bool", ValueTypeBoolean, 0, false, false},

		// json 转换
		{"string to json map", ValueTypeJSON, `{"key": "value"}`, map[string]any{"key": "value"}, false},
		{"string to json array", ValueTypeJSON, `[1, 2, 3]`, []any{float64(1), float64(2), float64(3)}, false},
		{"invalid json string", ValueTypeJSON, "not json", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ValueType: tt.valueType}
			got, err := s.CoerceValue(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// =============================================================================
// Setting.IsVisibleAtScope 测试
// =============================================================================

func TestSetting_IsVisibleAtScope(t *testing.T) {
	// 新语义（权限等级模式）：
	// - VisibleAt 表示"最低可见权限"
	// - queryScope 是查询者的权限级别
	// - 可见条件：queryScope <= VisibleAt（权限足够高才能看到）
	// 层级顺序（权限从高到低）：system(0) < org(1) < team(2) < user(3) < public(4)
	//
	// 例如：
	//   - VisibleAt="system"(0) → 只有 system 可见
	//   - VisibleAt="user"(3)   → system/org/team/user 都可见
	//   - VisibleAt="public"(4) → 所有人可见（包括未登录）
	tests := []struct {
		name        string
		visibleAt   string
		queryScope  ScopeLevel
		wantVisible bool
	}{
		// system 级别设置：只有管理员可见
		{
			name:        "system setting visible to system",
			visibleAt:   string(ScopeLevelSystem),
			queryScope:  ScopeLevelSystem,
			wantVisible: true, // system(0) <= system(0) ✓
		},
		{
			name:        "system setting NOT visible to org",
			visibleAt:   string(ScopeLevelSystem),
			queryScope:  ScopeLevelOrg,
			wantVisible: false, // org(1) > system(0) ✗
		},
		{
			name:        "system setting NOT visible to user",
			visibleAt:   string(ScopeLevelSystem),
			queryScope:  ScopeLevelUser,
			wantVisible: false, // user(3) > system(0) ✗
		},

		// org 级别设置：管理员和组织管理者可见
		{
			name:        "org setting visible to system",
			visibleAt:   string(ScopeLevelOrg),
			queryScope:  ScopeLevelSystem,
			wantVisible: true, // system(0) <= org(1) ✓
		},
		{
			name:        "org setting visible to org",
			visibleAt:   string(ScopeLevelOrg),
			queryScope:  ScopeLevelOrg,
			wantVisible: true, // org(1) <= org(1) ✓
		},
		{
			name:        "org setting NOT visible to team",
			visibleAt:   string(ScopeLevelOrg),
			queryScope:  ScopeLevelTeam,
			wantVisible: false, // team(2) > org(1) ✗
		},

		// team 级别设置：团队成员及以上可见
		{
			name:        "team setting visible to system",
			visibleAt:   string(ScopeLevelTeam),
			queryScope:  ScopeLevelSystem,
			wantVisible: true, // system(0) <= team(2) ✓
		},
		{
			name:        "team setting visible to team",
			visibleAt:   string(ScopeLevelTeam),
			queryScope:  ScopeLevelTeam,
			wantVisible: true, // team(2) <= team(2) ✓
		},
		{
			name:        "team setting NOT visible to user",
			visibleAt:   string(ScopeLevelTeam),
			queryScope:  ScopeLevelUser,
			wantVisible: false, // user(3) > team(2) ✗
		},

		// user 级别设置：所有登录用户可见
		{
			name:        "user setting visible to system",
			visibleAt:   string(ScopeLevelUser),
			queryScope:  ScopeLevelSystem,
			wantVisible: true, // system(0) <= user(3) ✓
		},
		{
			name:        "user setting visible to user",
			visibleAt:   string(ScopeLevelUser),
			queryScope:  ScopeLevelUser,
			wantVisible: true, // user(3) <= user(3) ✓
		},

		// public 级别设置：任何人可见
		{
			name:        "public setting visible to system",
			visibleAt:   string(ScopeLevelPublic),
			queryScope:  ScopeLevelSystem,
			wantVisible: true, // system(0) <= public(4) ✓
		},
		{
			name:        "public setting visible to user",
			visibleAt:   string(ScopeLevelPublic),
			queryScope:  ScopeLevelUser,
			wantVisible: true, // user(3) <= public(4) ✓
		},
		{
			name:        "public setting visible to public",
			visibleAt:   string(ScopeLevelPublic),
			queryScope:  ScopeLevelPublic,
			wantVisible: true, // public(4) <= public(4) ✓
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt}
			assert.Equal(t, tt.wantVisible, s.IsVisibleAtScope(tt.queryScope))
		})
	}
}

// =============================================================================
// Setting.IsConfigurableAtScope 测试
// =============================================================================

func TestSetting_IsConfigurableAtScope(t *testing.T) {
	tests := []struct {
		name             string
		configurableAt   string
		queryScope       ScopeLevel
		wantConfigurable bool
	}{
		{
			name:             "system configurable only at system",
			configurableAt:   string(ScopeLevelSystem),
			queryScope:       ScopeLevelSystem,
			wantConfigurable: true, // system(0) <= system(0)
		},
		{
			name:             "system NOT configurable at org",
			configurableAt:   string(ScopeLevelSystem),
			queryScope:       ScopeLevelOrg,
			wantConfigurable: false, // org(1) > system(0)
		},
		{
			name:             "org configurable at system",
			configurableAt:   string(ScopeLevelOrg),
			queryScope:       ScopeLevelSystem,
			wantConfigurable: true, // system(0) <= org(1)
		},
		{
			name:             "org configurable at org",
			configurableAt:   string(ScopeLevelOrg),
			queryScope:       ScopeLevelOrg,
			wantConfigurable: true, // org(1) <= org(1)
		},
		{
			name:             "org NOT configurable at team",
			configurableAt:   string(ScopeLevelOrg),
			queryScope:       ScopeLevelTeam,
			wantConfigurable: false, // team(2) > org(1)
		},
		{
			name:             "team configurable at org",
			configurableAt:   string(ScopeLevelTeam),
			queryScope:       ScopeLevelOrg,
			wantConfigurable: true, // org(1) <= team(2)
		},
		{
			name:             "team configurable at team",
			configurableAt:   string(ScopeLevelTeam),
			queryScope:       ScopeLevelTeam,
			wantConfigurable: true, // team(2) <= team(2)
		},
		{
			name:             "team NOT configurable at user",
			configurableAt:   string(ScopeLevelTeam),
			queryScope:       ScopeLevelUser,
			wantConfigurable: false, // user(3) > team(2)
		},
		{
			name:             "user configurable at all scopes",
			configurableAt:   string(ScopeLevelUser),
			queryScope:       ScopeLevelSystem,
			wantConfigurable: true, // system(0) <= user(3)
		},
		{
			name:             "user configurable at user",
			configurableAt:   string(ScopeLevelUser),
			queryScope:       ScopeLevelUser,
			wantConfigurable: true, // user(3) <= user(3)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ConfigurableAt: tt.configurableAt}
			assert.Equal(t, tt.wantConfigurable, s.IsConfigurableAtScope(tt.queryScope))
		})
	}
}

// =============================================================================
// Setting 层级检测方法测试
// =============================================================================

func TestSetting_IsSystemLevel(t *testing.T) {
	tests := []struct {
		name      string
		visibleAt string
		want      bool
	}{
		{"system level", string(ScopeLevelSystem), true},
		{"org level", string(ScopeLevelOrg), false},
		{"team level", string(ScopeLevelTeam), false},
		{"user level", string(ScopeLevelUser), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt}
			assert.Equal(t, tt.want, s.IsSystemLevel())
		})
	}
}

func TestSetting_IsOrgLevel(t *testing.T) {
	tests := []struct {
		name      string
		visibleAt string
		want      bool
	}{
		{"system level", string(ScopeLevelSystem), false},
		{"org level", string(ScopeLevelOrg), true},
		{"team level", string(ScopeLevelTeam), false},
		{"user level", string(ScopeLevelUser), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt}
			assert.Equal(t, tt.want, s.IsOrgLevel())
		})
	}
}

func TestSetting_IsTeamLevel(t *testing.T) {
	tests := []struct {
		name      string
		visibleAt string
		want      bool
	}{
		{"system level", string(ScopeLevelSystem), false},
		{"org level", string(ScopeLevelOrg), false},
		{"team level", string(ScopeLevelTeam), true},
		{"user level", string(ScopeLevelUser), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt}
			assert.Equal(t, tt.want, s.IsTeamLevel())
		})
	}
}

func TestSetting_IsUserLevel(t *testing.T) {
	tests := []struct {
		name      string
		visibleAt string
		want      bool
	}{
		{"system level", string(ScopeLevelSystem), false},
		{"org level", string(ScopeLevelOrg), false},
		{"team level", string(ScopeLevelTeam), false},
		{"user level", string(ScopeLevelUser), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt}
			assert.Equal(t, tt.want, s.IsUserLevel())
		})
	}
}

// =============================================================================
// Setting 可配置性检测方法测试
// =============================================================================

func TestSetting_CanOrgConfigure(t *testing.T) {
	tests := []struct {
		name           string
		configurableAt string
		want           bool
	}{
		// ConfigurableAt 定义最大可配置级别（到此级别为止可配置）
		// system(0) 只有 system 可配置
		{"system configurable", string(ScopeLevelSystem), false},
		// org(1) system 和 org 可配置
		{"org configurable", string(ScopeLevelOrg), true},
		// team(2) system、org、team 可配置
		{"team configurable", string(ScopeLevelTeam), true},
		// user(3) 所有级别都可配置
		{"user configurable", string(ScopeLevelUser), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ConfigurableAt: tt.configurableAt}
			assert.Equal(t, tt.want, s.CanOrgConfigure())
		})
	}
}

func TestSetting_CanTeamConfigure(t *testing.T) {
	tests := []struct {
		name           string
		configurableAt string
		want           bool
	}{
		// system(0) 只有 system 可配置
		{"system configurable", string(ScopeLevelSystem), false},
		// org(1) system 和 org 可配置，team 不可配置
		{"org configurable", string(ScopeLevelOrg), false},
		// team(2) system、org、team 可配置
		{"team configurable", string(ScopeLevelTeam), true},
		// user(3) 所有级别都可配置
		{"user configurable", string(ScopeLevelUser), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ConfigurableAt: tt.configurableAt}
			assert.Equal(t, tt.want, s.CanTeamConfigure())
		})
	}
}

func TestSetting_CanUserConfigure(t *testing.T) {
	tests := []struct {
		name           string
		configurableAt string
		want           bool
	}{
		// system(0) 只有 system 可配置
		{"system configurable", string(ScopeLevelSystem), false},
		// org(1) system 和 org 可配置
		{"org configurable", string(ScopeLevelOrg), false},
		// team(2) system、org、team 可配置
		{"team configurable", string(ScopeLevelTeam), false},
		// user(3) 所有级别都可配置
		{"user configurable", string(ScopeLevelUser), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{ConfigurableAt: tt.configurableAt}
			assert.Equal(t, tt.want, s.CanUserConfigure())
		})
	}
}

// =============================================================================
// Setting 组合场景方法测试
// =============================================================================

func TestSetting_IsOrgOnly(t *testing.T) {
	tests := []struct {
		name           string
		visibleAt      string
		configurableAt string
		want           bool
	}{
		{"org visible, org configurable", string(ScopeLevelOrg), string(ScopeLevelOrg), true},
		{"org visible, team configurable", string(ScopeLevelOrg), string(ScopeLevelTeam), false},
		{"system visible, system configurable", string(ScopeLevelSystem), string(ScopeLevelSystem), false},
		{"user visible, user configurable", string(ScopeLevelUser), string(ScopeLevelUser), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt, ConfigurableAt: tt.configurableAt}
			assert.Equal(t, tt.want, s.IsOrgOnly())
		})
	}
}

func TestSetting_IsTeamDefaultForUser(t *testing.T) {
	tests := []struct {
		name           string
		visibleAt      string
		configurableAt string
		want           bool
	}{
		{"user visible, team configurable", string(ScopeLevelUser), string(ScopeLevelTeam), true},
		{"user visible, user configurable", string(ScopeLevelUser), string(ScopeLevelUser), false},
		{"team visible, team configurable", string(ScopeLevelTeam), string(ScopeLevelTeam), false},
		{"user visible, org configurable", string(ScopeLevelUser), string(ScopeLevelOrg), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt, ConfigurableAt: tt.configurableAt}
			assert.Equal(t, tt.want, s.IsTeamDefaultForUser())
		})
	}
}

func TestSetting_IsVisibleToUser(t *testing.T) {
	// 新语义：只有 VisibleAt >= user(3) 的设置才对普通用户可见
	// 即：user 和 public 级别
	tests := []struct {
		name      string
		visibleAt string
		want      bool
	}{
		{"system NOT visible to user", string(ScopeLevelSystem), false}, // user(3) > system(0)
		{"org NOT visible to user", string(ScopeLevelOrg), false},       // user(3) > org(1)
		{"team NOT visible to user", string(ScopeLevelTeam), false},     // user(3) > team(2)
		{"user visible to user", string(ScopeLevelUser), true},          // user(3) <= user(3) ✓
		{"public visible to user", string(ScopeLevelPublic), true},      // user(3) <= public(4) ✓
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.visibleAt}
			assert.Equal(t, tt.want, s.IsVisibleToUser())
		})
	}
}

// =============================================================================
// 向后兼容方法测试
// =============================================================================

func TestSetting_IsSystemScope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"system scope", ScopeSystem, true},
		{"user scope", ScopeUser, false},
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.scope}
			assert.Equal(t, tt.want, s.IsSystemScope())
		})
	}
}

func TestSetting_IsUserScope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"user scope", ScopeUser, true},
		{"system scope", ScopeSystem, false},
		{"org scope", ScopeOrg, false},
		{"team scope", ScopeTeam, false},
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.scope}
			assert.Equal(t, tt.want, s.IsUserScope())
		})
	}
}

func TestSetting_IsOrgScope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"org scope", ScopeOrg, true},
		{"system scope", ScopeSystem, false},
		{"user scope", ScopeUser, false},
		{"team scope", ScopeTeam, false},
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.scope}
			assert.Equal(t, tt.want, s.IsOrgScope())
		})
	}
}

func TestSetting_IsTeamScope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"team scope", ScopeTeam, true},
		{"system scope", ScopeSystem, false},
		{"org scope", ScopeOrg, false},
		{"user scope", ScopeUser, false},
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{VisibleAt: tt.scope}
			assert.Equal(t, tt.want, s.IsTeamScope())
		})
	}
}

// =============================================================================
// SettingCategory.Validate 测试
// =============================================================================

func TestSettingCategory_Validate(t *testing.T) {
	tests := []struct {
		name     string
		category *SettingCategory
		wantErr  bool
	}{
		{
			name: "valid category",
			category: &SettingCategory{
				Key:   "general",
				Label: "常规设置",
				Icon:  "mdi-cog",
			},
			wantErr: false,
		},
		{
			name: "empty key",
			category: &SettingCategory{
				Key:   "",
				Label: "常规设置",
				Icon:  "mdi-cog",
			},
			wantErr: true,
		},
		{
			name: "empty label",
			category: &SettingCategory{
				Key:   "general",
				Label: "",
				Icon:  "mdi-cog",
			},
			wantErr: true,
		},
		{
			name: "empty icon",
			category: &SettingCategory{
				Key:   "general",
				Label: "常规设置",
				Icon:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSettingCategory_IsValidKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"general is valid", CategoryGeneral, true},
		{"security is valid", CategorySecurity, true},
		{"notification is valid", CategoryNotification, true},
		{"backup is valid", CategoryBackup, true},
		{"unknown is invalid", "unknown", false},
		{"empty is invalid", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SettingCategory{Key: tt.key}
			assert.Equal(t, tt.want, c.IsValidKey())
		})
	}
}
