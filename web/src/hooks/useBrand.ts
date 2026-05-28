import { useState, useEffect, useCallback } from 'react';

interface Brand {
  id: number;
  tenant_id: number;
  name: string;
  slug: string;
  description: string;
  logo_url: string;
  website: string;
  industry: string;
  status: number;
  created_at: string;
  updated_at: string;
}

interface BrandMetadata {
  brand_id: number;
  vi_profile: {
    primary_color: string;
    secondary_color: string;
    logo_url: string;
    font_family: string;
    brand_keywords: string[];
    slogan: string;
  };
  tone_profile: {
    formality: string;
    personality: string;
    avoid_words: string[];
    preferred_phrases: string[];
    style_guide: string;
  };
  audience_profiles: Array<{
    name: string;
    age_range: string;
    interests: string[];
    pain_points: string[];
  }>;
}

interface GlossaryEntry {
  id: number;
  brand_id: number;
  term: string;
  definition: string;
  category: string;
  aliases: string[];
  context: string;
  is_forbidden: boolean;
  is_preferred: boolean;
}

export function useBrands() {
  const [brands, setBrands] = useState<Brand[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchBrands = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch('/api/v1/brands');
      if (!response.ok) throw new Error('Failed to fetch brands');
      const data = await response.json();
      setBrands(data.brands || []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchBrands();
  }, [fetchBrands]);

  return { brands, loading, error, refetch: fetchBrands };
}

export function useBrand(brandId: number) {
  const [brand, setBrand] = useState<Brand | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!brandId) return;
    setLoading(true);
    fetch(`/api/v1/brand/${brandId}`)
      .then(res => res.json())
      .then(data => setBrand(data))
      .catch(err => setError(err))
      .finally(() => setLoading(false));
  }, [brandId]);

  return { brand, loading, error };
}

export function useBrandMetadata(brandId: number) {
  const [metadata, setMetadata] = useState<BrandMetadata | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!brandId) return;
    setLoading(true);
    fetch(`/api/v1/brand/${brandId}/metadata`)
      .then(res => res.json())
      .then(data => setMetadata(data))
      .finally(() => setLoading(false));
  }, [brandId]);

  return { metadata, loading };
}

export function useGlossary(brandId: number) {
  const [entries, setEntries] = useState<GlossaryEntry[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!brandId) return;
    setLoading(true);
    fetch(`/api/v1/brand/${brandId}/glossary`)
      .then(res => res.json())
      .then(data => setEntries(data.entries || []))
      .finally(() => setLoading(false));
  }, [brandId]);

  return { entries, loading };
}
