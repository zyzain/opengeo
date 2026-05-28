package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/service/scheduler/internal/domain/model"
)

// ScheduleRepository 调度仓储
type ScheduleRepository struct {
	db *gorm.DB
}

// NewScheduleRepository 创建调度仓储
func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Create 创建调度
func (r *ScheduleRepository) Create(ctx context.Context, schedule *model.Schedule) error {
	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	return nil
}

// GetByID 根据ID获取调度
func (r *ScheduleRepository) GetByID(ctx context.Context, id int64) (*model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.WithContext(ctx).First(&schedule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("schedule not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	return &schedule, nil
}

// Update 更新调度
func (r *ScheduleRepository) Update(ctx context.Context, schedule *model.Schedule) error {
	if err := r.db.WithContext(ctx).Save(schedule).Error; err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	return nil
}

// Delete 删除调度
func (r *ScheduleRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Schedule{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}

// List 列出调度
func (r *ScheduleRepository) List(ctx context.Context, userID int64, isEnabled *bool, page, pageSize int) ([]*model.Schedule, int32, error) {
	var schedules []*model.Schedule
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Schedule{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&schedules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list schedules: %w", err)
	}

	return schedules, int32(total), nil
}

// GetDueSchedules 获取到期的调度
func (r *ScheduleRepository) GetDueSchedules(ctx context.Context, limit int) ([]*model.Schedule, error) {
	var schedules []*model.Schedule

	err := r.db.WithContext(ctx).
		Where("is_enabled = ? AND next_run_time <= ?", true, time.Now()).
		Limit(limit).
		Order("next_run_time ASC").
		Find(&schedules).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get due schedules: %w", err)
	}

	return schedules, nil
}

// UpdateNextRunTime 更新下次运行时间
func (r *ScheduleRepository) UpdateNextRunTime(ctx context.Context, scheduleID int64, nextRunTime time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Schedule{}).
		Where("id = ?", scheduleID).
		Updates(map[string]interface{}{
			"next_run_time": nextRunTime,
			"last_run_time": time.Now(),
			"run_count":     gorm.Expr("run_count + 1"),
		}).Error; err != nil {
		return fmt.Errorf("failed to update next run time: %w", err)
	}
	return nil
}

// ScheduleTaskRepository 调度任务仓储
type ScheduleTaskRepository struct {
	db *gorm.DB
}

// NewScheduleTaskRepository 创建调度任务仓储
func NewScheduleTaskRepository(db *gorm.DB) *ScheduleTaskRepository {
	return &ScheduleTaskRepository{db: db}
}

// Create 创建调度任务
func (r *ScheduleTaskRepository) Create(ctx context.Context, task *model.ScheduleTask) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return fmt.Errorf("failed to create schedule task: %w", err)
	}
	return nil
}

// GetPendingTasks 获取待执行任务
func (r *ScheduleTaskRepository) GetPendingTasks(ctx context.Context, scheduleID int64, limit int) ([]*model.ScheduleTask, error) {
	var tasks []*model.ScheduleTask

	err := r.db.WithContext(ctx).
		Where("schedule_id = ? AND status = ? AND scheduled_at <= ?", scheduleID, "pending", time.Now()).
		Limit(limit).
		Order("scheduled_at ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending tasks: %w", err)
	}

	return tasks, nil
}

