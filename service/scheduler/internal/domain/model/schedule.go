package model

import "time"

// Schedule 调度任务
type Schedule struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64      `json:"user_id" gorm:"index;not null"`
	ScheduleName   string     `json:"schedule_name" gorm:"size:128;not null"`
	ScheduleType   string     `json:"schedule_type" gorm:"size:32;not null"`
	CronExpression string     `json:"cron_expression" gorm:"size:128"`
	Config         string     `json:"config" gorm:"type:text"`
	IsEnabled      bool       `json:"is_enabled" gorm:"default:true;index"`
	NextRunTime    *time.Time `json:"next_run_time" gorm:"index"`
	LastRunTime    *time.Time `json:"last_run_time"`
	RunCount       int64      `json:"run_count" gorm:"default:0"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ScheduleType 调度类型
const (
	ScheduleTypeFixed   = "fixed"   // 固定时间
	ScheduleTypeInterval = "interval" // 间隔循环
	ScheduleTypeEvent   = "event"   // 事件触发
	ScheduleTypeHeat    = "heat"    // 热力图推荐
)

// ScheduleTask 调度任务关联
type ScheduleTask struct {
	ScheduleID  int64     `json:"schedule_id" gorm:"primaryKey;index:idx_schedule_status_sched"`
	TaskID      int64     `json:"task_id" gorm:"primaryKey"`
	ScheduledAt time.Time `json:"scheduled_at" gorm:"index:idx_schedule_status_sched"`
	Status      string    `json:"status" gorm:"size:32;index:idx_schedule_status_sched"` // pending, executed, failed
	CreatedAt   time.Time `json:"created_at"`
}

// AIActivityHeatmap AI活跃度热力图
type AIActivityHeatmap struct {
	ID            int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	Platform      string  `json:"platform" gorm:"size:64;not null;index:idx_heatmap_lookup"`
	AIModel       string  `json:"ai_model" gorm:"size:64;not null;index:idx_heatmap_lookup"`
	TimeSlot      int32   `json:"time_slot" gorm:"not null;index:idx_heatmap_lookup"` // 0-23小时
	DayOfWeek     int32   `json:"day_of_week" gorm:"not null;index:idx_heatmap_lookup"` // 1-7星期
	ActivityScore float32 `json:"activity_score"`
	SampleCount   int64   `json:"sample_count"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PublishCalendar 发布日历
type PublishCalendar struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64      `json:"user_id" gorm:"index;not null"`
	Title       string     `json:"title" gorm:"size:256;not null"`
	Description string     `json:"description" gorm:"type:text"`
	StartTime   time.Time  `json:"start_time" gorm:"not null;index"`
	EndTime     *time.Time `json:"end_time"`
	ScheduleID  *int64     `json:"schedule_id" gorm:"index"`
	ContentID   *int64     `json:"content_id"`
	ChannelID   *int64     `json:"channel_id"`
	Status      string     `json:"status" gorm:"size:32;default:pending"` // pending, executing, completed, cancelled
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ScheduleLog 调度日志
type ScheduleLog struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID int64     `json:"schedule_id" gorm:"index;not null"`
	TaskID     *int64    `json:"task_id"`
	Action     string    `json:"action" gorm:"size:32"` // trigger, skip, error
	Message    string    `json:"message" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
}

// Priority 任务优先级
type Priority int32

const (
	PriorityLow    Priority = 0
	PriorityMedium Priority = 1
	PriorityHigh   Priority = 2
	PriorityUrgent Priority = 3
)