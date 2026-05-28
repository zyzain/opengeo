import { setupWorker } from 'msw/browser';
import { tenantHandlers } from './handlers/tenant.handlers';
import { brandHandlers } from './handlers/brand.handlers';
import { contentHandlers } from './handlers/content.handlers';

export const worker = setupWorker(
  ...tenantHandlers,
  ...brandHandlers,
  ...contentHandlers,
);
