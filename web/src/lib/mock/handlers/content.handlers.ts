import { http, HttpResponse } from 'msw';
import { mockContent, mockContentList } from '../fixtures/content';

export const contentHandlers = [
  http.get('/api/v1/contents', () => {
    return HttpResponse.json({
      contents: mockContentList(10),
      total: 10,
      page: 1,
      page_size: 20,
      total_pages: 1,
    });
  }),

  http.get('/api/v1/content/:contentId', ({ params }) => {
    return HttpResponse.json(mockContent({ id: Number(params.contentId) }));
  }),

  http.post('/api/v1/contents', async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json(mockContent({ ...body as any, id: Date.now() }), { status: 201 });
  }),

  http.put('/api/v1/content/:contentId', async ({ params, request }) => {
    const body = await request.json();
    return HttpResponse.json(mockContent({ ...body as any, id: Number(params.contentId) }));
  }),

  http.delete('/api/v1/content/:contentId', () => {
    return HttpResponse.json({ success: true });
  }),
];
