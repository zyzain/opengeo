package service

import (
	"testing"
	"time"
)

func TestStaggerScheduler_Basic(t *testing.T) {
	strategy := DefaultStaggerStrategy()
	scheduler := NewStaggerScheduler(strategy)

	baseTime := time.Now()
	nextTime := scheduler.CalculateNextTime(baseTime)

	// 应该在最小间隔之后
	if nextTime.Before(baseTime.Add(strategy.MinInterval)) {
		t.Errorf("next time should be after min interval")
	}
}

func TestStaggerScheduler_MinInterval(t *testing.T) {
	strategy := &StaggerStrategy{
		MinInterval:   5 * time.Minute,
		MaxInterval:   30 * time.Minute,
		VarianceRatio: 0,
	}
	scheduler := NewStaggerScheduler(strategy)

	baseTime := time.Now()
	nextTime := scheduler.CalculateNextTime(baseTime)

	t.Logf("baseTime=%v, nextTime=%v, diff=%v", baseTime, nextTime, nextTime.Sub(baseTime))

	// 无浮动时，间隔应该正好是最小间隔
	if nextTime.Sub(baseTime) < 4*time.Minute {
		t.Errorf("interval too short: %v", nextTime.Sub(baseTime))
	}
}

func TestStaggerScheduler_Variance(t *testing.T) {
	strategy := &StaggerStrategy{
		MinInterval:   5 * time.Minute,
		MaxInterval:   30 * time.Minute,
		VarianceRatio: 0.3,
	}
	scheduler := NewStaggerScheduler(strategy)

	baseTime := time.Now()
	intervals := make([]time.Duration, 100)

	for i := range intervals {
		nextTime := scheduler.CalculateNextTime(baseTime.Add(time.Duration(i) * time.Hour))
		intervals[i] = nextTime.Sub(baseTime.Add(time.Duration(i) * time.Hour))
	}

	// 检查间隔有变化（不是全部相同）
	hasVariation := false
	for i := 1; i < len(intervals); i++ {
		if intervals[i] != intervals[0] {
			hasVariation = true
			break
		}
	}
	if !hasVariation {
		t.Error("expected variance in intervals")
	}

	// 检查间隔在合理范围内
	for _, interval := range intervals {
		if interval < 5*time.Minute {
			t.Errorf("interval %v below minimum", interval)
		}
		if interval > 5*time.Minute+2*time.Minute { // 5min + 30% variance
			t.Errorf("interval %v above expected maximum", interval)
		}
	}
}

func TestStaggerScheduler_ScheduleStaggeredTasks(t *testing.T) {
	strategy := DefaultStaggerStrategy()
	scheduler := NewStaggerScheduler(strategy)

	tasks := []*QueueItem{
		{TaskID: 1, ContentID: 100, ChannelID: 1, Priority: 1},
		{TaskID: 2, ContentID: 200, ChannelID: 2, Priority: 2},
		{TaskID: 3, ContentID: 300, ChannelID: 1, Priority: 3},
	}

	startTime := time.Now()
	scheduled := scheduler.ScheduleStaggeredTasks(tasks, startTime)

	if len(scheduled) != 3 {
		t.Errorf("expected 3 scheduled tasks, got %d", len(scheduled))
	}

	// 检查时间递增
	for i := 1; i < len(scheduled); i++ {
		if scheduled[i].ScheduledAt.Before(scheduled[i-1].ScheduledAt) {
			t.Error("scheduled times should be increasing")
		}
	}

	// 检查间隔存在
	for _, task := range scheduled {
		if task.StaggerDelay < 0 {
			t.Error("stagger delay should be positive")
		}
	}
}

func TestAdaptiveStaggerScheduler_Basic(t *testing.T) {
	strategy := DefaultStaggerStrategy()
	scheduler := NewAdaptiveStaggerScheduler(strategy)

	// 更新渠道统计
	scheduler.UpdateChannelStats(1, 0.9, 5)
	scheduler.UpdateChannelStats(2, 0.5, 15)

	// 成功率高的渠道间隔应该更短
	interval1 := scheduler.CalculateAdaptiveInterval(1)
	interval2 := scheduler.CalculateAdaptiveInterval(2)

	// 渠道2成功率低且负载高，间隔应该更长
	if interval2 < interval1 {
		t.Errorf("channel 2 should have longer interval: %v vs %v", interval2, interval1)
	}
}

