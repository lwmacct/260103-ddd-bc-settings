package manualtest

import (
	"os"
	"testing"
)

// SkipIfNotManual 如果 MANUAL 环境变量未设置则跳过测试。
func SkipIfNotManual(t *testing.T) {
	t.Helper()
	if os.Getenv("MANUAL") == "" {
		t.SkipNow()
	}
}
