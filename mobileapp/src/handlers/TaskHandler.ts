import { TaskRepository } from '../repositories/TaskRepository';

export class TaskHandler {
  private repo = new TaskRepository();

  async getTasks() {
    return this.repo.getTasks();
  }

  async getTaskById(id: number) {
    return this.repo.getTaskById(id);
  }

  async createTask(task: any) {
    return this.repo.createTask(task);
  }

  async updateTask(task: any) {
    return this.repo.updateTask(task);
  }

  async deleteTask(id: number) {
    return this.repo.deleteTask(id);
  }
}
