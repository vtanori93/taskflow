import api from './api';
import { Session } from './session';

export class AuthService {
  async login(email: string, password: string) {
    try {
      const response = await api.post('/auth/login', { email, password });
      const { access_token, refresh_token, user } = response.data.data;
      await Session.setSession({ access_token, refresh_token, user });
      return { access_token, refresh_token, user };
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Error en login');
    }
  }

  async register(email: string, name: string, password: string) {
    try {
      const response = await api.post('/auth/register', { email, name, password });
      return response.data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Error en registro');
    }
  }

  async logout() {
    await Session.clearSession();
    return true;
  }
}
