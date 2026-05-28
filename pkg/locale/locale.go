package locale

import (
	"context"
	"strings"
	"sync"
)

type Locale string

const (
	ZhCN Locale = "zh-CN"
	EnUS Locale = "en-US"
)

type Messages struct {
	messages map[Locale]map[string]string
	mu       sync.RWMutex
}

var globalMessages = &Messages{
	messages: make(map[Locale]map[string]string),
}

func init() {
	globalMessages.messages[ZhCN] = zhCNMessages
	globalMessages.messages[EnUS] = enUSMessages
}

// T translates a key to the specified locale
func T(locale Locale, key string) string {
	globalMessages.mu.RLock()
	defer globalMessages.mu.RUnlock()

	if msgs, ok := globalMessages.messages[locale]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// fallback to zh-CN
	if msg, ok := globalMessages.messages[ZhCN][key]; ok {
		return msg
	}
	return key
}

// Tf translates a key with format parameters
func Tf(locale Locale, key string, args ...interface{}) string {
	format := T(locale, key)
	if len(args) == 0 {
		return format
	}
	return sprintf(format, args...)
}

func sprintf(format string, args ...interface{}) string {
	// Simple sprintf implementation to avoid fmt dependency
	result := format
	for i, arg := range args {
		placeholder := "{" + itoa(i) + "}"
		result = strings.Replace(result, placeholder, toString(arg), 1)
	}
	return result
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return itoa(val)
	case int64:
		return itoa(int(val))
	case float64:
		return itoa(int(val))
	default:
		return ""
	}
}

// ParseAcceptLanguage parses Accept-Language header and returns best matching locale
func ParseAcceptLanguage(header string) Locale {
	if header == "" {
		return ZhCN
	}

	// Split by comma and parse quality values
	langs := strings.Split(header, ",")
	for _, lang := range langs {
		// Remove quality value (e.g., ";q=0.9")
		code := strings.TrimSpace(strings.Split(lang, ";")[0])
		code = strings.ToLower(code)

		if strings.HasPrefix(code, "zh") {
			return ZhCN
		}
		if strings.HasPrefix(code, "en") {
			return EnUS
		}
	}

	return ZhCN
}

// Context key for locale
type contextKey struct{}

// WithLocale adds locale to context
func WithLocale(ctx context.Context, locale Locale) context.Context {
	return context.WithValue(ctx, contextKey{}, locale)
}

// FromContext extracts locale from context
func FromContext(ctx context.Context) Locale {
	if locale, ok := ctx.Value(contextKey{}).(Locale); ok {
		return locale
	}
	return ZhCN
}
