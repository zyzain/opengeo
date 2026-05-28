import {
  Brand,
  BrandMetadata,
  GlossaryEntry,
  BrandSnapshot,
  KnowledgeEntity,
  KnowledgeRelation,
  CreateBrandRequest,
  UpdateBrandRequest,
  CreateGlossaryEntryRequest,
  BulkImportGlossaryRequest,
  CreateKnowledgeEntityRequest,
  CreateKnowledgeRelationRequest,
} from '../../types/brand';

const BASE_URL = '/api/v1';

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }));
    throw new Error(error.message || `HTTP error! status: ${response.status}`);
  }

  return response.json();
}

// ========== 品牌管理 ==========

export const brandApi = {
  // 列出品牌
  list: (params?: { page?: number; page_size?: number; status?: number; industry?: string; keyword?: string }) => {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    const query = searchParams.toString();
    return request<{ brands: Brand[]; total: number; page: number; page_size: number }>(
      `${BASE_URL}/brands${query ? `?${query}` : ''}`
    );
  },

  // 获取品牌
  get: (id: number) => request<Brand>(`${BASE_URL}/brand/${id}`),

  // 创建品牌
  create: (data: CreateBrandRequest) =>
    request<Brand>(`${BASE_URL}/brands`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新品牌
  update: (id: number, data: UpdateBrandRequest) =>
    request<Brand>(`${BASE_URL}/brand/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除品牌
  delete: (id: number) =>
    request<{ success: boolean }>(`${BASE_URL}/brand/${id}`, {
      method: 'DELETE',
    }),
};

// ========== 品牌元数据 ==========

export const metadataApi = {
  // 获取元数据
  get: (brandId: number) =>
    request<BrandMetadata>(`${BASE_URL}/brand/${brandId}/metadata`),

  // 更新元数据
  update: (brandId: number, data: Partial<BrandMetadata>) =>
    request<BrandMetadata>(`${BASE_URL}/brand/${brandId}/metadata`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
};

// ========== 术语表 ==========

export const glossaryApi = {
  // 列出术语
  list: (brandId: number, params?: { page?: number; page_size?: number; category?: string; keyword?: string }) => {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    const query = searchParams.toString();
    return request<{ entries: GlossaryEntry[]; total: number }>(
      `${BASE_URL}/brand/${brandId}/glossary${query ? `?${query}` : ''}`
    );
  },

  // 创建术语
  create: (brandId: number, data: CreateGlossaryEntryRequest) =>
    request<GlossaryEntry>(`${BASE_URL}/brand/${brandId}/glossary`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新术语
  update: (brandId: number, entryId: number, data: Partial<GlossaryEntry>) =>
    request<GlossaryEntry>(`${BASE_URL}/brand/${brandId}/glossary/${entryId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除术语
  delete: (brandId: number, entryId: number) =>
    request<{ success: boolean }>(`${BASE_URL}/brand/${brandId}/glossary/${entryId}`, {
      method: 'DELETE',
    }),

  // 批量导入
  bulkImport: (brandId: number, data: BulkImportGlossaryRequest) =>
    request<{ imported_count: number; skipped_count: number; error_count: number }>(
      `${BASE_URL}/brand/${brandId}/glossary/bulk-import`,
      {
        method: 'POST',
        body: JSON.stringify(data),
      }
    ),

  // 批量导出
  bulkExport: (brandId: number, params?: { category?: string; format?: string }) => {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    const query = searchParams.toString();
    return request<{ entries: GlossaryEntry[]; download_url?: string }>(
      `${BASE_URL}/brand/${brandId}/glossary/export${query ? `?${query}` : ''}`
    );
  },
};

// ========== 品牌快照 ==========

export const snapshotApi = {
  // 列出快照
  list: (brandId: number, params?: { page?: number; page_size?: number }) => {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    const query = searchParams.toString();
    return request<{ snapshots: BrandSnapshot[]; total: number }>(
      `${BASE_URL}/brand/${brandId}/snapshots${query ? `?${query}` : ''}`
    );
  },

  // 创建快照
  create: (brandId: number, data: { version: string; change_log: string }) =>
    request<BrandSnapshot>(`${BASE_URL}/brand/${brandId}/snapshots`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 对比快照
  compare: (brandId: number, snapshotId1: number, snapshotId2: number) =>
    request<{ snapshot_1: BrandSnapshot; snapshot_2: BrandSnapshot; diff: any }>(
      `${BASE_URL}/brand/${brandId}/snapshots/compare`,
      {
        method: 'POST',
        body: JSON.stringify({ snapshot_id_1: snapshotId1, snapshot_id_2: snapshotId2 }),
      }
    ),
};

// ========== 知识图谱 ==========

export const knowledgeApi = {
  // 列出实体
  listEntities: (brandId: number, params?: { page?: number; page_size?: number; type?: string; keyword?: string }) => {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    const query = searchParams.toString();
    return request<{ entities: KnowledgeEntity[]; total: number }>(
      `${BASE_URL}/brand/${brandId}/knowledge/entities${query ? `?${query}` : ''}`
    );
  },

  // 创建实体
  createEntity: (brandId: number, data: CreateKnowledgeEntityRequest) =>
    request<KnowledgeEntity>(`${BASE_URL}/brand/${brandId}/knowledge/entities`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新实体
  updateEntity: (brandId: number, entityId: number, data: Partial<KnowledgeEntity>) =>
    request<KnowledgeEntity>(`${BASE_URL}/brand/${brandId}/knowledge/entities/${entityId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除实体
  deleteEntity: (brandId: number, entityId: number) =>
    request<{ success: boolean }>(`${BASE_URL}/brand/${brandId}/knowledge/entities/${entityId}`, {
      method: 'DELETE',
    }),

  // 搜索实体
  searchEntities: (brandId: number, query: string, params?: { type?: string; limit?: number }) => {
    const searchParams = new URLSearchParams({ query });
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    return request<{ results: Array<{ entity: KnowledgeEntity; similarity: number }> }>(
      `${BASE_URL}/brand/${brandId}/knowledge/entities/search?${searchParams.toString()}`
    );
  },

  // 列出关系
  listRelations: (brandId: number, params?: { entity_id?: number; type?: string }) => {
    const searchParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, String(value));
        }
      });
    }
    const query = searchParams.toString();
    return request<{ relations: KnowledgeRelation[]; total: number }>(
      `${BASE_URL}/brand/${brandId}/knowledge/relations${query ? `?${query}` : ''}`
    );
  },

  // 创建关系
  createRelation: (brandId: number, data: CreateKnowledgeRelationRequest) =>
    request<KnowledgeRelation>(`${BASE_URL}/brand/${brandId}/knowledge/relations`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 删除关系
  deleteRelation: (brandId: number, relationId: number) =>
    request<{ success: boolean }>(`${BASE_URL}/brand/${brandId}/knowledge/relations/${relationId}`, {
      method: 'DELETE',
    }),

  // 查询图谱
  queryGraph: (brandId: number, params: { center_entity_id: number; max_depth?: number; relation_types?: string[]; entity_types?: string[] }) =>
    request<{ graph: { nodes: any[]; edges: any[]; entity_details: Record<number, KnowledgeEntity> } }>(
      `${BASE_URL}/brand/${brandId}/knowledge/graph`,
      {
        method: 'POST',
        body: JSON.stringify(params),
      }
    ),
};

export default {
  brand: brandApi,
  metadata: metadataApi,
  glossary: glossaryApi,
  snapshot: snapshotApi,
  knowledge: knowledgeApi,
};