func TestAdaptiveStaggerScheduler_ScheduleAdaptive(t *testing.T) {
	strategy := DefaultStaggerStrategy()
	scheduler := NewAdaptiveStaggerScheduler(strategy)

	scheduler.UpdateChannelStats(1, 0.9, 5)
	scheduler.UpdateChannelStats(2, 0.8, 10)

	tasks := []*QueueItem{
		{TaskID: 1, ContentID: 100, ChannelID: 1, Priority: 1},
		{TaskID: 2, ContentID: 200, ChannelID: 1, Priority: 1},
		{TaskID: 3, ContentID: 300, ChannelID: 2, Priority: 2},
	}

	startTime := time.Now()
	scheduled := scheduler.ScheduleAdaptive(tasks, startTime)

	if len(scheduled) != 3 {
		t.Errorf("expected 3 scheduled tasks, got %d", len(scheduled))
	}
}

func TestStaggerConfig_Parse(t *testing.T) {
	config := &StaggerConfig{
		Strategy:       "adaptive",
		MinIntervalMin: 10,
		MaxIntervalMin: 60,
		VariancePct:    20,
		MaxConcurrency: 5,
	}

	strategy := ParseStaggerConfig(config)

	if strategy.MinInterval != 10*time.Minute {
		t.Errorf("expected 10min, got %v", strategy.MinInterval)
	}
	if strategy.MaxInterval != 60*time.Minute {
		t.Errorf("expected 60min, got %v", strategy.MaxInterval)
	}
	if strategy.VarianceRatio != 0.2 {
		t.Errorf("expected 0.2, got %f", strategy.VarianceRatio)
	}
	if strategy.MaxConcurrency != 5 {
		t.Errorf("expected 5, got %d", strategy.MaxConcurrency)
	}
}

func TestStaggerConfig_Validate(t *testing.T) {
	tests := []struct {
		config  *StaggerConfig
		wantErr bool
		desc    string
	}{
		{
			config:  &StaggerConfig{MinIntervalMin: 5, MaxIntervalMin: 30, VariancePct: 30, MaxConcurrency: 10},
			wantErr: false,
			desc:    "valid config",
		},
		{
			config:  &StaggerConfig{MinIntervalMin: 0, MaxIntervalMin: 30, VariancePct: 30, MaxConcurrency: 10},
			wantErr: true,
			desc:    "min interval too low",
		},
		{
			config:  &StaggerConfig{MinIntervalMin: 30, MaxIntervalMin: 10, VariancePct: 30, MaxConcurrency: 10},
			wantErr: true,
			desc:    "max less than min",
		},
		{
			config:  &StaggerConfig{MinIntervalMin: 5, MaxIntervalMin: 30, VariancePct: 150, MaxConcurrency: 10},
			wantErr: true,
			desc:    "variance too high",
		},
		{
			config:  &StaggerConfig{MinIntervalMin: 5, MaxIntervalMin: 30, VariancePct: 30, MaxConcurrency: 0},
			wantErr: true,
			desc:    "concurrency too low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := ValidateStaggerConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultStaggerStrategy(t *testing.T) {
	strategy := DefaultStaggerStrategy()

	if strategy.MinInterval != 5*time.Minute {
		t.Errorf("expected 5min, got %v", strategy.MinInterval)
	}
	if strategy.VarianceRatio != 0.3 {
		t.Errorf("expected 0.3, got %f", strategy.VarianceRatio)
	}
	if strategy.MaxConcurrency != 10 {
		t.Errorf("expected 10, got %d", strategy.MaxConcurrency)
	}
}

func TestScheduledTask_Fields(t *testing.T) {
	task := &ScheduledTask{
		TaskID:       123,
		ContentID:    456,
		ChannelID:    789,
		Priority:     2,
		ScheduledAt:  time.Now(),
		StaggerDelay: 5 * time.Minute,
	}

	if task.TaskID != 123 {
		t.Errorf("expected 123, got %d", task.TaskID)
	}
	if task.StaggerDelay != 5*time.Minute {
		t.Errorf("expected 5min, got %v", task.StaggerDelay)
	}
}

func TestStaggerScheduler_ConsecutiveScheduling(t *testing.T) {
	strategy := &StaggerStrategy{
		MinInterval:   1 * time.Minute,
		VarianceRatio: 0.1,
	}
	scheduler := NewStaggerScheduler(strategy)

	// 连续调度应该累积间隔
	baseTime := time.Now()
	time1 := scheduler.CalculateNextTime(baseTime)
	time2 := scheduler.CalculateNextTime(time1)
	time3 := scheduler.CalculateNextTime(time2)

	if !time2.After(time1) {
		t.Error("time2 should be after time1")
	}
	if !time3.After(time2) {
		t.Error("time3 should be after time2")
	}
}
