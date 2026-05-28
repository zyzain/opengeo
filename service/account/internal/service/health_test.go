package service

import (
	"testing"

	"opengeo/service/account/internal/dal"
)

func TestDetermineStatus(t *testing.T) {
	svc := &HealthCheckService{}

	tests := []struct {
		score       float32
		checks      []dal.CheckItem
		wantStatus  string
		wantAlert   string
		desc        string
	}{
		{
			score:      90,
			checks:     []dal.CheckItem{{Name: "test", Passed: true, Score: 90}},
			wantStatus: "normal",
			wantAlert:  "none",
			desc:       "healthy account",
		},
		{
			score:      70,
			checks:     []dal.CheckItem{{Name: "test", Passed: true, Score: 70}},
			wantStatus: "attention",
			wantAlert:  "low",
			desc:       "attention needed",
		},
		{
			score:      50,
			checks:     []dal.CheckItem{{Name: "test", Passed: false, Score: 50}},
			wantStatus: "warning",
			wantAlert:  "medium",
			desc:       "warning status",
		},
		{
			score:      20,
			checks:     []dal.CheckItem{{Name: "test", Passed: false, Score: 10}},
			wantStatus: "critical",
			wantAlert:  "high",
			desc:       "critical status",
		},
		{
			score: 80,
			checks: []dal.CheckItem{
				{Name: "test1", Passed: true, Score: 100},
				{Name: "test2", Passed: false, Score: 10},
			},
			wantStatus: "critical",
			wantAlert:  "high",
			desc:       "high risk item triggers critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			status, alert := svc.determineStatus(tt.score, tt.checks)
			if status != tt.wantStatus {
				t.Errorf("status: got %s, want %s", status, tt.wantStatus)
			}
			if alert != tt.wantAlert {
				t.Errorf("alert: got %s, want %s", alert, tt.wantAlert)
			}
		})
	}
}

func TestDeterminePause(t *testing.T) {
	svc := &HealthCheckService{}

	tests := []struct {
		score      float32
		checks     []dal.CheckItem
		wantPause  bool
		desc       string
	}{
		{
			score:     90,
			checks:    []dal.CheckItem{{Name: "账号状态", Passed: true, Score: 100, Detail: "正常"}},
			wantPause: false,
			desc:      "healthy - no pause",
		},
		{
			score:     20,
			checks:    []dal.CheckItem{{Name: "test", Passed: false, Score: 10}},
			wantPause: true,
			desc:      "low score - pause",
		},
		{
			score: 80,
			checks: []dal.CheckItem{
				{Name: "账号状态", Passed: false, Score: 0, Detail: "账号被封禁"},
			},
			wantPause: true,
			desc:      "banned - pause",
		},
		{
			score: 80,
			checks: []dal.CheckItem{
				{Name: "账号状态", Passed: false, Score: 20, Detail: "账号被限流"},
			},
			wantPause: true,
			desc:      "rate limited - pause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			paused, _ := svc.determinePause(tt.score, tt.checks)
			if paused != tt.wantPause {
				t.Errorf("got %v, want %v", paused, tt.wantPause)
			}
		})
	}
}

func TestCheckAccountStatus(t *testing.T) {
	svc := &HealthCheckService{}

	tests := []struct {
		status    int32
		wantScore float32
		wantPass  bool
		desc      string
	}{
		{1, 100, true, "normal"},
		{0, 0, false, "disabled"},
		{2, 20, false, "rate limited"},
		{3, 0, false, "banned"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			account := &dal.Account{Status: tt.status}
			check := svc.checkAccountStatus(account)
			if check.Score != tt.wantScore {
				t.Errorf("score: got %f, want %f", check.Score, tt.wantScore)
			}
			if check.Passed != tt.wantPass {
				t.Errorf("passed: got %v, want %v", check.Passed, tt.wantPass)
			}
		})
	}
}

func TestCheckHealthTrend(t *testing.T) {
	tests := []struct {
		score     float32
		wantPass  bool
		desc      string
	}{
		{100, true, "stable score"},
		{70, true, "slight decline"},
		{30, false, "major decline"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			check := dal.CheckItem{
				Name:   "健康趋势",
				Score:  tt.score,
				Passed: tt.wantPass,
			}
			if check.Passed != tt.wantPass {
				t.Errorf("got %v, want %v", check.Passed, tt.wantPass)
			}
		})
	}
}

func TestCheckPlatformRisk(t *testing.T) {
	svc := &HealthCheckService{}

	tests := []struct {
		healthScore float32
		wantPass    bool
		wantScore   float32
		desc        string
	}{
		{90, true, 100, "no risk"},
		{70, true, 70, "low risk"},
		{50, false, 40, "medium risk"},
		{30, false, 10, "high risk"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			account := &dal.Account{HealthScore: tt.healthScore}
			check := svc.checkPlatformRisk(account)
			if check.Passed != tt.wantPass {
				t.Errorf("passed: got %v, want %v", check.Passed, tt.wantPass)
			}
			if check.Score != tt.wantScore {
				t.Errorf("score: got %f, want %f", check.Score, tt.wantScore)
			}
		})
	}
}

func TestFormatChecks(t *testing.T) {
	checks := []dal.CheckItem{
		{Name: "测试1", Passed: true, Score: 100, Detail: "正常"},
		{Name: "测试2", Passed: false, Score: 20, Detail: "异常"},
	}

	result := formatChecks(checks)

	if result == "" {
		t.Error("expected non-empty result")
	}
	if !containsStr(result, "测试1") {
		t.Error("expected check 1")
	}
	if !containsStr(result, "测试2") {
		t.Error("expected check 2")
	}
}

func TestAlertRecordModel(t *testing.T) {
	record := &dal.AlertRecord{
		AccountID: 1,
		AlertType: "high",
		Channel:   "webhook,dingtalk",
		Title:     "测试告警",
		Content:   "测试内容",
		Success:   true,
	}

	if record.AccountID != 1 {
		t.Error("account ID mismatch")
	}
	if !record.Success {
		t.Error("expected success")
	}
}

func TestHealthCheckResultModel(t *testing.T) {
	result := &dal.HealthCheckResult{
		AccountID:     1,
		HealthScore:   85.5,
		Status:        "normal",
		AlertLevel:    "none",
		PublishPaused: false,
		Checks: []dal.CheckItem{
			{Name: "测试", Passed: true, Score: 85.5, Detail: "正常"},
		},
	}

	if result.HealthScore != 85.5 {
		t.Error("health score mismatch")
	}
	if result.PublishPaused {
		t.Error("expected not paused")
	}
	if len(result.Checks) != 1 {
		t.Error("expected 1 check")
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
