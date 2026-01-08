package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// SettingKey 测试
// =============================================================================

func TestNewSettingKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantErr  bool
		wantCat  string
		wantName string
	}{
		{
			name:     "valid key",
			key:      "general.site_name",
			wantErr:  false,
			wantCat:  "general",
			wantName: "site_name",
		},
		{
			name:     "valid key with multiple dots",
			key:      "general.nested.key",
			wantErr:  false,
			wantCat:  "general",
			wantName: "nested.key",
		},
		{
			name:    "empty key",
			key:     "",
			wantErr: true,
		},
		{
			name:    "no dot",
			key:     "invalidkey",
			wantErr: true,
		},
		{
			name:    "dot at start",
			key:     ".key",
			wantErr: true,
		},
		{
			name:    "dot at end",
			key:     "category.",
			wantErr: true,
		},
		{
			name:    "only dot",
			key:     ".",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk, err := NewSettingKey(tt.key)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidKeyFormat)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.key, sk.String())
				assert.Equal(t, tt.wantCat, sk.Category())
				assert.Equal(t, tt.wantName, sk.Name())
			}
		})
	}
}

func TestMustSettingKey(t *testing.T) {
	t.Run("valid key does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			sk := MustSettingKey("general.site_name")
			assert.Equal(t, "general.site_name", sk.String())
		})
	})

	t.Run("invalid key panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustSettingKey("invalid")
		})
	})
}

func TestSettingKey_String(t *testing.T) {
	sk, err := NewSettingKey("general.theme")
	require.NoError(t, err)
	assert.Equal(t, "general.theme", sk.String())
}

func TestSettingKey_Category(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"simple key", "general.theme", "general"},
		{"nested key", "security.password.min_length", "security"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk, err := NewSettingKey(tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.want, sk.Category())
		})
	}
}

func TestSettingKey_Name(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"simple key", "general.theme", "theme"},
		{"nested key", "security.password.min_length", "password.min_length"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk, err := NewSettingKey(tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.want, sk.Name())
		})
	}
}

func TestSettingKey_IsEmpty(t *testing.T) {
	t.Run("valid key is not empty", func(t *testing.T) {
		sk, err := NewSettingKey("general.theme")
		require.NoError(t, err)
		assert.False(t, sk.IsEmpty())
	})

	t.Run("zero value is empty", func(t *testing.T) {
		var sk SettingKey
		assert.True(t, sk.IsEmpty())
	})
}

func TestSettingKey_Equal(t *testing.T) {
	sk1, _ := NewSettingKey("general.theme")
	sk2, _ := NewSettingKey("general.theme")
	sk3, _ := NewSettingKey("general.other")

	t.Run("equal keys", func(t *testing.T) {
		assert.True(t, sk1.Equal(sk2))
	})

	t.Run("not equal keys", func(t *testing.T) {
		assert.False(t, sk1.Equal(sk3))
	})

	t.Run("equal with zero value", func(t *testing.T) {
		var zero SettingKey
		assert.False(t, sk1.Equal(zero))
	})
}

// =============================================================================
// Category 测试
// =============================================================================

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name    string
		cat     string
		wantErr bool
	}{
		{"general is valid", "general", false},
		{"security is valid", "security", false},
		{"notification is valid", "notification", false},
		{"backup is valid", "backup", false},
		{"custom category is valid", "custom", false},
		{"unknown is valid", "unknown", false},
		{"empty is invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCategory(tt.cat)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidCategory)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.cat, c.String())
			}
		})
	}
}

func TestMustCategory(t *testing.T) {
	t.Run("valid category does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			c := MustCategory("general")
			assert.Equal(t, "general", c.String())
		})
	})

	t.Run("empty category panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustCategory("")
		})
	})
}

func TestCategory_String(t *testing.T) {
	c, err := NewCategory("general")
	require.NoError(t, err)
	assert.Equal(t, "general", c.String())
}

func TestCategory_IsValid(t *testing.T) {
	tests := []struct {
		name string
		cat  string
		want bool
	}{
		{"general is valid", "general", true},
		{"security is valid", "security", true},
		{"notification is valid", "notification", true},
		{"backup is valid", "backup", true},
		{"custom category is valid", "custom", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCategory(tt.cat)
			require.NoError(t, err)
			assert.Equal(t, tt.want, c.IsValid())
		})
	}

	t.Run("zero value is invalid", func(t *testing.T) {
		var c Category
		assert.False(t, c.IsValid())
	})
}

func TestCategory_IsEmpty(t *testing.T) {
	t.Run("valid category is not empty", func(t *testing.T) {
		c, err := NewCategory("general")
		require.NoError(t, err)
		assert.False(t, c.IsEmpty())
	})

	t.Run("zero value is empty", func(t *testing.T) {
		var c Category
		assert.True(t, c.IsEmpty())
	})
}

func TestCategory_Equal(t *testing.T) {
	c1, _ := NewCategory("general")
	c2, _ := NewCategory("general")
	c3, _ := NewCategory("security")

	t.Run("equal categories", func(t *testing.T) {
		assert.True(t, c1.Equal(c2))
	})

	t.Run("not equal categories", func(t *testing.T) {
		assert.False(t, c1.Equal(c3))
	})

	t.Run("equal with zero value", func(t *testing.T) {
		var zero Category
		assert.False(t, c1.Equal(zero))
	})
}
