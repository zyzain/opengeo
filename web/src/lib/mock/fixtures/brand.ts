export function mockBrand(overrides: Record<string, any> = {}) {
  return {
    id: 1,
    tenant_id: 1,
    name: 'OpenGEO默认品牌',
    slug: 'opengeo-default',
    description: '系统自动创建的默认品牌',
    logo_url: '',
    website: 'https://opengeo.com',
    industry: '科技',
    founded_year: 2024,
    headquarters: '北京',
    status: 1,
    settings: {},
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-05-28T00:00:00Z',
    ...overrides,
  };
}

export function mockBrandList(count: number) {
  return Array.from({ length: count }, (_, i) =>
    mockBrand({ id: i + 1, name: `品牌 ${i + 1}`, slug: `brand-${i + 1}` })
  );
}

export function mockBrandMetadata(brandId: number) {
  return {
    brand_id: brandId,
    vi_profile: {
      primary_color: '#1890ff',
      secondary_color: '#52c41a',
      logo_url: '',
      font_family: 'PingFang SC',
      brand_keywords: ['GEO', 'AI优化', '品牌治理'],
      slogan: 'AI驱动的品牌治理平台',
    },
    tone_profile: {
      formality: 'professional',
      personality: 'friendly',
      avoid_words: [],
      preferred_phrases: ['AI驱动', '智能化', '高效'],
      style_guide: '使用简洁专业的语言',
    },
    audience_profiles: [{
      name: '企业用户',
      age_range: '25-45',
      interests: ['技术', '营销', 'AI'],
      pain_points: ['内容效率', '品牌一致性'],
      preferred_channels: ['微信', '知乎', '官网'],
      locations: ['北京', '上海', '深圳'],
      languages: ['zh-CN'],
    }],
    competitor_list: [],
    brand_values: ['创新', '专业', '可靠'],
    unique_selling_points: ['AI原生', '开源开放', '品牌治理'],
    schema_version: '1.0',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-05-28T00:00:00Z',
  };
}

export function mockGlossaryEntry(overrides: Record<string, any> = {}) {
  return {
    id: 1,
    brand_id: 1,
    term: 'GEO',
    definition: 'Generative Engine Optimization，生成式引擎优化',
    category: 'technology',
    aliases: ['生成式引擎优化'],
    context: 'GEO是针对AI搜索引擎的内容优化策略',
    is_forbidden: false,
    is_preferred: true,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-05-28T00:00:00Z',
    ...overrides,
  };
}

export function mockGlossaryList(count: number) {
  const terms = ['GEO', 'SEO', 'AIGC', 'LLM', '品牌治理', '内容优化', '知识图谱', '术语表', '元数据', '可信度'];
  return Array.from({ length: count }, (_, i) =>
    mockGlossaryEntry({ id: i + 1, term: terms[i % terms.length], definition: `${terms[i % terms.length]}的定义说明` })
  );
}
