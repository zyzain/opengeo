package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"opengeo/pkg/config"
	"opengeo/service/account/internal/dal"
)

// HealthCheckService 健康检测服务
type HealthCheckService struct {
	accountRepo *dal.AccountRepository
	alertRepo   *dal.AlertRecordRepository
	httpClient  *http.Client
}

// NewHealthCheckService 创建健康检测服务
func NewHealthCheckService(accountRepo *dal.AccountRepository, alertRepo *dal.AlertRecordRepository) *HealthCheckService {
	cfg := config.GetConfig()
	return &HealthCheckService{
		accountRepo: accountRepo,
		alertRepo:   alertRepo,
		httpClient: &http.Client{
			Timeout: cfg.HealthCheck.HTTPTimeout,
		},
	}
}

// CheckAccountHealth 检测单个账号健康状态
func (s *HealthCheckService) CheckAccountHealth(ctx context.Context, accountID int64) (*dal.HealthCheckResult, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	checks := make([]dal.CheckItem, 0)
	totalScore := float32(0)

	// 检查1: 账号状态
	statusCheck := s.checkAccountStatus(account)
	checks = append(checks, statusCheck)
	totalScore += statusCheck.Score

	// 检查2: 健康分数趋势
	trendCheck := s.checkHealthTrend(ctx, accountID)
	checks = append(checks, trendCheck)
	totalScore += trendCheck.Score

	// 检查3: 最近发布成功率
	publishCheck := s.checkPublishSuccessRate(ctx, accountID)
	checks = append(checks, publishCheck)
	totalScore += publishCheck.Score

	// 检查4: 平台风控检测
	riskCheck := s.checkPlatformRisk(account)
	checks = append(checks, riskCheck)
	totalScore += riskCheck.Score

	// 计算平均分
	avgScore := totalScore / float32(len(checks))
	avgScore = float32(math.Min(math.Max(float64(avgScore), 0), 100))

	// 确定状态和告警级别
	status, alertLevel := s.determineStatus(avgScore, checks)

	// 确定是否需要暂停发布
	publishPaused, pauseReason := s.determinePause(avgScore, checks)

	result := &dal.HealthCheckResult{
		AccountID:     accountID,
		HealthScore:   avgScore,
		Status:        status,
		AlertLevel:    alertLevel,
		PublishPaused: publishPaused,
		PauseReason:   pauseReason,
		Checks:        checks,
	}

	// 保存检测结果
	health := &dal.AccountHealth{
		AccountID:     accountID,
		HealthScore:   avgScore,
		Status:        status,
		CheckType:     "auto",
		CheckDetails:  s.serializeChecks(checks),
		AlertLevel:    alertLevel,
		PublishPaused: publishPaused,
		PauseReason:   pauseReason,
		CheckedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}

	if err := s.accountRepo.SaveHealth(ctx, health); err != nil {
		return nil, fmt.Errorf("save health: %w", err)
	}

	// 更新账号健康分数
	if err := s.accountRepo.UpdateAccountHealthScore(ctx, accountID, avgScore, publishPaused); err != nil {
		return nil, fmt.Errorf("update health score: %w", err)
	}

	return result, nil
}

