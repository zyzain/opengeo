package service

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// MockHeatmapRepository 模拟热力图仓库
type MockHeatmapRepository struct {
	data map[string]*HeatmapData
}

func NewMockHeatmapRepository() *MockHeatmapRepository {
	return &MockHeatmapRepository{
		data: make(map[string]*HeatmapData),
	}
}

func (r *MockHeatmapRepository) key(platform, aiModel string, hour, dayOfWeek int32) string {
	return fmt.Sprintf("%s-%s-%d-%d", platform, aiModel, hour, dayOfWeek)
}

func (r *MockHeatmapRepository) Get(ctx context.Context, platform, aiModel string, hour, dayOfWeek int32) (*HeatmapData, error) {
	key := r.key(platform, aiModel, hour, dayOfWeek)
	if data, exists := r.data[key]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("not found")
}

func (r *MockHeatmapRepository) Upsert(ctx context.Context, data *HeatmapData) error {
	key := r.key(data.Platform, data.AIModel, data.Hour, data.DayOfWeek)
	r.data[key] = data
	return nil
}

func (r *MockHeatmapRepository) GetByPlatform(ctx context.Context, platform, aiModel string) ([]*HeatmapData, error) {
	result := make([]*HeatmapData, 0)
	for _, d := range r.data {
		if d.Platform == platform && d.AIModel == aiModel {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *MockHeatmapRepository) GetBestSlots(ctx context.Context, platform, aiModel string, limit int) ([]*HeatmapData, error) {
	all, _ := r.GetByPlatform(ctx, platform, aiModel)

	// 按分数排序
	for i := 0; i < len(all)-1; i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].ActivityScore > all[i].ActivityScore {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	if limit > len(all) {
		limit = len(all)
	}
	return all[:limit], nil
}

func TestHeatmapService_RecordActivity(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	// 记录第一次
	err := svc.RecordActivity(ctx, "wechat", "deepseek", 10, 1, 80.0)
	if err != nil {
		t.Fatalf("record failed: %v", err)
	}

	data, _ := repo.Get(ctx, "wechat", "deepseek", 10, 1)
	if data.SampleCount != 1 {
		t.Errorf("expected 1 sample, got %d", data.SampleCount)
	}
	if data.ActivityScore != 80.0 {
		t.Errorf("expected 80.0, got %f", data.ActivityScore)
	}

	// 记录第二次（移动平均）
	err = svc.RecordActivity(ctx, "wechat", "deepseek", 10, 1, 60.0)
	if err != nil {
		t.Fatalf("record failed: %v", err)
	}

	data, _ = repo.Get(ctx, "wechat", "deepseek", 10, 1)
	if data.SampleCount != 2 {
		t.Errorf("expected 2 samples, got %d", data.SampleCount)
	}
	if data.ActivityScore != 70.0 { // (80+60)/2
		t.Errorf("expected 70.0, got %f", data.ActivityScore)
	}
}

func TestHeatmapService_RecordActivity_InvalidInput(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	tests := []struct {
		hour      int32
		dayOfWeek int32
		wantErr   bool
	}{
		{-1, 1, true},
		{24, 1, true},
		{0, 0, true},
		{0, 8, true},
		{10, 1, false},
		{23, 7, false},
	}

	for _, tt := range tests {
		err := svc.RecordActivity(ctx, "test", "test", tt.hour, tt.dayOfWeek, 50.0)
		if (err != nil) != tt.wantErr {
			t.Errorf("hour=%d, day=%d: error=%v, wantErr=%v", tt.hour, tt.dayOfWeek, err, tt.wantErr)
		}
	}
}

func TestHeatmapService_BatchRecordActivity(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	records := []*ActivityRecord{
		{Platform: "wechat", AIModel: "deepseek", Hour: 9, DayOfWeek: 1, Score: 80},
		{Platform: "wechat", AIModel: "deepseek", Hour: 10, DayOfWeek: 1, Score: 90},
		{Platform: "wechat", AIModel: "deepseek", Hour: 14, DayOfWeek: 2, Score: 70},
	}

	err := svc.BatchRecordActivity(ctx, records)
	if err != nil {
		t.Fatalf("batch record failed: %v", err)
	}

	data, _ := repo.GetByPlatform(ctx, "wechat", "deepseek")
	if len(data) != 3 {
		t.Errorf("expected 3 records, got %d", len(data))
	}
}

func TestHeatmapService_GetHeatmap(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	// 添加测试数据
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 9, DayOfWeek: 1, ActivityScore: 85})
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 14, DayOfWeek: 3, ActivityScore: 90})

	heatmap, err := svc.GetHeatmap(ctx, "wechat", "deepseek")
	if err != nil {
		t.Fatalf("get heatmap failed: %v", err)
	}

	if len(heatmap) != 7 {
		t.Errorf("expected 7 rows, got %d", len(heatmap))
	}
	if len(heatmap[0]) != 24 {
		t.Errorf("expected 24 columns, got %d", len(heatmap[0]))
	}

	// 检查数据
	if heatmap[0][9] != 85 { // 周一（索引0）9点
		t.Errorf("expected 85, got %f", heatmap[0][9])
	}
	if heatmap[2][14] != 90 { // 周三（索引2）14点
		t.Errorf("expected 90, got %f", heatmap[2][14])
	}
}

