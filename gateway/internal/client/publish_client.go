package client

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strconv"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"

	"opengeo/pkg/similarity"
)

type PublishClient struct {
	db *gorm.DB
}

func NewPublishClient(db *gorm.DB) *PublishClient {
	return &PublishClient{db: db}
}

// 渠道
func (c *PublishClient) CreateChannel(ctx context.Context, userID int64, channelType, channelName, channelConfig string) (map[string]interface{}, error) {
	channel := map[string]interface{}{
		"user_id":        userID,
		"channel_type":   channelType,
		"channel_name":   channelName,
		"channel_config": channelConfig,
		"is_enabled":     true,
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}
	if err := c.db.WithContext(ctx).Table("publish_channels").Create(channel).Error; err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}
	return channel, nil
}

func (c *PublishClient) GetChannel(ctx context.Context, id int64) (map[string]interface{}, error) {
	var channel map[string]interface{}
	if err := c.db.WithContext(ctx).Table("publish_channels").Where("id = ?", id).First(&channel).Error; err != nil {
		return nil, fmt.Errorf("channel not found")
	}
	return channel, nil
}

func (c *PublishClient) ListChannels(ctx context.Context, userID int64, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("publish_channels")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *PublishClient) GetChannelPlatforms(ctx context.Context) ([]map[string]interface{}, error) {
	var platforms []map[string]interface{}
	c.db.WithContext(ctx).Table("publish_channels").Select("DISTINCT channel_type").Find(&platforms)
	return platforms, nil
}

