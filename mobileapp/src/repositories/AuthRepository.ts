import { AuthService } from '../services/AuthService';

export class AuthRepository {
  private service = new AuthService();

  login(email: string, password: string) {
    return this.service.login(email, password);
  }

  register(email: string, password: string) {
    return this.service.register(email, password);
  }

  logout() {
    return this.service.logout();
  }
}