// CheckAllAccounts 批量检测所有账号
func (s *HealthCheckService) CheckAllAccounts(ctx context.Context, intervalMinutes int, limit int) ([]*dal.HealthCheckResult, error) {
	accounts, err := s.accountRepo.GetAccountsToCheck(ctx, intervalMinutes, limit)
	if err != nil {
		return nil, fmt.Errorf("get accounts to check: %w", err)
	}

	results := make([]*dal.HealthCheckResult, 0, len(accounts))
	for _, account := range accounts {
		result, err := s.CheckAccountHealth(ctx, account.ID)
		if err != nil {
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// GetHealthHistory 获取健康历史
func (s *HealthCheckService) GetHealthHistory(ctx context.Context, accountID int64, limit int) ([]*dal.AccountHealth, error) {
	return s.accountRepo.GetHealthHistory(ctx, accountID, limit)
}

// ==================== 检查项实现 ====================

// checkAccountStatus 检查账号状态
func (s *HealthCheckService) checkAccountStatus(account *dal.Account) dal.CheckItem {
	check := dal.CheckItem{
		Name: "账号状态",
	}

	switch account.Status {
	case 1:
		check.Passed = true
		check.Score = 100
		check.Detail = "账号状态正常"
	case 0:
		check.Passed = false
		check.Score = 0
		check.Detail = "账号已禁用"
	case 2:
		check.Passed = false
		check.Score = 20
		check.Detail = "账号被限流"
	case 3:
		check.Passed = false
		check.Score = 0
		check.Detail = "账号被封禁"
	default:
		check.Passed = false
		check.Score = 50
		check.Detail = fmt.Sprintf("未知状态: %d", account.Status)
	}

	return check
}

// checkHealthTrend 检查健康分数趋势
func (s *HealthCheckService) checkHealthTrend(ctx context.Context, accountID int64) dal.CheckItem {
	check := dal.CheckItem{
		Name: "健康趋势",
	}

	history, err := s.accountRepo.GetHealthHistory(ctx, accountID, 5)
	if err != nil || len(history) < 2 {
		check.Passed = true
		check.Score = 80
		check.Detail = "历史数据不足，默认通过"
		return check
	}

	// 计算趋势
	latest := history[0].HealthScore
	previous := history[1].HealthScore

	if latest >= previous {
		check.Passed = true
		check.Score = 100
		check.Detail = fmt.Sprintf("健康分数稳定或上升 (%.1f → %.1f)", previous, latest)
	} else if latest >= previous-10 {
		check.Passed = true
		check.Score = 70
		check.Detail = fmt.Sprintf("健康分数轻微下降 (%.1f → %.1f)", previous, latest)
	} else {
		check.Passed = false
		check.Score = 30
		check.Detail = fmt.Sprintf("健康分数大幅下降 (%.1f → %.1f)", previous, latest)
	}

	return check
}

// checkPublishSuccessRate 检查发布成功率
func (s *HealthCheckService) checkPublishSuccessRate(ctx context.Context, accountID int64) dal.CheckItem {
	check := dal.CheckItem{
		Name: "发布成功率",
	}

	// 从账号健康历史中获取发布相关数据
	history, err := s.accountRepo.GetHealthHistory(ctx, accountID, 10)
	if err != nil || len(history) == 0 {
		check.Passed = true
		check.Score = 80
		check.Detail = "无发布历史，默认通过"
		return check
	}

	// 统计最近的健康状态
	normalCount := 0
	for _, h := range history {
		if h.Status == "normal" {
			normalCount++
		}
	}

	rate := float32(normalCount) / float32(len(history)) * 100

	if rate >= 80 {
		check.Passed = true
		check.Score = 100
		check.Detail = fmt.Sprintf("状态正常率 %.0f%%", rate)
	} else if rate >= 50 {
		check.Passed = true
		check.Score = 60
		check.Detail = fmt.Sprintf("状态正常率 %.0f%%，需关注", rate)
	} else {
		check.Passed = false
		check.Score = 20
		check.Detail = fmt.Sprintf("状态正常率 %.0f%%，建议检查账号", rate)
	}

	return check
}

// checkPlatformRisk 检查平台风控
func (s *HealthCheckService) checkPlatformRisk(account *dal.Account) dal.CheckItem {
	check := dal.CheckItem{
		Name: "平台风控",
	}

	// 基于健康分数判断
	score := account.HealthScore

	if score >= 80 {
		check.Passed = true
		check.Score = 100
		check.Detail = "无风控风险"
	} else if score >= 60 {
		check.Passed = true
		check.Score = 70
		check.Detail = "低风险，建议关注"
	} else if score >= 40 {
		check.Passed = false
		check.Score = 40
		check.Detail = "中等风险，建议减少发布频率"
	} else {
		check.Passed = false
		check.Score = 10
		check.Detail = "高风险，建议暂停发布"
	}

	return check
}

// ==================== 状态判定 ====================

// determineStatus 确定账号状态
func (s *HealthCheckService) determineStatus(score float32, checks []dal.CheckItem) (string, string) {
	cfg := config.GetConfig()

	// 检查是否有高危项
	hasHighRisk := false
	for _, check := range checks {
		if !check.Passed && check.Score < 20 {
			hasHighRisk = true
			break
		}
	}

	if hasHighRisk || score < cfg.HealthCheck.PauseScoreThreshold {
		return "critical", "high"
	}

	if score < cfg.HealthCheck.AlertScoreThreshold {
		return "warning", "medium"
	}

	if score < 80 {
		return "attention", "low"
	}

	return "normal", "none"
}

// determinePause 确定是否暂停发布
func (s *HealthCheckService) determinePause(score float32, checks []dal.CheckItem) (bool, string) {
	cfg := config.GetConfig()

	// 分数过低，暂停发布
	if score < cfg.HealthCheck.PauseScoreThreshold {
		return true, fmt.Sprintf("健康分数过低 (%.1f)，自动暂停发布", score)
	}

	// 检查是否有封禁/限流
	for _, check := range checks {
		if check.Name == "账号状态" && !check.Passed {
			if strings.Contains(check.Detail, "封禁") || strings.Contains(check.Detail, "限流") {
				return true, check.Detail
			}
		}
	}

	return false, ""
}

// ==================== 告警系统 ====================

// SendAlert 发送告警
func (s *HealthCheckService) SendAlert(ctx context.Context, result *dal.HealthCheckResult) error {
	if result.AlertLevel == "none" {
		return nil
	}

	checksJSON, _ := json.Marshal(result.Checks)

	health := &dal.AccountHealth{
		AccountID:     result.AccountID,
		HealthScore:   result.HealthScore,
		Status:        result.Status,
		CheckType:     "auto",
		CheckDetails:  string(checksJSON),
		AlertLevel:    result.AlertLevel,
		PublishPaused: result.PublishPaused,
		PauseReason:   result.PauseReason,
		CheckedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}

	channels := make([]string, 0)

	// 发送Webhook告警
	if err := s.sendWebhookAlert(ctx, result); err == nil {
		channels = append(channels, "webhook")
	}

	// 发送钉钉告警
	if err := s.sendDingTalkAlert(ctx, result); err == nil {
		channels = append(channels, "dingtalk")
	}

	// 发送邮件告警
	if err := s.sendEmailAlert(ctx, result); err == nil {
		channels = append(channels, "email")
	}

	health.AlertSent = len(channels) > 0
	health.AlertChannels = strings.Join(channels, ",")

	// 保存告警记录
	record := &dal.AlertRecord{
		AccountID: result.AccountID,
		HealthID:  health.ID,
		AlertType: result.AlertLevel,
		Channel:   strings.Join(channels, ","),
		Title:     fmt.Sprintf("账号健康告警 - %s", result.Status),
		Content:   fmt.Sprintf("账号ID: %d, 健康分数: %.1f, 状态: %s", result.AccountID, result.HealthScore, result.Status),
		Success:   len(channels) > 0,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.alertRepo.Create(ctx, record); err != nil {
		return fmt.Errorf("save alert record: %w", err)
	}

	return nil
}

// sendWebhookAlert 发送Webhook告警
func (s *HealthCheckService) sendWebhookAlert(ctx context.Context, result *dal.HealthCheckResult) error {
	payload := map[string]interface{}{
		"event":        "account_health_alert",
		"account_id":   result.AccountID,
		"health_score": result.HealthScore,
		"status":       result.Status,
		"alert_level":  result.AlertLevel,
		"paused":       result.PublishPaused,
		"timestamp":    time.Now().Unix(),
	}

	payloadJSON, _ := json.Marshal(payload)
	_ = payloadJSON

	// TODO: 从配置获取Webhook URL并发送
	// webhookURL := s.config.GetWebhookURL()
	// req, _ := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewReader(payloadJSON))
	// req.Header.Set("Content-Type", "application/json")
	// resp, err := s.httpClient.Do(req)

	return nil
}

// sendDingTalkAlert 发送钉钉告警
func (s *HealthCheckService) sendDingTalkAlert(ctx context.Context, result *dal.HealthCheckResult) error {
	// TODO: 从配置获取钉钉机器人Webhook URL
	// 构建钉钉消息格式
	title := fmt.Sprintf("⚠️ 账号健康告警")
	text := fmt.Sprintf("### %s\n\n- **账号ID**: %d\n- **健康分数**: %.1f\n- **状态**: %s\n- **告警级别**: %s\n- **发布时间**: %s",
		title, result.AccountID, result.HealthScore, result.Status, result.AlertLevel, time.Now().Format("2006-01-02 15:04:05"))

	if result.PublishPaused {
		text += fmt.Sprintf("\n- **⚠️ 已暂停发布**: %s", result.PauseReason)
	}

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}

	msgJSON, _ := json.Marshal(msg)
	_ = msgJSON

	// TODO: 发送到钉钉Webhook
	// dingTalkURL := s.config.GetDingTalkWebhook()
	// req, _ := http.NewRequestWithContext(ctx, "POST", dingTalkURL, bytes.NewReader(msgJSON))
	// req.Header.Set("Content-Type", "application/json")
	// resp, err := s.httpClient.Do(req)

	return nil
}

// sendEmailAlert 发送邮件告警
func (s *HealthCheckService) sendEmailAlert(ctx context.Context, result *dal.HealthCheckResult) error {
	subject := fmt.Sprintf("[OpenGEO] 账号健康告警 - 账号ID: %d", result.AccountID)
	body := fmt.Sprintf(`
账号健康检测告警

账号ID: %d
健康分数: %.1f
状态: %s
告警级别: %s
检测时间: %s

检测详情:
%s

%s
`,
		result.AccountID,
		result.HealthScore,
		result.Status,
		result.AlertLevel,
		time.Now().Format("2006-01-02 15:04:05"),
		formatChecks(result.Checks),
		func() string {
			if result.PublishPaused {
				return fmt.Sprintf("⚠️ 已自动暂停发布: %s", result.PauseReason)
			}
			return ""
		}(),
	)

	_ = subject
	_ = body

	// TODO: 从配置获取SMTP设置并发送邮件
	// smtpHost := s.config.GetSMTPHost()
	// smtpPort := s.config.GetSMTPPort()
	// to := s.config.GetAlertEmail()
	// err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, []byte(message))

	return nil
}

// ==================== 辅助函数 ====================

// serializeChecks 序列化检查项
func (s *HealthCheckService) serializeChecks(checks []dal.CheckItem) string {
	data, _ := json.Marshal(checks)
	return string(data)
}

// formatChecks 格式化检查项
func formatChecks(checks []dal.CheckItem) string {
	var sb strings.Builder
	for _, check := range checks {
		status := "✅"
		if !check.Passed {
			status = "❌"
		}
		sb.WriteString(fmt.Sprintf("%s %s: %s (分数: %.0f)\n", status, check.Name, check.Detail, check.Score))
	}
	return sb.String()
}
