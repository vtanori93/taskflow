import { useState } from 'react';
import { TaskHandler } from '../handlers/TaskHandler';

export function useTaskViewModel() {
  const [tasks, setTasks] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const handler = new TaskHandler();

  const fetchTasks = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await handler.getTasks();
      setTasks(res);
      setLoading(false);
    } catch (e: any) {
      setError(e.message);
      setLoading(false);
    }
  };

  return { tasks, loading, error, fetchTasks };
}
