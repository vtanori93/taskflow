import { useState } from 'react';
import { AuthHandler } from '../handlers/AuthHandler';

export function useAuthViewModel() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [user, setUser] = useState<any>(null);
  const handler = new AuthHandler();

  const login = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await handler.login(email, password);
      setUser(res.user);
      setLoading(false);
      return res;
    } catch (e: any) {
      setError(e.message);
      setLoading(false);
    }
  };

  const register = async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await handler.register(email, password);
      setUser(res.user);
      setLoading(false);
      return res;
    } catch (e: any) {
      setError(e.message);
      setLoading(false);
    }
  };

  const logout = async () => {
    setLoading(true);
    setError(null);
    try {
      await handler.logout();
      setUser(null);
      setLoading(false);
    } catch (e: any) {
      setError(e.message);
      setLoading(false);
    }
  };

  return { loading, error, user, login, register, logout };
}