func TestHeatmapService_GetBestTimeSlots(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	// 添加测试数据
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 9, DayOfWeek: 1, ActivityScore: 85})
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 14, DayOfWeek: 1, ActivityScore: 90})
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 20, DayOfWeek: 5, ActivityScore: 75})

	slots, err := svc.GetBestTimeSlots(ctx, "wechat", "deepseek", 2)
	if err != nil {
		t.Fatalf("get best slots failed: %v", err)
	}

	if len(slots) != 2 {
		t.Errorf("expected 2 slots, got %d", len(slots))
	}

	// 第一个应该是最高分
	if slots[0].Score != 90 {
		t.Errorf("expected 90, got %f", slots[0].Score)
	}
	if slots[0].Hour != 14 {
		t.Errorf("expected hour 14, got %d", slots[0].Hour)
	}
}

func TestHeatmapService_GetRecommendedSchedule(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	// 添加测试数据
	for day := int32(1); day <= 7; day++ {
		for hour := int32(0); hour < 24; hour++ {
			score := float32(50 + hour + day*2)
			repo.Upsert(ctx, &HeatmapData{
				Platform:      "wechat",
				AIModel:       "deepseek",
				Hour:          hour,
				DayOfWeek:     day,
				ActivityScore: score,
			})
		}
	}

	recommended, err := svc.GetRecommendedSchedule(ctx, "wechat", "deepseek", 5)
	if err != nil {
		t.Fatalf("get recommended failed: %v", err)
	}

	if len(recommended) != 5 {
		t.Errorf("expected 5 recommendations, got %d", len(recommended))
	}

	// 检查置信度
	for _, rec := range recommended {
		if rec.Confidence < 0 || rec.Confidence > 1 {
			t.Errorf("confidence out of range: %f", rec.Confidence)
		}
	}
}

func TestHeatmapService_GetHeatmapStats(t *testing.T) {
	repo := NewMockHeatmapRepository()
	svc := NewHeatmapService(repo)
	ctx := context.Background()

	// 添加测试数据
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 9, DayOfWeek: 1, ActivityScore: 80, SampleCount: 10})
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 14, DayOfWeek: 1, ActivityScore: 90, SampleCount: 20})
	repo.Upsert(ctx, &HeatmapData{Platform: "wechat", AIModel: "deepseek", Hour: 20, DayOfWeek: 5, ActivityScore: 70, SampleCount: 15})

	stats, err := svc.GetHeatmapStats(ctx, "wechat", "deepseek")
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}

	if stats.DataPoints != 3 {
		t.Errorf("expected 3 data points, got %d", stats.DataPoints)
	}
	if stats.TotalSamples != 45 {
		t.Errorf("expected 45 samples, got %d", stats.TotalSamples)
	}
	if stats.MaxScore != 90 {
		t.Errorf("expected max 90, got %f", stats.MaxScore)
	}
	if stats.MinScore != 70 {
		t.Errorf("expected min 70, got %f", stats.MinScore)
	}
}

func TestGetDayName(t *testing.T) {
	tests := []struct {
		day  int32
		want string
	}{
		{1, "周一"},
		{2, "周二"},
		{3, "周三"},
		{4, "周四"},
		{5, "周五"},
		{6, "周六"},
		{7, "周日"},
		{8, "未知"},
	}

	for _, tt := range tests {
		if got := getDayName(tt.day); got != tt.want {
			t.Errorf("getDayName(%d) = %s, want %s", tt.day, got, tt.want)
		}
	}
}

func TestGetTimeRange(t *testing.T) {
	tests := []struct {
		hour int32
		want string
	}{
		{0, "00:00-00:59"},
		{9, "09:00-09:59"},
		{23, "23:00-23:59"},
	}

	for _, tt := range tests {
		if got := getTimeRange(tt.hour); got != tt.want {
			t.Errorf("getTimeRange(%d) = %s, want %s", tt.hour, got, tt.want)
		}
	}
}

func TestGetRecommendation(t *testing.T) {
	tests := []struct {
		score float32
		want  string
	}{
		{90, "强烈推荐"},
		{80, "强烈推荐"},
		{70, "推荐"},
		{60, "推荐"},
		{50, "一般"},
		{40, "一般"},
		{30, "不推荐"},
	}

	for _, tt := range tests {
		if got := getRecommendation(tt.score); got != tt.want {
			t.Errorf("getRecommendation(%f) = %s, want %s", tt.score, got, tt.want)
		}
	}
}

func TestCalculateConfidence(t *testing.T) {
	tests := []struct {
		score float32
		want  float32
	}{
		{95, 0.95},
		{85, 0.85},
		{75, 0.75},
		{65, 0.65},
		{55, 0.55},
		{45, 0.45},
	}

	for _, tt := range tests {
		got := calculateConfidence(tt.score)
		if got != tt.want {
			t.Errorf("calculateConfidence(%f) = %f, want %f", tt.score, got, tt.want)
		}
	}
}

func TestHeatmapData_Fields(t *testing.T) {
	data := &HeatmapData{
		Platform:      "wechat",
		AIModel:       "deepseek",
		Hour:          10,
		DayOfWeek:     3,
		ActivityScore: 85.5,
		SampleCount:   100,
		UpdatedAt:     time.Now(),
	}

	if data.Platform != "wechat" {
		t.Errorf("expected wechat, got %s", data.Platform)
	}
	if data.Hour != 10 {
		t.Errorf("expected 10, got %d", data.Hour)
	}
}

func TestActivityRecord_Fields(t *testing.T) {
	record := &ActivityRecord{
		Platform:  "wechat",
		AIModel:   "deepseek",
		Hour:      14,
		DayOfWeek: 1,
		Score:     90.0,
	}

	if record.Platform != "wechat" {
		t.Errorf("expected wechat, got %s", record.Platform)
	}
	if record.Score != 90.0 {
		t.Errorf("expected 90.0, got %f", record.Score)
	}
}
