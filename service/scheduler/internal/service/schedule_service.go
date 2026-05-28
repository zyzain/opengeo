package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/scheduler/internal/dal"
	"opengeo/service/scheduler/internal/domain/model"
)

// ScheduleService 调度服务
type ScheduleService struct {
	scheduleRepo    *dal.ScheduleRepository
	taskRepo        *dal.ScheduleTaskRepository
	heatmapRepo     *dal.AIActivityHeatmapRepository
	calendarRepo    *dal.PublishCalendarRepository
	logRepo         *dal.ScheduleLogRepository
	cronParser      *CronParser
}

// NewScheduleService 创建调度服务
func NewScheduleService(
	scheduleRepo *dal.ScheduleRepository,
	taskRepo *dal.ScheduleTaskRepository,
	heatmapRepo *dal.AIActivityHeatmapRepository,
	calendarRepo *dal.PublishCalendarRepository,
	logRepo *dal.ScheduleLogRepository,
) *ScheduleService {
	return &ScheduleService{
		scheduleRepo:  scheduleRepo,
		taskRepo:      taskRepo,
		heatmapRepo:   heatmapRepo,
		calendarRepo:  calendarRepo,
		logRepo:       logRepo,
		cronParser:    NewCronParser(),
	}
}

// ==================== 调度管理 ====================

// CreateSchedule 创建调度
func (s *ScheduleService) CreateSchedule(ctx context.Context, userID int64, name, scheduleType, cronExpression, config string) (*model.Schedule, error) {
	// 验证调度类型
	if !isValidScheduleType(scheduleType) {
		return nil, fmt.Errorf("invalid schedule type: %s", scheduleType)
	}

	schedule := &model.Schedule{
		UserID:         userID,
		ScheduleName:   name,
		ScheduleType:   scheduleType,
		CronExpression: cronExpression,
		Config:         config,
		IsEnabled:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 计算下次运行时间
	if scheduleType == model.ScheduleTypeFixed || scheduleType == model.ScheduleTypeInterval {
		nextRunTime := s.calculateNextRunTime(scheduleType, cronExpression)
		schedule.NextRunTime = &nextRunTime
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	// 记录日志
	s.logRepo.Create(ctx, &model.ScheduleLog{
		ScheduleID: schedule.ID,
		Action:     "created",
		Message:    fmt.Sprintf("调度创建: %s", name),
		CreatedAt:  time.Now(),
	})

	return schedule, nil
}

// GetSchedule 获取调度
func (s *ScheduleService) GetSchedule(ctx context.Context, scheduleID int64) (*model.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return schedule, nil
}

// UpdateSchedule 更新调度
func (s *ScheduleService) UpdateSchedule(ctx context.Context, scheduleID int64, name, cronExpression, config string) (*model.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	if name != "" {
		schedule.ScheduleName = name
	}
	if cronExpression != "" {
		schedule.CronExpression = cronExpression
		// 重新计算下次运行时间
		nextRunTime := s.calculateNextRunTime(schedule.ScheduleType, cronExpression)
		schedule.NextRunTime = &nextRunTime
	}
	if config != "" {
		schedule.Config = config
	}
	schedule.UpdatedAt = time.Now()

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return schedule, nil
}

// DeleteSchedule 删除调度
func (s *ScheduleService) DeleteSchedule(ctx context.Context, scheduleID int64) error {
	if err := s.scheduleRepo.Delete(ctx, scheduleID); err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	return nil
}

// ListSchedules 列出调度
func (s *ScheduleService) ListSchedules(ctx context.Context, userID int64, isEnabled *bool, page, pageSize int) ([]*model.Schedule, int32, error) {
	schedules, total, err := s.scheduleRepo.List(ctx, userID, isEnabled, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list schedules: %w", err)
	}

	return schedules, total, nil
}

// EnableSchedule 启用调度
func (s *ScheduleService) EnableSchedule(ctx context.Context, scheduleID int64) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	schedule.IsEnabled = true
	schedule.UpdatedAt = time.Now()

	// 重新计算下次运行时间
	nextRunTime := s.calculateNextRunTime(schedule.ScheduleType, schedule.CronExpression)
	schedule.NextRunTime = &nextRunTime

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return fmt.Errorf("failed to enable schedule: %w", err)
	}

	return nil
}

// DisableSchedule 禁用调度
func (s *ScheduleService) DisableSchedule(ctx context.Context, scheduleID int64) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	schedule.IsEnabled = false
	schedule.NextRunTime = nil
	schedule.UpdatedAt = time.Now()

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return fmt.Errorf("failed to disable schedule: %w", err)
	}

	return nil
}

// ==================== 任务执行 ====================

// GetDueSchedules 获取到期的调度
func (s *ScheduleService) GetDueSchedules(ctx context.Context, limit int) ([]*model.Schedule, error) {
	schedules, err := s.scheduleRepo.GetDueSchedules(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get due schedules: %w", err)
	}

	return schedules, nil
}

