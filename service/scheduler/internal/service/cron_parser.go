package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronParser Cron表达式解析器
type CronParser struct{}

// NewCronParser 创建Cron解析器
func NewCronParser() *CronParser {
	return &CronParser{}
}

// Parse 解析Cron表达式，返回下次运行时间
func (p *CronParser) Parse(cronExpr string, from time.Time) (time.Time, error) {
	cronExpr = strings.TrimSpace(cronExpr)

	// 处理特殊表达式
	if strings.HasPrefix(cronExpr, "@") {
		return p.parseSpecialExpression(cronExpr, from)
	}

	// 解析标准5字段Cron
	return p.parseStandardCron(cronExpr, from)
}

// ParseInterval 解析间隔表达式 (如 "30s", "5m", "1h")
func (p *CronParser) ParseInterval(interval string) (time.Duration, error) {
	interval = strings.TrimSpace(interval)

	if len(interval) < 2 {
		return 0, fmt.Errorf("invalid interval: %s", interval)
	}

	unit := interval[len(interval)-1]
	valueStr := interval[:len(interval)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid interval value: %s", valueStr)
	}

	switch unit {
	case 's':
		return time.Duration(value) * time.Second, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	case 'd':
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid interval unit: %c", unit)
	}
}

// ==================== 特殊表达式 ====================

func (p *CronParser) parseSpecialExpression(expr string, from time.Time) (time.Time, error) {
	switch {
	case expr == "@yearly" || expr == "@annually":
		return p.nextYear(from), nil
	case expr == "@monthly":
		return p.nextMonth(from), nil
	case expr == "@weekly":
		return p.nextWeek(from), nil
	case expr == "@daily" || expr == "@midnight":
		return p.nextDay(from), nil
	case expr == "@hourly":
		return from.Add(1 * time.Hour), nil
	case strings.HasPrefix(expr, "@every "):
		durationStr := strings.TrimPrefix(expr, "@every ")
		duration, err := p.ParseInterval(durationStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid @every expression: %w", err)
		}
		return from.Add(duration), nil
	default:
		return time.Time{}, fmt.Errorf("unknown special expression: %s", expr)
	}
}

// ==================== 标准Cron解析 ====================

func (p *CronParser) parseStandardCron(expr string, from time.Time) (time.Time, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("cron expression must have 5 fields, got %d", len(fields))
	}

	minute, err := p.parseField(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute field: %w", err)
	}

	hour, err := p.parseField(fields[1], 0, 23)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour field: %w", err)
	}

	dayOfMonth, err := p.parseField(fields[2], 1, 31)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day field: %w", err)
	}

	month, err := p.parseField(fields[3], 1, 12)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month field: %w", err)
	}

	dayOfWeek, err := p.parseField(fields[4], 0, 7)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid weekday field: %w", err)
	}

	// 将星期日(7)统一为0
	for i, v := range dayOfWeek {
		if v == 7 {
			dayOfWeek[i] = 0
		}
	}

	return p.findNextTime(from, minute, hour, dayOfMonth, month, dayOfWeek)
}

// parseField 解析单个字段
func (p *CronParser) parseField(field string, min, max int) ([]int, error) {
	// 处理通配符
	if field == "*" {
		result := make([]int, max-min+1)
		for i := min; i <= max; i++ {
			result[i-min] = i
		}
		return result, nil
	}

	// 处理步进 (如 */5, 1-10/2)
	if strings.Contains(field, "/") {
		parts := strings.Split(field, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid step expression: %s", field)
		}

		step, err := strconv.Atoi(parts[1])
		if err != nil || step <= 0 {
			return nil, fmt.Errorf("invalid step value: %s", parts[1])
		}

		var base []int
		if parts[0] == "*" {
			base = make([]int, max-min+1)
			for i := min; i <= max; i++ {
				base[i-min] = i
			}
		} else {
			base, err = p.parseField(parts[0], min, max)
			if err != nil {
				return nil, err
			}
		}

		result := make([]int, 0)
		for i := 0; i < len(base); i += step {
			result = append(result, base[i])
		}
		return result, nil
	}

	// 处理范围 (如 1-5)
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range expression: %s", field)
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid range start: %s", parts[0])
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid range end: %s", parts[1])
		}

		if start < min || end > max || start > end {
			return nil, fmt.Errorf("invalid range: %d-%d", start, end)
		}

		result := make([]int, end-start+1)
		for i := start; i <= end; i++ {
			result[i-start] = i
		}
		return result, nil
	}

	// 处理列表 (如 1,3,5)
	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		result := make([]int, len(parts))
		for i, part := range parts {
			val, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("invalid list value: %s", part)
			}
			if val < min || val > max {
				return nil, fmt.Errorf("value out of range: %d", val)
			}
			result[i] = val
		}
		return result, nil
	}

	// 处理单个值
	val, err := strconv.Atoi(field)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %s", field)
	}
	if val < min || val > max {
		return nil, fmt.Errorf("value out of range: %d", val)
	}
	return []int{val}, nil
}

