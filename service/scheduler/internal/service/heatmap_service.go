package service

import (
	"context"
	"fmt"
	"math"
	"time"
)

// HeatmapService 热力图服务
type HeatmapService struct {
	// 数据仓库接口
	repo HeatmapRepository
}

// HeatmapRepository 热力图数据仓库接口
type HeatmapRepository interface {
	Get(ctx context.Context, platform, aiModel string, hour, dayOfWeek int32) (*HeatmapData, error)
	Upsert(ctx context.Context, data *HeatmapData) error
	GetByPlatform(ctx context.Context, platform, aiModel string) ([]*HeatmapData, error)
	GetBestSlots(ctx context.Context, platform, aiModel string, limit int) ([]*HeatmapData, error)
}

// HeatmapData 热力图数据
type HeatmapData struct {
	Platform      string    `json:"platform"`
	AIModel       string    `json:"ai_model"`
	Hour          int32     `json:"hour"`         // 0-23
	DayOfWeek     int32     `json:"day_of_week"`  // 1=周一, 7=周日
	ActivityScore float32   `json:"activity_score"`
	SampleCount   int64     `json:"sample_count"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewHeatmapService 创建热力图服务
func NewHeatmapService(repo HeatmapRepository) *HeatmapService {
	return &HeatmapService{repo: repo}
}

// ==================== 数据采集 ====================

// RecordActivity 记录活跃度数据
func (s *HeatmapService) RecordActivity(ctx context.Context, platform, aiModel string, hour, dayOfWeek int32, score float32) error {
	if hour < 0 || hour > 23 {
		return fmt.Errorf("invalid hour: %d", hour)
	}
	if dayOfWeek < 1 || dayOfWeek > 7 {
		return fmt.Errorf("invalid day of week: %d", dayOfWeek)
	}

	// 获取现有数据
	existing, err := s.repo.Get(ctx, platform, aiModel, hour, dayOfWeek)
	if err != nil {
		// 创建新记录
		return s.repo.Upsert(ctx, &HeatmapData{
			Platform:      platform,
			AIModel:       aiModel,
			Hour:          hour,
			DayOfWeek:     dayOfWeek,
			ActivityScore: score,
			SampleCount:   1,
			UpdatedAt:     time.Now(),
		})
	}

	// 更新现有记录（移动平均）
	newCount := existing.SampleCount + 1
	newScore := (existing.ActivityScore*float32(existing.SampleCount) + score) / float32(newCount)

	return s.repo.Upsert(ctx, &HeatmapData{
		Platform:      platform,
		AIModel:       aiModel,
		Hour:          hour,
		DayOfWeek:     dayOfWeek,
		ActivityScore: newScore,
		SampleCount:   newCount,
		UpdatedAt:     time.Now(),
	})
}

// BatchRecordActivity 批量记录活跃度
func (s *HeatmapService) BatchRecordActivity(ctx context.Context, records []*ActivityRecord) error {
	for _, record := range records {
		if err := s.RecordActivity(ctx, record.Platform, record.AIModel, record.Hour, record.DayOfWeek, record.Score); err != nil {
			return fmt.Errorf("record activity: %w", err)
		}
	}
	return nil
}

// ActivityRecord 活跃度记录
type ActivityRecord struct {
	Platform  string  `json:"platform"`
	AIModel   string  `json:"ai_model"`
	Hour      int32   `json:"hour"`
	DayOfWeek int32   `json:"day_of_week"`
	Score     float32 `json:"score"`
}

// ==================== 数据分析 ====================

// GetHeatmap 获取热力图数据
func (s *HeatmapService) GetHeatmap(ctx context.Context, platform, aiModel string) ([][]float32, error) {
	data, err := s.repo.GetByPlatform(ctx, platform, aiModel)
	if err != nil {
		return nil, fmt.Errorf("get heatmap: %w", err)
	}

	// 构建 7x24 矩阵
	heatmap := make([][]float32, 7)
	for i := range heatmap {
		heatmap[i] = make([]float32, 24)
	}

	for _, d := range data {
		dayIdx := d.DayOfWeek - 1 // 转换为0-6索引
		if dayIdx >= 0 && dayIdx < 7 && d.Hour >= 0 && d.Hour < 24 {
			heatmap[dayIdx][d.Hour] = d.ActivityScore
		}
	}

	return heatmap, nil
}

// GetBestTimeSlots 获取最佳时间段
func (s *HeatmapService) GetBestTimeSlots(ctx context.Context, platform, aiModel string, limit int) ([]*TimeSlot, error) {
	slots, err := s.repo.GetBestSlots(ctx, platform, aiModel, limit)
	if err != nil {
		return nil, fmt.Errorf("get best slots: %w", err)
	}

	result := make([]*TimeSlot, len(slots))
	for i, slot := range slots {
		result[i] = &TimeSlot{
			DayOfWeek:     slot.DayOfWeek,
			Hour:          slot.Hour,
			Score:         slot.ActivityScore,
			DayName:       getDayName(slot.DayOfWeek),
			TimeRange:     getTimeRange(slot.Hour),
			Recommendation: getRecommendation(slot.ActivityScore),
		}
	}

	return result, nil
}

// TimeSlot 时间段
type TimeSlot struct {
	DayOfWeek      int32   `json:"day_of_week"`
	Hour           int32   `json:"hour"`
	Score          float32 `json:"score"`
	DayName        string  `json:"day_name"`
	TimeRange      string  `json:"time_range"`
	Recommendation string  `json:"recommendation"`
}

// GetRecommendedSchedule 获取推荐调度时间
func (s *HeatmapService) GetRecommendedSchedule(ctx context.Context, platform, aiModel string, count int) ([]*RecommendedTime, error) {
	slots, err := s.GetBestTimeSlots(ctx, platform, aiModel, count*2) // 获取更多候选
	if err != nil {
		return nil, err
	}

	// 去重和筛选
	result := make([]*RecommendedTime, 0, count)
	seen := make(map[string]bool)

	for _, slot := range slots {
		key := fmt.Sprintf("%d-%d", slot.DayOfWeek, slot.Hour)
		if seen[key] {
			continue
		}
		seen[key] = true

		result = append(result, &RecommendedTime{
			DayOfWeek:  slot.DayOfWeek,
			Hour:       slot.Hour,
			Score:      slot.Score,
			Confidence: calculateConfidence(slot.Score),
		})

		if len(result) >= count {
			break
		}
	}

	return result, nil
}

// RecommendedTime 推荐时间
type RecommendedTime struct {
	DayOfWeek  int32   `json:"day_of_week"`
	Hour       int32   `json:"hour"`
	Score      float32 `json:"score"`
	Confidence float32 `json:"confidence"` // 0-1
}

// ==================== 统计分析 ====================

// GetHeatmapStats 获取热力图统计
func (s *HeatmapService) GetHeatmapStats(ctx context.Context, platform, aiModel string) (*HeatmapStats, error) {
	data, err := s.repo.GetByPlatform(ctx, platform, aiModel)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &HeatmapStats{}, nil
	}

	var totalScore float32
	var maxScore float32
	var minScore float32 = math.MaxFloat32
	var totalSamples int64

	for _, d := range data {
		totalScore += d.ActivityScore
		totalSamples += d.SampleCount
		if d.ActivityScore > maxScore {
			maxScore = d.ActivityScore
		}
		if d.ActivityScore < minScore {
			minScore = d.ActivityScore
		}
	}

	avgScore := totalScore / float32(len(data))

	// 计算标准差
	var variance float32
	for _, d := range data {
		diff := d.ActivityScore - avgScore
		variance += diff * diff
	}
	variance /= float32(len(data))
	stdDev := float32(math.Sqrt(float64(variance)))

	return &HeatmapStats{
		Platform:     platform,
		AIModel:      aiModel,
		DataPoints:   len(data),
		TotalSamples: totalSamples,
		AvgScore:     avgScore,
		MaxScore:     maxScore,
		MinScore:     minScore,
		StdDev:       stdDev,
		LastUpdated:  time.Now(),
	}, nil
}

// HeatmapStats 热力图统计
type HeatmapStats struct {
	Platform     string    `json:"platform"`
	AIModel      string    `json:"ai_model"`
	DataPoints   int       `json:"data_points"`
	TotalSamples int64     `json:"total_samples"`
	AvgScore     float32   `json:"avg_score"`
	MaxScore     float32   `json:"max_score"`
	MinScore     float32   `json:"min_score"`
	StdDev       float32   `json:"std_dev"`
	LastUpdated  time.Time `json:"last_updated"`
}

// ==================== 辅助函数 ====================

func getDayName(dayOfWeek int32) string {
	names := map[int32]string{
		1: "周一", 2: "周二", 3: "周三", 4: "周四",
		5: "周五", 6: "周六", 7: "周日",
	}
	if name, ok := names[dayOfWeek]; ok {
		return name
	}
	return "未知"
}

func getTimeRange(hour int32) string {
	return fmt.Sprintf("%02d:00-%02d:59", hour, hour)
}

func getRecommendation(score float32) string {
	if score >= 80 {
		return "强烈推荐"
	} else if score >= 60 {
		return "推荐"
	} else if score >= 40 {
		return "一般"
	} else {
		return "不推荐"
	}
}

func calculateConfidence(score float32) float32 {
	// 将分数转换为置信度 (0-1)
	if score >= 90 {
		return 0.95
	} else if score >= 80 {
		return 0.85
	} else if score >= 70 {
		return 0.75
	} else if score >= 60 {
		return 0.65
	} else if score >= 50 {
		return 0.55
	} else {
		return 0.45
	}
}
