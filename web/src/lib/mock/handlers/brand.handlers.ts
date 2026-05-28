import { http, HttpResponse } from 'msw';
import { mockBrand, mockBrandList, mockBrandMetadata, mockGlossaryEntry, mockGlossaryList } from '../fixtures/brand';

export const brandHandlers = [
  http.get('/api/v1/brands', () => {
    return HttpResponse.json({
      brands: mockBrandList(5),
      total: 5,
      page: 1,
      page_size: 20,
      total_pages: 1,
    });
  }),

  http.get('/api/v1/brand/:brandId', ({ params }) => {
    return HttpResponse.json(mockBrand({ id: Number(params.brandId) }));
  }),

  http.post('/api/v1/brands', async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json(mockBrand({ ...body as any, id: Date.now() }), { status: 201 });
  }),

  http.put('/api/v1/brand/:brandId', async ({ params, request }) => {
    const body = await request.json();
    return HttpResponse.json(mockBrand({ ...body as any, id: Number(params.brandId) }));
  }),

  http.delete('/api/v1/brand/:brandId', () => {
    return HttpResponse.json({ success: true });
  }),

  http.get('/api/v1/brand/:brandId/metadata', ({ params }) => {
    return HttpResponse.json(mockBrandMetadata(Number(params.brandId)));
  }),

  http.put('/api/v1/brand/:brandId/metadata', async ({ params, request }) => {
    const body = await request.json();
    return HttpResponse.json({ ...body, brand_id: Number(params.brandId) });
  }),

  http.get('/api/v1/brand/:brandId/glossary', () => {
    return HttpResponse.json({
      entries: mockGlossaryList(10),
      total: 10,
      page: 1,
      page_size: 20,
      total_pages: 1,
    });
  }),

  http.post('/api/v1/brand/:brandId/glossary', async ({ params, request }) => {
    const body = await request.json();
    return HttpResponse.json(mockGlossaryEntry({ ...body as any, brand_id: Number(params.brandId), id: Date.now() }), { status: 201 });
  }),
];