// 发布任务
func (c *PublishClient) CreatePublishTask(ctx context.Context, userID, contentID, channelID int64, scheduledTime string) (map[string]interface{}, error) {
	task := map[string]interface{}{
		"user_id":     userID,
		"content_id":  contentID,
		"channel_id":  channelID,
		"status":      0,
		"max_retries": 3,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	// 空字符串转换为nil，避免MySQL datetime错误
	if scheduledTime != "" {
		task["scheduled_time"] = scheduledTime
	}
	if err := c.db.WithContext(ctx).Table("publish_tasks").Create(task).Error; err != nil {
		return nil, fmt.Errorf("create publish task: %w", err)
	}
	return task, nil
}

func (c *PublishClient) GetPublishTask(ctx context.Context, id int64) (map[string]interface{}, error) {
	var task map[string]interface{}
	if err := c.db.WithContext(ctx).Table("publish_tasks").Where("id = ?", id).First(&task).Error; err != nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (c *PublishClient) ListPublishTasks(ctx context.Context, userID int64, status, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("publish_tasks")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *PublishClient) CancelPublishTask(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("publish_tasks").Where("id = ?", id).Update("status", 4).Error
}

func (c *PublishClient) RetryPublishTask(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("publish_tasks").Where("id = ?", id).Update("status", 0).Error
}

// 平台
func (c *PublishClient) ListPlatforms(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	c.db.WithContext(ctx).Table("platforms").Count(&total)
	c.db.WithContext(ctx).Table("platforms").Offset((page-1)*pageSize).Limit(pageSize).Order("sort_order DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *PublishClient) GetPlatform(ctx context.Context, id int64) (map[string]interface{}, error) {
	var platform map[string]interface{}
	if err := c.db.WithContext(ctx).Table("platforms").Where("id = ?", id).First(&platform).Error; err != nil {
		return nil, fmt.Errorf("platform not found")
	}
	return platform, nil
}

func (c *PublishClient) CreatePlatform(ctx context.Context, req *CreatePlatformRequest) (map[string]interface{}, error) {
	platform := map[string]interface{}{
		"code":        req.Code,
		"name":        req.Name,
		"icon":        req.Icon,
		"color":       req.Color,
		"description": req.Description,
		"is_enabled":  true,
		"sort_order":  req.SortOrder,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	if err := c.db.WithContext(ctx).Table("platforms").Create(platform).Error; err != nil {
		return nil, fmt.Errorf("create platform: %w", err)
	}
	return platform, nil
}

func (c *PublishClient) UpdatePlatform(ctx context.Context, id int64, req *UpdatePlatformRequest) (map[string]interface{}, error) {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if err := c.db.WithContext(ctx).Table("platforms").Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update platform: %w", err)
	}
	var platform map[string]interface{}
	c.db.WithContext(ctx).Table("platforms").Where("id = ?", id).First(&platform)
	return platform, nil
}

func (c *PublishClient) DeletePlatform(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("platforms").Where("id = ?", id).Delete(nil).Error
}

func (c *PublishClient) EnablePlatform(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("platforms").Where("id = ?", id).Update("is_enabled", true).Error
}

func (c *PublishClient) DisablePlatform(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("platforms").Where("id = ?", id).Update("is_enabled", false).Error
}

func (c *PublishClient) PreviewPublish(ctx context.Context, channelID int64, title, body string) (map[string]interface{}, error) {
	// 生成预览HTML（转义用户输入防止XSS）
	safeTitle := html.EscapeString(title)
	safeBody := html.EscapeString(body)
	previewHTML := fmt.Sprintf("<div class='preview'><h1>%s</h1><div>%s</div></div>", safeTitle, safeBody)
	return map[string]interface{}{
		"success":      true,
		"preview_html": previewHTML,
		"preview_url":  "",
	}, nil
}

func (c *PublishClient) ValidatePublish(ctx context.Context, channelID int64, title, body string) (map[string]interface{}, error) {
	errors := make([]string, 0)
	warnings := make([]string, 0)

	if title == "" {
		errors = append(errors, "标题不能为空")
	}
	if body == "" {
		errors = append(errors, "正文不能为空")
	}
	if len(title) > 64 {
		warnings = append(warnings, "标题超过64字符可能被截断")
	}

	return map[string]interface{}{
		"valid":    len(errors) == 0,
		"errors":   errors,
		"warnings": warnings,
	}, nil
}

// ==================== 内容去重相关 ====================

// CheckContentDedup 检查内容去重
func (c *PublishClient) CheckContentDedup(ctx context.Context, userID int64, text string) (map[string]interface{}, error) {
	// 计算内容指纹
	simCalc := similarity.NewCombinedSimilarity()
	fingerprint := simCalc.ComputeFingerprint(text)

	// 查找相似内容
	var similarContents []map[string]interface{}

	// 查询内容指纹表
	var fingerprints []map[string]interface{}
	query := c.db.WithContext(ctx).Table("content_fingerprints")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	query.Order("created_at DESC").Limit(100).Find(&fingerprints)

	// 计算相似度
	similarityThreshold := 0.7
	for _, fp := range fingerprints {
		if bodyHash, ok := fp["body_hash"].(string); ok && bodyHash != "" {
			// 使用SimHash计算相似度
			simHashCalc := similarity.NewSimHash()
			// 将字符串hash转换为uint64
			hashValue, _ := strconv.ParseUint(bodyHash, 10, 64)
			sim := simHashCalc.Similarity(hashValue, fingerprint.SimHash)
			if sim >= similarityThreshold {
				similarContents = append(similarContents, map[string]interface{}{
					"id":         fp["content_id"],
					"title":      fp["title"],
					"similarity": sim,
					"source":     "内容库",
					"status":     getSimilarityStatus(sim),
				})
			}
		}
	}

	// 如果没有找到相似内容，查询历史发布内容
	if len(similarContents) == 0 {
		var publishTasks []map[string]interface{}
		c.db.WithContext(ctx).Table("publish_tasks").
			Where("user_id = ? AND status = ?", userID, 2). // status=2 表示已发布
			Order("created_at DESC").
			Limit(50).
			Find(&publishTasks)

		for _, task := range publishTasks {
			if contentID, ok := task["content_id"].(int64); ok && contentID > 0 {
				// 查询内容详情
				var content map[string]interface{}
				if err := c.db.WithContext(ctx).Table("contents").Where("id = ?", contentID).First(&content).Error; err == nil {
					if body, ok := content["body"].(string); ok {
						// 计算余弦相似度
						cosineSim := similarity.NewCosineSimilarity()
						sim := cosineSim.Compute(text, body)
						if sim >= similarityThreshold {
							similarContents = append(similarContents, map[string]interface{}{
								"id":         contentID,
								"title":      content["title"],
								"similarity": sim,
								"source":     "已发布内容",
								"status":     getSimilarityStatus(sim),
							})
						}
					}
				}
			}
		}
	}

	// 生成建议
	suggestions := generateDedupSuggestions(text, similarContents)

	return map[string]interface{}{
		"text_length":      utf8.RuneCountInString(text),
		"similarity_score": calculateOverallSimilarity(similarContents),
		"duplicates":       similarContents,
		"suggestions":      suggestions,
		"fingerprint":      fingerprint,
	}, nil
}

// SaveContentFingerprint 保存内容指纹
func (c *PublishClient) SaveContentFingerprint(ctx context.Context, userID, contentID int64, title, body, contentType string) error {
	simCalc := similarity.NewCombinedSimilarity()
	fingerprint := simCalc.ComputeFingerprint(body)
	titleFingerprint := simCalc.ComputeFingerprint(title)

	keywordsJSON, _ := json.Marshal(fingerprint.Keywords)

	fp := map[string]interface{}{
		"user_id":           userID,
		"content_id":        contentID,
		"title_hash":        fmt.Sprintf("%d", titleFingerprint.SimHash),
		"body_hash":         fmt.Sprintf("%d", fingerprint.SimHash),
		"title_fingerprint": fmt.Sprintf("%d", titleFingerprint.SimHash),
		"body_fingerprint":  fmt.Sprintf("%d", fingerprint.SimHash),
		"keywords":          string(keywordsJSON),
		"word_count":        fingerprint.WordCount,
		"content_type":      contentType,
		"created_at":        time.Now(),
	}

	return c.db.WithContext(ctx).Table("content_fingerprints").Create(fp).Error
}

// ==================== 请求模型 ====================

type CreatePlatformRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sort_order"`
}

type UpdatePlatformRequest struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sort_order"`
}

// 辅助函数
func getSimilarityStatus(sim float64) string {
	if sim >= 0.9 {
		return "high"
	}
	if sim >= 0.7 {
		return "medium"
	}
	return "low"
}

func calculateOverallSimilarity(duplicates []map[string]interface{}) float64 {
	if len(duplicates) == 0 {
		return 0
	}

	maxSim := 0.0
	for _, dup := range duplicates {
		if sim, ok := dup["similarity"].(float64); ok && sim > maxSim {
			maxSim = sim
		}
	}
	return maxSim
}

func generateDedupSuggestions(text string, duplicates []map[string]interface{}) []string {
	suggestions := make([]string, 0)

	if len(duplicates) > 0 {
		suggestions = append(suggestions, "发现相似内容，建议添加独特的观点和分析")
		suggestions = append(suggestions, "可使用同义词替换部分高频词汇")
		suggestions = append(suggestions, "建议调整段落顺序或重新组织内容结构")
	} else {
		suggestions = append(suggestions, "内容原创度较高，建议保持")
	}

	// 基于内容长度的建议
	textLen := utf8.RuneCountInString(text)
	if textLen < 100 {
		suggestions = append(suggestions, "内容较短，建议补充更多细节")
	}

	return suggestions
}
