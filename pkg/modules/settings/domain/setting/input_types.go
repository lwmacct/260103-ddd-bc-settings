package setting

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

// InputType 补充常量（未在 constants.go 中定义的）。
const (
	InputTypeSwitch   = "switch"   // 开关，由 ValueType 校验
	InputTypeRadio    = "radio"    // 单选按钮组，校验值在 options 列表中
	InputTypeCheckbox = "checkbox" // 多选复选框
	InputTypeJSON     = "json"     // JSON 编辑器，由 ValueType 校验（别名为 InputTypeJSONEditor）
)

// validInputTypes 有效的 InputType 集合。
var validInputTypes = map[string]bool{
	InputTypeText:     true,
	InputTypeTextarea: true,
	InputTypeNumber:   true,
	InputTypeSwitch:   true,
	InputTypeSelect:   true,
	InputTypeRadio:    true,
	InputTypeCheckbox: true,
	InputTypePassword: true,
	InputTypeEmail:    true,
	InputTypeURL:      true,
	InputTypeJSON:     true,
	InputTypeColor:    true,
}

// IsValidInputType 报告给定的 InputType 是否有效。
func IsValidInputType(t string) bool {
	if t == "" {
		return true // 空值默认为 text
	}
	return validInputTypes[t]
}

// ValidateByInputType 基于控件类型验证值的格式。
//
// 校验规则：
//   - email: RFC 5322 简化版邮箱格式
//   - url: http/https 开头的有效 URL
//   - password: 最小长度 6 个字符
//   - select/radio: 值必须在 options 列表中（需外部校验）
//
// 其他类型由 [Setting.ValidateValue] 按 ValueType 校验。
func (s *Setting) ValidateByInputType(value any) error {
	if value == nil {
		return nil // nil 值总是允许的
	}

	// 获取实际的 InputType（空值默认为 text）
	inputType := s.InputType
	if inputType == "" {
		inputType = InputTypeText
	}

	// 非字符串值由 ValidateValue 处理
	str, ok := value.(string)
	if !ok {
		return nil
	}

	switch inputType {
	case InputTypeEmail:
		if !isValidEmail(str) {
			return fmt.Errorf("%w: invalid email format", ErrInvalidValue)
		}
	case InputTypeURL:
		if !isValidURL(str) {
			return fmt.Errorf("%w: invalid URL format", ErrInvalidValue)
		}
	case InputTypePassword:
		if len(str) < 6 {
			return fmt.Errorf("%w: password too short (min 6 characters)", ErrInvalidValue)
		}
		// text, textarea, number, switch, json, color 等无内置格式校验
		// select, radio 的 options 校验需要外部处理（需要解析 UIConfig）
	}

	return nil
}

// IsValidInputType 报告实体的 InputType 是否有效。
func (s *Setting) IsValidInputType() bool {
	return IsValidInputType(s.InputType)
}

// isValidEmail 使用 RFC 5322 标准验证邮箱格式。
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

// isValidURL 验证 URL 格式（仅接受 http/https）。
func isValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	// 仅接受 http 和 https 协议
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return false
	}
	// 必须有主机名
	return u.Host != ""
}