// findNextTime 查找下一个匹配时间
func (p *CronParser) findNextTime(from time.Time, minutes, hours, days, months, weekdays []int) (time.Time, error) {
	// 从当前时间开始，最多搜索2年
	for t := from.Add(1 * time.Minute); t.Before(from.AddDate(2, 0, 0)); t = t.Add(1 * time.Minute) {
		if p.matchesTime(t, minutes, hours, days, months, weekdays) {
			// 将秒归零
			return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location()), nil
		}
	}

	return time.Time{}, fmt.Errorf("no matching time found within 2 years")
}

// matchesTime 检查时间是否匹配
func (p *CronParser) matchesTime(t time.Time, minutes, hours, days, months, weekdays []int) bool {
	// 检查月份
	if !containsInt(months, int(t.Month())) {
		return false
	}

	// 检查日期
	if !containsInt(days, t.Day()) {
		return false
	}

	// 检查星期 (0=周日, 1=周一, ..., 6=周六)
	weekday := int(t.Weekday())
	if !containsInt(weekdays, weekday) {
		return false
	}

	// 检查小时
	if !containsInt(hours, t.Hour()) {
		return false
	}

	// 检查分钟
	if !containsInt(minutes, t.Minute()) {
		return false
	}

	return true
}

// ==================== 辅助方法 ====================

func (p *CronParser) nextYear(from time.Time) time.Time {
	return time.Date(from.Year()+1, 1, 1, 0, 0, 0, 0, from.Location())
}

func (p *CronParser) nextMonth(from time.Time) time.Time {
	year, month, _ := from.Date()
	if month == 12 {
		return time.Date(year+1, 1, 1, 0, 0, 0, 0, from.Location())
	}
	return time.Date(year, month+1, 1, 0, 0, 0, 0, from.Location())
}

func (p *CronParser) nextWeek(from time.Time) time.Time {
	daysUntilSunday := (7 - int(from.Weekday())) % 7
	if daysUntilSunday == 0 {
		daysUntilSunday = 7
	}
	next := from.AddDate(0, 0, daysUntilSunday)
	return time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
}

func (p *CronParser) nextDay(from time.Time) time.Time {
	next := from.AddDate(0, 0, 1)
	return time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
}

func containsInt(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// ValidateCron 验证Cron表达式
func (p *CronParser) ValidateCron(expr string) error {
	expr = strings.TrimSpace(expr)

	// 特殊表达式
	if strings.HasPrefix(expr, "@") {
		validSpecial := map[string]bool{
			"@yearly":   true,
			"@annually": true,
			"@monthly":  true,
			"@weekly":   true,
			"@daily":    true,
			"@midnight": true,
			"@hourly":   true,
		}
		if validSpecial[expr] {
			return nil
		}
		if strings.HasPrefix(expr, "@every ") {
			_, err := p.ParseInterval(strings.TrimPrefix(expr, "@every "))
			return err
		}
		return fmt.Errorf("unknown special expression: %s", expr)
	}

	// 标准5字段
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return fmt.Errorf("cron expression must have 5 fields, got %d", len(fields))
	}

	fieldConfigs := []struct {
		field string
		min   int
		max   int
		name  string
	}{
		{fields[0], 0, 59, "minute"},
		{fields[1], 0, 23, "hour"},
		{fields[2], 1, 31, "day"},
		{fields[3], 1, 12, "month"},
		{fields[4], 0, 7, "weekday"},
	}

	for _, cfg := range fieldConfigs {
		if _, err := p.parseField(cfg.field, cfg.min, cfg.max); err != nil {
			return fmt.Errorf("invalid %s field: %w", cfg.name, err)
		}
	}

	return nil
}
