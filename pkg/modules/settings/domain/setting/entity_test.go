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
			name: "valid setting",
			setting: &Setting{
				Key:          "general.site_name",
				CategoryID:   1,
				Scope:        ScopeSystem,
				ValueType:    ValueTypeString,
				InputType:    InputTypeText,
				DefaultValue: "My Site",
			},
			wantErr: nil,
		},
		{
			name: "empty key",
			setting: &Setting{
				Key:          "",
				CategoryID:   1,
				Scope:        ScopeSystem,
				ValueType:    ValueTypeString,
				InputType:    InputTypeText,
				DefaultValue: "test",
			},
			wantErr: ErrInvalidValue,
		},
		{
			name: "invalid key format - no dot",
			setting: &Setting{
				Key:          "invalidkey",
				CategoryID:   1,
				Scope:        ScopeSystem,
				ValueType:    ValueTypeString,
				InputType:    InputTypeText,
				DefaultValue: "test",
			},
			wantErr: ErrInvalidKeyFormat,
		},
		{
			name: "zero category id",
			setting: &Setting{
				Key:          "general.site_name",
				CategoryID:   0,
				Scope:        ScopeSystem,
				ValueType:    ValueTypeString,
				InputType:    InputTypeText,
				DefaultValue: "test",
			},
			wantErr: ErrCategoryNotFound,
		},
		{
			name: "invalid scope",
			setting: &Setting{
				Key:          "general.site_name",
				CategoryID:   1,
				Scope:        "invalid",
				ValueType:    ValueTypeString,
				InputType:    InputTypeText,
				DefaultValue: "test",
			},
			wantErr: ErrInvalidScope,
		},
		{
			name: "invalid value type",
			setting: &Setting{
				Key:          "general.site_name",
				CategoryID:   1,
				Scope:        ScopeSystem,
				ValueType:    "invalid",
				InputType:    InputTypeText,
				DefaultValue: "test",
			},
			wantErr: ErrInvalidValueType,
		},
		{
			name: "invalid input type",
			setting: &Setting{
				Key:          "general.site_name",
				CategoryID:   1,
				Scope:        ScopeSystem,
				ValueType:    ValueTypeString,
				InputType:    "invalid_input",
				DefaultValue: "test",
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
		{"int already valid for number", ValueTypeNumber, 100, 100, false}, // int 已是有效 number 类型，直接返回,

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
// Setting Scope 方法测试
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
			s := &Setting{Scope: tt.scope}
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
		{"invalid scope", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{Scope: tt.scope}
			assert.Equal(t, tt.want, s.IsUserScope())
		})
	}
}

func TestSetting_IsVisibleToUser(t *testing.T) {
	tests := []struct {
		name   string
		scope  string
		public bool
		want   bool
	}{
		{"user scope is visible", ScopeUser, false, true},
		{"system + public is visible", ScopeSystem, true, true},
		{"system + private is not visible", ScopeSystem, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Setting{Scope: tt.scope, Public: tt.public}
			assert.Equal(t, tt.want, s.IsVisibleToUser())
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