// UpdateStatus 更新任务状态
func (r *ScheduleTaskRepository) UpdateStatus(ctx context.Context, scheduleID, taskID int64, status string) error {
	if err := r.db.WithContext(ctx).
		Model(&model.ScheduleTask{}).
		Where("schedule_id = ? AND task_id = ?", scheduleID, taskID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

// AIActivityHeatmapRepository 热力图仓储
type AIActivityHeatmapRepository struct {
	db *gorm.DB
}

// NewAIActivityHeatmapRepository 创建热力图仓储
func NewAIActivityHeatmapRepository(db *gorm.DB) *AIActivityHeatmapRepository {
	return &AIActivityHeatmapRepository{db: db}
}

// Get 获取热力图数据
func (r *AIActivityHeatmapRepository) Get(ctx context.Context, platform, aiModel string, timeSlot, dayOfWeek int32) (*model.AIActivityHeatmap, error) {
	var heatmap model.AIActivityHeatmap

	err := r.db.WithContext(ctx).
		Where("platform = ? AND ai_model = ? AND time_slot = ? AND day_of_week = ?", platform, aiModel, timeSlot, dayOfWeek).
		First(&heatmap).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.AIActivityHeatmap{
				Platform:      platform,
				AIModel:       aiModel,
				TimeSlot:      timeSlot,
				DayOfWeek:     dayOfWeek,
				ActivityScore: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get heatmap: %w", err)
	}

	return &heatmap, nil
}

// Upsert 更新或插入热力图数据
func (r *AIActivityHeatmapRepository) Upsert(ctx context.Context, heatmap *model.AIActivityHeatmap) error {
	if err := r.db.WithContext(ctx).
		Where("platform = ? AND ai_model = ? AND time_slot = ? AND day_of_week = ?", heatmap.Platform, heatmap.AIModel, heatmap.TimeSlot, heatmap.DayOfWeek).
		Assign(map[string]interface{}{
			"activity_score": heatmap.ActivityScore,
			"sample_count":   heatmap.SampleCount,
			"updated_at":     time.Now(),
		}).
		FirstOrCreate(heatmap).Error; err != nil {
		return fmt.Errorf("failed to upsert heatmap: %w", err)
	}
	return nil
}

// GetBestTimeSlots 获取最佳时间段
func (r *AIActivityHeatmapRepository) GetBestTimeSlots(ctx context.Context, platform, aiModel string, limit int) ([]*model.AIActivityHeatmap, error) {
	var heatmaps []*model.AIActivityHeatmap

	err := r.db.WithContext(ctx).
		Where("platform = ? AND ai_model = ?", platform, aiModel).
		Order("activity_score DESC").
		Limit(limit).
		Find(&heatmaps).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get best time slots: %w", err)
	}

	return heatmaps, nil
}

// PublishCalendarRepository 发布日历仓储
type PublishCalendarRepository struct {
	db *gorm.DB
}

// NewPublishCalendarRepository 创建发布日历仓储
func NewPublishCalendarRepository(db *gorm.DB) *PublishCalendarRepository {
	return &PublishCalendarRepository{db: db}
}

// Create 创建日历事件
func (r *PublishCalendarRepository) Create(ctx context.Context, event *model.PublishCalendar) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return fmt.Errorf("failed to create calendar event: %w", err)
	}
	return nil
}

// GetByID 根据ID获取日历事件
func (r *PublishCalendarRepository) GetByID(ctx context.Context, id int64) (*model.PublishCalendar, error) {
	var event model.PublishCalendar
	if err := r.db.WithContext(ctx).First(&event, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("calendar event not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get calendar event: %w", err)
	}
	return &event, nil
}

// Update 更新日历事件
func (r *PublishCalendarRepository) Update(ctx context.Context, event *model.PublishCalendar) error {
	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		return fmt.Errorf("failed to update calendar event: %w", err)
	}
	return nil
}

// Delete 删除日历事件
func (r *PublishCalendarRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.PublishCalendar{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}
	return nil
}

// List 列出日历事件
func (r *PublishCalendarRepository) List(ctx context.Context, userID int64, startTime, endTime time.Time) ([]*model.PublishCalendar, error) {
	var events []*model.PublishCalendar

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if !startTime.IsZero() {
		query = query.Where("start_time >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("start_time <= ?", endTime)
	}

	if err := query.Order("start_time ASC").Limit(500).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to list calendar events: %w", err)
	}

	return events, nil
}

// ScheduleLogRepository 调度日志仓储
type ScheduleLogRepository struct {
	db *gorm.DB
}

// NewScheduleLogRepository 创建调度日志仓储
func NewScheduleLogRepository(db *gorm.DB) *ScheduleLogRepository {
	return &ScheduleLogRepository{db: db}
}

// Create 创建日志
func (r *ScheduleLogRepository) Create(ctx context.Context, log *model.ScheduleLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create schedule log: %w", err)
	}
	return nil
}

// List 列出日志
func (r *ScheduleLogRepository) List(ctx context.Context, scheduleID int64, page, pageSize int) ([]*model.ScheduleLog, int32, error) {
	var logs []*model.ScheduleLog
	var total int64

	query := r.db.WithContext(ctx).Where("schedule_id = ?", scheduleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count logs: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list logs: %w", err)
	}

	return logs, int32(total), nil
}