import { useState, useEffect, useCallback } from 'react';

interface Snapshot {
  id: number;
  brand_id: number;
  version: string;
  snapshot_data: string;
  change_log: string;
  created_by: number;
  created_at: string;
}

interface SnapshotDiff {
  changes: Array<{
    field_path: string;
    old_value: string;
    new_value: string;
    change_type: 'added' | 'modified' | 'deleted';
  }>;
  additions: Array<{
    field_path: string;
    new_value: string;
  }>;
  deletions: Array<{
    field_path: string;
    old_value: string;
  }>;
}

export function useSnapshots(brandId: number) {
  const [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchSnapshots = useCallback(async () => {
    if (!brandId) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/snapshots`);
      if (!response.ok) throw new Error('Failed to fetch snapshots');
      const data = await response.json();
      setSnapshots(data.snapshots || []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, [brandId]);

  useEffect(() => {
    fetchSnapshots();
  }, [fetchSnapshots]);

  const createSnapshot = async (values: { version: string; change_log: string }) => {
    const response = await fetch(`/api/v1/brand/${brandId}/snapshots`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });
    if (!response.ok) throw new Error('Failed to create snapshot');
    await fetchSnapshots();
    return response.json();
  };

  const getSnapshot = async (snapshotId: number) => {
    const response = await fetch(`/api/v1/brand/${brandId}/snapshots/${snapshotId}`);
    if (!response.ok) throw new Error('Failed to get snapshot');
    return response.json();
  };

  const compareSnapshots = async (snapshotId1: number, snapshotId2: number): Promise<SnapshotDiff> => {
    const response = await fetch(`/api/v1/brand/${brandId}/snapshots/compare`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ snapshot_id_1: snapshotId1, snapshot_id_2: snapshotId2 }),
    });
    if (!response.ok) throw new Error('Failed to compare snapshots');
    return response.json();
  };

  return {
    snapshots,
    loading,
    error,
    refetch: fetchSnapshots,
    createSnapshot,
    getSnapshot,
    compareSnapshots,
  };
}
