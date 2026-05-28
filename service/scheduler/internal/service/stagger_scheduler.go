package service

import (
	"fmt"
	"math/rand"
	"time"

	"opengeo/pkg/config"
)

// StaggerStrategy 错峰策略
type StaggerStrategy struct {
	MinInterval    time.Duration `json:"min_interval"`     // 最小间隔（默认5分钟）
	MaxInterval    time.Duration `json:"max_interval"`     // 最大间隔
	VarianceRatio  float32       `json:"variance_ratio"`   // 浮动比例（默认0.3即±30%）
	MaxConcurrency int           `json:"max_concurrency"`  // 最大并发数
	BurstLimit     int           `json:"burst_limit"`      // 突发限制
}

// DefaultStaggerStrategy 默认错峰策略
func DefaultStaggerStrategy() *StaggerStrategy {
	cfg := config.GetConfig()
	return &StaggerStrategy{
		MinInterval:    cfg.Stagger.MinInterval,
		MaxInterval:    cfg.Stagger.MaxInterval,
		VarianceRatio:  cfg.Stagger.VarianceRatio,
		MaxConcurrency: cfg.Stagger.MaxConcurrency,
		BurstLimit:     cfg.Stagger.BurstLimit,
	}
}

// StaggerScheduler 错峰调度器
type StaggerScheduler struct {
	strategy    *StaggerStrategy
	lastPublish time.Time
	rng         *rand.Rand
}

