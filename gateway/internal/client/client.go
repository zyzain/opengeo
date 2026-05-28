package client

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// validateWebhookURL 校验webhook URL，防止SSRF攻击
func validateWebhookURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("only http/https URLs are allowed")
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("missing hostname")
	}

	// 禁止内网地址
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("internal/private network addresses are not allowed")
		}
	}

	// 禁止常见内网域名
	lowerHost := strings.ToLower(host)
	blockedHosts := []string{"localhost", "127.0.0.1", "0.0.0.0", "169.254.169.254", "metadata.google.internal"}
	for _, blocked := range blockedHosts {
		if lowerHost == blocked {
			return fmt.Errorf("internal addresses are not allowed")
		}
	}

	return nil
}
