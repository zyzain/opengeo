#!/bin/bash

# OpenGEO BrandOS - MSW Mock 生成脚本
# 基于 Proto 文件自动生成 MSW Handler

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
PROTO_DIR="$ROOT_DIR/proto"
MOCK_DIR="$ROOT_DIR/web/src/lib/mock"

echo "=== OpenGEO Mock 生成 ==="
echo "Proto 目录: $PROTO_DIR"
echo "Mock 目录: $MOCK_DIR"

# 创建 Mock 目录结构
mkdir -p "$MOCK_DIR/handlers"
mkdir -p "$MOCK_DIR/fixtures"

# 生成 Tenant Mock
cat > "$MOCK_DIR/handlers/tenant.handlers.ts" << 'EOF'
import { http, HttpResponse } from 'msw';
import { mockTenant, mockTenantList } from '../fixtures/tenant';

export const tenantHandlers = [
  // 获取租户信息
  http.get('/api/v1/tenant/:tenantId', ({ params }) => {
    const tenant = mockTenant({ id: Number(params.tenantId) });
    return HttpResponse.json(tenant);
  }),

  // 列出租户
  http.get('/api/v1/tenants', () => {
    return HttpResponse.json({
      tenants: mockTenantList(5),
      total: 5,
      page: 1,
      page_size: 20,
    });
  }),

  // 创建租户
  http.post('/api/v1/tenants', async ({ request }) => {
    const body = await request.json();
    const tenant = mockTenant({ ...body as any, id: Date.now() });
    return HttpResponse.json(tenant, { status: 201 });
  }),

  // 更新租户
  http.put('/api/v1/tenant/:tenantId', async ({ params, request }) => {
    const body = await request.json();
    const tenant = mockTenant({ ...body as any, id: Number(params.tenantId) });
    return HttpResponse.json(tenant);
  }),

  // 获取租户配额
  http.get('/api/v1/tenant/:tenantId/quota', ({ params }) => {
    return HttpResponse.json({
      tenant_id: Number(params.tenantId),
      brand_limit: 10,
      brand_count: 3,
      user_limit: 20,
      user_count: 5,
      storage_limit: 1073741824,
      storage_used: 536870912,
      api_quota: 1000,
      api_used: 250,
      quota_reset_at: '2026-06-01T00:00:00Z',
    });
  }),
];
EOF

# 生成 Brand Mock
cat > "$MOCK_DIR/handlers/brand.handlers.ts" << 'EOF'
import { http, HttpResponse } from 'msw';
import { mockBrand, mockBrandList, mockBrandMetadata, mockGlossaryEntry } from '../fixtures/brand';

export const brandHandlers = [
  // 列出品牌
  http.get('/api/v1/brands', () => {
    return HttpResponse.json({
      brands: mockBrandList(5),
      total: 5,
      page: 1,
      page_size: 20,
    });
  }),

  // 获取品牌信息
  http.get('/api/v1/brand/:brandId', ({ params }) => {
    const brand = mockBrand({ id: Number(params.brandId) });
    return HttpResponse.json(brand);
  }),

  // 创建品牌
  http.post('/api/v1/brands', async ({ request }) => {
    const body = await request.json();
    const brand = mockBrand({ ...body as any, id: Date.now() });
    return HttpResponse.json(brand, { status: 201 });
  }),

  // 更新品牌
  http.put('/api/v1/brand/:brandId', async ({ params, request }) => {
    const body = await request.json();
    const brand = mockBrand({ ...body as any, id: Number(params.brandId) });
    return HttpResponse.json(brand);
  }),

  // 删除品牌
  http.delete('/api/v1/brand/:brandId', () => {
    return HttpResponse.json({ success: true });
  }),

  // 获取品牌元数据
  http.get('/api/v1/brand/:brandId/metadata', ({ params }) => {
    return HttpResponse.json(mockBrandMetadata(Number(params.brandId)));
  }),

  // 更新品牌元数据
  http.put('/api/v1/brand/:brandId/metadata', async ({ params, request }) => {
    const body = await request.json();
    return HttpResponse.json({ ...body, brand_id: Number(params.brandId) });
  }),

  // 列出术语表
  http.get('/api/v1/brand/:brandId/glossary', ({ params }) => {
    return HttpResponse.json({
      entries: Array.from({ length: 10 }, (_, i) => mockGlossaryEntry({ brand_id: Number(params.brandId) })),
      total: 10,
      page: 1,
      page_size: 20,
    });
  }),

  // 创建术语
  http.post('/api/v1/brand/:brandId/glossary', async ({ params, request }) => {
    const body = await request.json();
    const entry = mockGlossaryEntry({ ...body as any, brand_id: Number(params.brandId), id: Date.now() });
    return HttpResponse.json(entry, { status: 201 });
  }),
];
EOF

# 生成 Tenant Fixtures
cat > "$MOCK_DIR/fixtures/tenant.ts" << 'EOF'
export function mockTenant(overrides: Record<string, any> = {}) {
  return {
    id: 1,
    name: '默认租户',
    slug: 'default',
    domain: 'default.opengeo.com',
    logo_url: '',
    plan: 'free',
    status: 1,
    brand_limit: 10,
    user_limit: 20,
    storage_limit: 1073741824,
    api_quota: 1000,
    api_used: 250,
    quota_reset_at: '2026-06-01T00:00:00Z',
    settings: {},
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-05-28T00:00:00Z',
    ...overrides,
  };
}

export function mockTenantList(count: number) {
  return Array.from({ length: count }, (_, i) =>
    mockTenant({
      id: i + 1,
      name: `租户 ${i + 1}`,
      slug: `tenant-${i + 1}`,
    })
  );
}
EOF

# 生成 Brand Fixtures
cat > "$MOCK_DIR/fixtures/brand.ts" << 'EOF'
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
    mockBrand({
      id: i + 1,
      name: `品牌 ${i + 1}`,
      slug: `brand-${i + 1}`,
    })
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
      style_guide: '使用简洁专业的语言，避免过于技术化的表述',
    },
    audience_profiles: [
      {
        name: '企业用户',
        age_range: '25-45',
        interests: ['技术', '营销', 'AI'],
        pain_points: ['内容效率', '品牌一致性'],
        preferred_channels: ['微信', '知乎', '官网'],
        locations: ['北京', '上海', '深圳'],
        languages: ['zh-CN'],
      },
    ],
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
EOF

# 生成 Mock 入口文件
cat > "$MOCK_DIR/index.ts" << 'EOF'
import { setupWorker } from 'msw/browser';
import { tenantHandlers } from './handlers/tenant.handlers';
import { brandHandlers } from './handlers/brand.handlers';

export const worker = setupWorker(
  ...tenantHandlers,
  ...brandHandlers,
);
EOF

echo "=== Mock 生成完成 ==="
echo "生成文件："
echo "  - $MOCK_DIR/handlers/tenant.handlers.ts"
echo "  - $MOCK_DIR/handlers/brand.handlers.ts"
echo "  - $MOCK_DIR/fixtures/tenant.ts"
echo "  - $MOCK_DIR/fixtures/brand.ts"
echo "  - $MOCK_DIR/index.ts"
