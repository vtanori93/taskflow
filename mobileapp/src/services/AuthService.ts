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

  async register(email: string, password: string) {
    // Implementar registro real si es necesario
    return { token: 'dummy-token', user: { id: 2, email } };
  }

  async logout() {
    await Session.clearSession();
    return true;
  }
}
