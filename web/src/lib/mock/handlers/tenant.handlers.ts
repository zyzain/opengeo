import { http, HttpResponse } from 'msw';
import { mockTenant, mockTenantList, mockTenantQuota } from '../fixtures/tenant';

export const tenantHandlers = [
  http.get('/api/v1/tenant/:tenantId', ({ params }) => {
    return HttpResponse.json(mockTenant({ id: Number(params.tenantId) }));
  }),

  http.get('/api/v1/tenants', () => {
    return HttpResponse.json({
      tenants: mockTenantList(5),
      total: 5,
      page: 1,
      page_size: 20,
      total_pages: 1,
    });
  }),

  http.post('/api/v1/tenants', async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json(mockTenant({ ...body as any, id: Date.now() }), { status: 201 });
  }),

  http.put('/api/v1/tenant/:tenantId', async ({ params, request }) => {
    const body = await request.json();
    return HttpResponse.json(mockTenant({ ...body as any, id: Number(params.tenantId) }));
  }),

  http.delete('/api/v1/tenant/:tenantId', () => {
    return HttpResponse.json({ success: true });
  }),

  http.get('/api/v1/tenant/:tenantId/quota', ({ params }) => {
    return HttpResponse.json(mockTenantQuota(Number(params.tenantId)));
  }),
];
