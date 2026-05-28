package application

import (
	"context"
	"fmt"
	"strings"
)

// PromptInjector Prompt 注入框架
type PromptInjector struct {
	brandRepo    BrandRepository
	glossaryRepo GlossaryRepository
}

// BrandRepository 品牌仓储接口
type BrandRepository interface {
	FindByID(ctx context.Context, tenantID, id int64) (*Brand, error)
}

// GlossaryRepository 术语仓储接口
type GlossaryRepository interface {
	ListByBrandID(ctx context.Context, tenantID, brandID int64) ([]*GlossaryEntry, error)
}

// Brand 品牌信息
type Brand struct {
	ID          int64
	Name        string
	Description string
	Industry    string
}

// GlossaryEntry 术语条目
type GlossaryEntry struct {
	Term        string
	Definition  string
	Category    string
	IsForbidden bool
}

// NewPromptInjector 创建 Prompt 注入器
func NewPromptInjector(brandRepo BrandRepository, glossaryRepo GlossaryRepository) *PromptInjector {
	return &PromptInjector{
		brandRepo:    brandRepo,
		glossaryRepo: glossaryRepo,
	}
}

// BuildPrompt 构建注入品牌信息的 Prompt
func (p *PromptInjector) BuildPrompt(ctx context.Context, tenantID, brandID int64, basePrompt string) (string, error) {
	// 获取品牌信息
	brand, err := p.brandRepo.FindByID(ctx, tenantID, brandID)
	if err != nil {
		return "", fmt.Errorf("failed to get brand: %w", err)
	}

	// 获取术语表
	glossary, err := p.glossaryRepo.ListByBrandID(ctx, tenantID, brandID)
	if err != nil {
		return "", fmt.Errorf("failed to get glossary: %w", err)
	}

	// 构建品牌规范部分
	brandPrompt := p.buildBrandPrompt(brand)

	// 构建术语表部分
	glossaryPrompt := p.buildGlossaryPrompt(glossary)

	// 构建禁用词部分
	forbiddenPrompt := p.buildForbiddenPrompt(glossary)

	// 组合完整 Prompt
	fullPrompt := fmt.Sprintf(`## 品牌规范
%s

## 术语表
%s

## 禁用词
%s

## 用户请求
%s

## 输出要求
请根据上述品牌规范和术语表，生成符合品牌调性的内容。避免使用禁用词。`,
		brandPrompt, glossaryPrompt, forbiddenPrompt, basePrompt)

	return fullPrompt, nil
}

// buildBrandPrompt 构建品牌 Prompt
func (p *PromptInjector) buildBrandPrompt(brand *Brand) string {
	if brand == nil {
		return "无特定品牌规范"
	}

	return fmt.Sprintf(`品牌名称：%s
品牌描述：%s
所属行业：%s`, brand.Name, brand.Description, brand.Industry)
}

// buildGlossaryPrompt 构建术语表 Prompt
func (p *PromptInjector) buildGlossaryPrompt(glossary []*GlossaryEntry) string {
	if len(glossary) == 0 {
		return "无特定术语要求"
	}

	var lines []string
	for _, entry := range glossary {
		if entry.IsForbidden {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s：%s", entry.Term, entry.Definition))
	}

	return strings.Join(lines, "\n")
}

// buildForbiddenPrompt 构建禁用词 Prompt
func (p *PromptInjector) buildForbiddenPrompt(glossary []*GlossaryEntry) string {
	var forbidden []string
	for _, entry := range glossary {
		if entry.IsForbidden {
			forbidden = append(forbidden, entry.Term)
		}
	}

	if len(forbidden) == 0 {
		return "无禁用词"
	}

	return fmt.Sprintf("请避免使用以下词语：%s", strings.Join(forbidden, "、"))
}
