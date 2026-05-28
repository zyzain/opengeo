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

  const createBrand = async (values: Partial<Brand>) => {
    const response = await fetch('/api/v1/brands', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to create brand');
    const data = await response.json();
    await fetchBrands();
    return data;
  };

  const updateBrand = async (id: number, values: Partial<Brand>) => {
    const response = await fetch(`/api/v1/brand/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to update brand');
    const data = await response.json();
    await fetchBrands();
    return data;
  };

  const deleteBrand = async (id: number) => {
    const response = await fetch(`/api/v1/brand/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete brand');
    await fetchBrands();
  };

  return { 
    brands, 
    loading, 
    error, 
    refetch: fetchBrands,
    createBrand,
    updateBrand,
    deleteBrand
  };
}

export function useBrand(brandId: number) {
  const [brand, setBrand] = useState<Brand | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchBrand = useCallback(async () => {
    if (!brandId) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/brand/${brandId}`);
      if (!response.ok) throw new Error('Failed to fetch brand');
      const data = await response.json();
      setBrand(data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, [brandId]);

  useEffect(() => {
    fetchBrand();
  }, [fetchBrand]);

  const updateBrand = async (values: Partial<Brand>) => {
    const response = await fetch(`/api/v1/brand/${brandId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to update brand');
    const data = await response.json();
    await fetchBrand();
    return data;
  };

  return { brand, loading, error, refetch: fetchBrand, updateBrand };
}

export function useBrandMetadata(brandId: number) {
  const [metadata, setMetadata] = useState<BrandMetadata | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchMetadata = useCallback(async () => {
    if (!brandId) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/metadata`);
      if (!response.ok) throw new Error('Failed to fetch metadata');
      const data = await response.json();
      setMetadata(data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, [brandId]);

  useEffect(() => {
    fetchMetadata();
  }, [fetchMetadata]);

  const updateMetadata = async (values: Partial<BrandMetadata>) => {
    const response = await fetch(`/api/v1/brand/${brandId}/metadata`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to update metadata');
    const data = await response.json();
    await fetchMetadata();
    return data;
  };

  return { metadata, loading, error, refetch: fetchMetadata, updateMetadata };
}

export function useGlossary(brandId: number) {
  const [entries, setEntries] = useState<GlossaryEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchGlossary = useCallback(async () => {
    if (!brandId) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/glossary`);
      if (!response.ok) throw new Error('Failed to fetch glossary');
      const data = await response.json();
      setEntries(data.entries || []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, [brandId]);

  useEffect(() => {
    fetchGlossary();
  }, [fetchGlossary]);

  const createEntry = async (values: Partial<GlossaryEntry>) => {
    const response = await fetch(`/api/v1/brand/${brandId}/glossary`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to create glossary entry');
    const data = await response.json();
    await fetchGlossary();
    return data;
  };

  const updateEntry = async (entryId: number, values: Partial<GlossaryEntry>) => {
    const response = await fetch(`/api/v1/brand/${brandId}/glossary/${entryId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to update glossary entry');
    const data = await response.json();
    await fetchGlossary();
    return data;
  };

  const deleteEntry = async (entryId: number) => {
    const response = await fetch(`/api/v1/brand/${brandId}/glossary/${entryId}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete glossary entry');
    await fetchGlossary();
  };

  return { entries, loading, error, refetch: fetchGlossary, createEntry, updateEntry, deleteEntry };
}
