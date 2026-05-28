export function mockContent(overrides: Record<string, any> = {}) {
  return {
    id: 1,
    tenant_id: 1,
    user_id: 1,
    brand_id: 1,
    title: 'OpenGEO BrandOS 品牌治理指南',
    body: '# OpenGEO BrandOS\n\n这是一篇关于品牌治理的示例文章。\n\n## 什么是品牌治理？\n\n品牌治理是指...',
    summary: 'OpenGEO BrandOS 品牌治理指南的摘要',
    content_type: 'article',
    status: 1,
    visibility: 'public',
    schema_markup: '',
    ai_optimization_score: 85.5,
    word_count: 1500,
    reading_time: 5,
    tags: '["品牌治理", "GEO", "AI优化"]',
    published_at: '2026-05-28T10:00:00Z',
    created_at: '2026-05-28T00:00:00Z',
    updated_at: '2026-05-28T10:00:00Z',
    ...overrides,
  };
}

export function mockContentList(count: number) {
  const titles = [
    'OpenGEO BrandOS 品牌治理指南',
    'AI 时代的内容优化策略',
    'GEO 与 SEO 的区别',
    '如何构建品牌知识图谱',
    '多渠道发布最佳实践',
    '品牌可信度评分详解',
    'AI 引用归因分析',
    '合规校验指南',
    '术语表管理教程',
    '品牌元数据配置说明',
  ];
  return Array.from({ length: count }, (_, i) =>
    mockContent({
      id: i + 1,
      title: titles[i % titles.length],
      status: i % 3 === 0 ? 0 : 1,
    })
  );
}