// NewStaggerScheduler 创建错峰调度器
func NewStaggerScheduler(strategy *StaggerStrategy) *StaggerScheduler {
	if strategy == nil {
		strategy = DefaultStaggerStrategy()
	}
	return &StaggerScheduler{
		strategy:    strategy,
		lastPublish: time.Time{},
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateNextTime 计算下次发布时间
func (s *StaggerScheduler) CalculateNextTime(baseTime time.Time) time.Time {
	// 确保最小间隔
	earliest := s.lastPublish.Add(s.strategy.MinInterval)
	if baseTime.Before(earliest) {
		baseTime = earliest
	}

	// 添加随机浮动
	interval := s.strategy.MinInterval
	variance := time.Duration(float32(interval) * s.strategy.VarianceRatio)

	if variance > 0 {
		// 随机增减浮动
		delta := time.Duration(s.rng.Int63n(int64(variance)))
		if s.rng.Float64() > 0.5 {
			interval += delta
		} else {
			interval -= delta
		}
	}

	// 确保不小于最小间隔
	if interval < s.strategy.MinInterval {
		interval = s.strategy.MinInterval
	}

	// 确保不超过最大间隔
	if interval > s.strategy.MaxInterval {
		interval = s.strategy.MaxInterval
	}

	nextTime := baseTime.Add(interval)
	s.lastPublish = nextTime

	return nextTime
}

// ScheduleStaggeredTasks 错峰调度多个任务
func (s *StaggerScheduler) ScheduleStaggeredTasks(tasks []*QueueItem, startTime time.Time) []*ScheduledTask {
	result := make([]*ScheduledTask, 0, len(tasks))
	currentTime := startTime

	for _, task := range tasks {
		scheduledTime := s.CalculateNextTime(currentTime)

		result = append(result, &ScheduledTask{
			TaskID:       task.TaskID,
			ContentID:    task.ContentID,
			ChannelID:    task.ChannelID,
			Priority:     task.Priority,
			ScheduledAt:  scheduledTime,
			StaggerDelay: scheduledTime.Sub(currentTime),
		})

		currentTime = scheduledTime
	}

	return result
}

// ScheduledTask 已调度任务
type ScheduledTask struct {
	TaskID       int64         `json:"task_id"`
	ContentID    int64         `json:"content_id"`
	ChannelID    int64         `json:"channel_id"`
	Priority     int32         `json:"priority"`
	ScheduledAt  time.Time     `json:"scheduled_at"`
	StaggerDelay time.Duration `json:"stagger_delay"`
}

// ==================== 高级错峰策略 ====================

// AdaptiveStaggerScheduler 自适应错峰调度器
type AdaptiveStaggerScheduler struct {
	baseStrategy  *StaggerStrategy
	successRates  map[int64]float32 // 渠道成功率
	channelLoads  map[int64]int     // 渠道当前负载
	rng           *rand.Rand
}

// NewAdaptiveStaggerScheduler 创建自适应错峰调度器
func NewAdaptiveStaggerScheduler(strategy *StaggerStrategy) *AdaptiveStaggerScheduler {
	if strategy == nil {
		strategy = DefaultStaggerStrategy()
	}
	return &AdaptiveStaggerScheduler{
		baseStrategy: strategy,
		successRates: make(map[int64]float32),
		channelLoads: make(map[int64]int),
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// UpdateChannelStats 更新渠道统计
func (s *AdaptiveStaggerScheduler) UpdateChannelStats(channelID int64, successRate float32, load int) {
	s.successRates[channelID] = successRate
	s.channelLoads[channelID] = load
}

// CalculateAdaptiveInterval 计算自适应间隔
func (s *AdaptiveStaggerScheduler) CalculateAdaptiveInterval(channelID int64) time.Duration {
	baseInterval := s.baseStrategy.MinInterval

	// 根据成功率调整
	successRate, exists := s.successRates[channelID]
	if exists && successRate < 0.8 {
		// 成功率低，增加间隔
		penalty := time.Duration((1 - successRate) * 10) * time.Minute
		baseInterval += penalty
	}

	// 根据负载调整
	load, exists := s.channelLoads[channelID]
	if exists && load > s.baseStrategy.MaxConcurrency {
		// 负载高，增加间隔
		overload := time.Duration(load-s.baseStrategy.MaxConcurrency) * time.Minute
		baseInterval += overload
	}

	// 添加随机浮动
	variance := time.Duration(float32(baseInterval) * s.baseStrategy.VarianceRatio)
	if variance > 0 {
		delta := time.Duration(s.rng.Int63n(int64(variance)))
		if s.rng.Float64() > 0.5 {
			baseInterval += delta
		} else {
			baseInterval -= delta
		}
	}

	// 确保范围
	if baseInterval < s.baseStrategy.MinInterval {
		baseInterval = s.baseStrategy.MinInterval
	}
	if baseInterval > s.baseStrategy.MaxInterval {
		baseInterval = s.baseStrategy.MaxInterval
	}

	return baseInterval
}

// ScheduleAdaptive 自适应错峰调度
func (s *AdaptiveStaggerScheduler) ScheduleAdaptive(tasks []*QueueItem, startTime time.Time) []*ScheduledTask {
	result := make([]*ScheduledTask, 0, len(tasks))
	currentTime := startTime

	// 按渠道分组
	channelTasks := make(map[int64][]*QueueItem)
	for _, task := range tasks {
		channelTasks[task.ChannelID] = append(channelTasks[task.ChannelID], task)
	}

	// 为每个渠道计算独立的调度
	for channelID, channelTaskList := range channelTasks {
		channelTime := currentTime

		for _, task := range channelTaskList {
			interval := s.CalculateAdaptiveInterval(channelID)
			scheduledTime := channelTime.Add(interval)

			result = append(result, &ScheduledTask{
				TaskID:       task.TaskID,
				ContentID:    task.ContentID,
				ChannelID:    task.ChannelID,
				Priority:     task.Priority,
				ScheduledAt:  scheduledTime,
				StaggerDelay: interval,
			})

			channelTime = scheduledTime
		}
	}

	return result
}

// ==================== 错峰配置管理 ====================

// StaggerConfig 错峰配置
type StaggerConfig struct {
	Strategy       string `json:"strategy"`        // fixed, adaptive
	MinIntervalMin int    `json:"min_interval_min"` // 最小间隔（分钟）
	MaxIntervalMin int    `json:"max_interval_min"` // 最大间隔（分钟）
	VariancePct    int    `json:"variance_pct"`     // 浮动百分比
	MaxConcurrency int    `json:"max_concurrency"`  // 最大并发
}

// ParseStaggerConfig 解析错峰配置
func ParseStaggerConfig(config *StaggerConfig) *StaggerStrategy {
	strategy := DefaultStaggerStrategy()

	if config.MinIntervalMin > 0 {
		strategy.MinInterval = time.Duration(config.MinIntervalMin) * time.Minute
	}
	if config.MaxIntervalMin > 0 {
		strategy.MaxInterval = time.Duration(config.MaxIntervalMin) * time.Minute
	}
	if config.VariancePct > 0 {
		strategy.VarianceRatio = float32(config.VariancePct) / 100.0
	}
	if config.MaxConcurrency > 0 {
		strategy.MaxConcurrency = config.MaxConcurrency
	}

	return strategy
}

// ValidateStaggerConfig 验证错峰配置
func ValidateStaggerConfig(config *StaggerConfig) error {
	if config.MinIntervalMin < 1 {
		return fmt.Errorf("min_interval_min must be >= 1")
	}
	if config.MaxIntervalMin < config.MinIntervalMin {
		return fmt.Errorf("max_interval_min must be >= min_interval_min")
	}
	if config.VariancePct < 0 || config.VariancePct > 100 {
		return fmt.Errorf("variance_pct must be between 0 and 100")
	}
	if config.MaxConcurrency < 1 {
		return fmt.Errorf("max_concurrency must be >= 1")
	}
	return nil
}
