import { useCallback, useState } from 'react';
import { getTaskStats } from '../services/api';

export function useStatsViewModel() {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getTaskStats();
      setStats(data);
    } catch (e: any) {
      setError(e.message || 'Error al obtener estadísticas');
    } finally {
      setLoading(false);
    }
  }, []);

  return { stats, loading, error, fetchStats };
}