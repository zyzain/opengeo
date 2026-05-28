export function mockTenant(overrides: Record<string, any> = {}) {
  return {
    id: 1,
    name: '默认租户',
    slug: 'default',
    domain: 'default.opengeo.com',
    logo_url: '',
    plan: 1,
    status: 1,
    brand_limit: 10,
    user_limit: 20,
    storage_limit: 1073741824,
    api_quota: 1000,
    api_used: 250,
    settings: {},
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-05-28T00:00:00Z',
    ...overrides,
  };
}

export function mockTenantList(count: number) {
  return Array.from({ length: count }, (_, i) =>
    mockTenant({ id: i + 1, name: `租户 ${i + 1}`, slug: `tenant-${i + 1}` })
  );
}

export function mockTenantQuota(tenantId: number) {
  return {
    tenant_id: tenantId,
    brand_limit: 10,
    brand_count: 3,
    user_limit: 20,
    user_count: 5,
    storage_limit: 1073741824,
    storage_used: 536870912,
    api_quota: 1000,
    api_used: 250,
    quota_reset_at: '2026-06-01T00:00:00Z',
  };
}
