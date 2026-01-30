import { AuthRepository } from '../repositories/AuthRepository';

export class AuthHandler {
  private repo = new AuthRepository();

  async login(email: string, password: string) {
    return this.repo.login(email, password);
  }

  async register(email: string, password: string) {
    return this.repo.register(email, password);
  }

  async logout() {
    return this.repo.logout();
  }
}
