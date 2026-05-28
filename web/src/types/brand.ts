// 品牌相关类型定义

export interface Brand {
  id: number;
  tenant_id: number;
  name: string;
  slug: string;
  description: string;
  logo_url: string;
  website: string;
  industry: string;
  founded_year: number;
  headquarters: string;
  status: BrandStatus;
  settings: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export enum BrandStatus {
  Unspecified = 0,
  Active = 1,
  Archived = 2,
  Disabled = 3,
}

export interface BrandMetadata {
  brand_id: number;
  vi_profile: VIProfile;
  tone_profile: ToneProfile;
  audience_profiles: AudienceProfile[];
  competitor_list: CompetitorInfo[];
  brand_values: string[];
  unique_selling_points: string[];
  schema_version: string;
  created_at: string;
  updated_at: string;
}

export interface VIProfile {
  primary_color: string;
  secondary_color: string;
  logo_url: string;
  font_family: string;
  brand_keywords: string[];
  slogan: string;
}

export interface ToneProfile {
  formality: 'formal' | 'casual' | 'technical';
  personality: 'friendly' | 'professional' | 'playful' | 'authoritative';
  avoid_words: string[];
  preferred_phrases: string[];
  style_guide: string;
}

export interface AudienceProfile {
  name: string;
  age_range: string;
  interests: string[];
  pain_points: string[];
  preferred_channels: string[];
  locations: string[];
  languages: string[];
}

export interface CompetitorInfo {
  name: string;
  domain: string;
  description: string;
  strengths: string[];
  weaknesses: string[];
}

export interface GlossaryEntry {
  id: number;
  brand_id: number;
  term: string;
  definition: string;
  category: string;
  aliases: string[];
  context: string;
  is_forbidden: boolean;
  is_preferred: boolean;
  created_at: string;
  updated_at: string;
}

export interface BrandSnapshot {
  id: number;
  brand_id: number;
  version: string;
  snapshot_data: string;
  change_log: string;
  created_by: number;
  created_at: string;
}

export interface KnowledgeEntity {
  id: number;
  brand_id: number;
  name: string;
  type: EntityType;
  description: string;
  properties: Record<string, string>;
  authority_links: string[];
  tags: string[];
  embedding_id: string;
  created_at: string;
  updated_at: string;
}

export enum EntityType {
  Unspecified = 'unspecified',
  Brand = 'brand',
  Product = 'product',
  Person = 'person',
  Org = 'org',
  Event = 'event',
  Concept = 'concept',
  Location = 'location',
  Technology = 'technology',
}

export interface KnowledgeRelation {
  id: number;
  from_entity_id: number;
  to_entity_id: number;
  type: RelationType;
  weight: number;
  description: string;
  properties: Record<string, string>;
  created_at: string;
}

export enum RelationType {
  Unspecified = 'unspecified',
  IsA = 'is_a',
  PartOf = 'part_of',
  RelatedTo = 'related_to',
  CompetesWith = 'competes_with',
  Mentions = 'mentions',
  DependsOn = 'depends_on',
  Owns = 'owns',
  CreatedBy = 'created_by',
}

export interface BrandTrustScore {
  id: number;
  brand_id: number;
  score: number;
  dimensions: {
    search: number;
    social: number;
    compliance: number;
    citation: number;
  };
  factors: string;
  api_request_id: string;
  expires_at: string;
  created_at: string;
}

// 创建品牌请求
export interface CreateBrandRequest {
  name: string;
  slug: string;
  description?: string;
  logo_url?: string;
  website?: string;
  industry?: string;
  founded_year?: number;
  headquarters?: string;
}

// 更新品牌请求
export interface UpdateBrandRequest {
  name?: string;
  description?: string;
  logo_url?: string;
  website?: string;
  industry?: string;
  founded_year?: number;
  headquarters?: string;
  status?: BrandStatus;
  settings?: Record<string, string>;
}

// 创建术语请求
export interface CreateGlossaryEntryRequest {
  term: string;
  definition: string;
  category?: string;
  aliases?: string[];
  context?: string;
  is_forbidden?: boolean;
}

// 批量导入术语请求
export interface BulkImportGlossaryRequest {
  entries: CreateGlossaryEntryRequest[];
  overwrite_existing?: boolean;
}

// 创建知识实体请求
export interface CreateKnowledgeEntityRequest {
  name: string;
  type: EntityType;
  description?: string;
  properties?: Record<string, string>;
  authority_links?: string[];
  tags?: string[];
}

// 创建知识关系请求
export interface CreateKnowledgeRelationRequest {
  from_entity_id: number;
  to_entity_id: number;
  type: RelationType;
  weight?: number;
  description?: string;
  properties?: Record<string, string>;
}