// ExecuteSchedule 执行调度
func (s *ScheduleService) ExecuteSchedule(ctx context.Context, scheduleID int64) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	if !schedule.IsEnabled {
		return fmt.Errorf("schedule is disabled")
	}

	// 记录执行日志
	s.logRepo.Create(ctx, &model.ScheduleLog{
		ScheduleID: scheduleID,
		Action:     "trigger",
		Message:    fmt.Sprintf("调度触发: %s", schedule.ScheduleName),
		CreatedAt:  time.Now(),
	})

	// 更新下次运行时间
	nextRunTime := s.calculateNextRunTime(schedule.ScheduleType, schedule.CronExpression)
	if err := s.scheduleRepo.UpdateNextRunTime(ctx, scheduleID, nextRunTime); err != nil {
		return fmt.Errorf("failed to update next run time: %w", err)
	}

	return nil
}

// ==================== 热力图 ====================

// GetAIActivityHeatmap 获取AI活跃度热力图
func (s *ScheduleService) GetAIActivityHeatmap(ctx context.Context, platform, aiModel string) ([][]float32, error) {
	// 创建7x24矩阵
	heatmap := make([][]float32, 7)
	for i := range heatmap {
		heatmap[i] = make([]float32, 24)
	}

	// 填充数据
	for day := 1; day <= 7; day++ {
		for hour := 0; hour < 24; hour++ {
			data, err := s.heatmapRepo.Get(ctx, platform, aiModel, int32(hour), int32(day))
			if err != nil {
				continue
			}
			heatmap[day-1][hour] = data.ActivityScore
		}
	}

	return heatmap, nil
}

// GetBestTimeSlots 获取最佳发布时间
func (s *ScheduleService) GetBestTimeSlots(ctx context.Context, platform, aiModel string, limit int) ([]*model.AIActivityHeatmap, error) {
	slots, err := s.heatmapRepo.GetBestTimeSlots(ctx, platform, aiModel, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get best time slots: %w", err)
	}

	return slots, nil
}

// ==================== 发布日历 ====================

// CreateCalendarEvent 创建日历事件
func (s *ScheduleService) CreateCalendarEvent(ctx context.Context, userID int64, title, description string, startTime time.Time, endTime *time.Time, scheduleID, contentID, channelID *int64) (*model.PublishCalendar, error) {
	event := &model.PublishCalendar{
		UserID:      userID,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		ScheduleID:  scheduleID,
		ContentID:   contentID,
		ChannelID:   channelID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.calendarRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create calendar event: %w", err)
	}

	return event, nil
}

// GetCalendarEvents 获取日历事件
func (s *ScheduleService) GetCalendarEvents(ctx context.Context, userID int64, startTime, endTime time.Time) ([]*model.PublishCalendar, error) {
	events, err := s.calendarRepo.List(ctx, userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar events: %w", err)
	}

	return events, nil
}

// UpdateCalendarEvent 更新日历事件
func (s *ScheduleService) UpdateCalendarEvent(ctx context.Context, eventID int64, title, description string, startTime time.Time, endTime *time.Time) (*model.PublishCalendar, error) {
	event, err := s.calendarRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar event: %w", err)
	}

	if title != "" {
		event.Title = title
	}
	if description != "" {
		event.Description = description
	}
	if !startTime.IsZero() {
		event.StartTime = startTime
	}
	if endTime != nil {
		event.EndTime = endTime
	}
	event.UpdatedAt = time.Now()

	if err := s.calendarRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update calendar event: %w", err)
	}

	return event, nil
}

// DeleteCalendarEvent 删除日历事件
func (s *ScheduleService) DeleteCalendarEvent(ctx context.Context, eventID int64) error {
	if err := s.calendarRepo.Delete(ctx, eventID); err != nil {
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}

	return nil
}

// ==================== 辅助函数 ====================

// calculateNextRunTime 计算下次运行时间
func (s *ScheduleService) calculateNextRunTime(scheduleType, cronExpression string) time.Time {
	now := time.Now()

	switch scheduleType {
	case model.ScheduleTypeFixed:
		// 固定时间：解析Cron表达式
		nextTime, err := s.cronParser.Parse(cronExpression, now)
		if err != nil {
			// 解析失败，返回1小时后
			return now.Add(1 * time.Hour)
		}
		return nextTime

	case model.ScheduleTypeInterval:
		// 间隔循环：解析间隔表达式
		duration, err := s.cronParser.ParseInterval(cronExpression)
		if err != nil {
			// 解析失败，返回1小时后
			return now.Add(1 * time.Hour)
		}
		return now.Add(duration)

	case model.ScheduleTypeEvent:
		// 事件触发：由外部触发，不预设时间
		return now.Add(24 * time.Hour) // 默认24小时后过期

	case model.ScheduleTypeHeat:
		// 热力图推荐：使用热力图数据推荐最佳时间
		// 这里简化处理，返回下一个整点
		next := now.Add(1 * time.Hour)
		return time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())

	default:
		return now.Add(1 * time.Hour)
	}
}

// IsValidCronExpression 验证Cron表达式是否有效
func (s *ScheduleService) IsValidCronExpression(expr string) bool {
	return s.cronParser.ValidateCron(expr) == nil
}

// isValidScheduleType 验证调度类型
func isValidScheduleType(scheduleType string) bool {
	validTypes := map[string]bool{
		model.ScheduleTypeFixed:    true,
		model.ScheduleTypeInterval: true,
		model.ScheduleTypeEvent:    true,
		model.ScheduleTypeHeat:     true,
	}
	return validTypes[scheduleType]
}