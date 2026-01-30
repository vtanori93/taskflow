import { TaskService } from '../services/TaskService';

export class TaskRepository {
  private service = new TaskService();

  getTasks(params?: { page?: number; page_size?: number }) {
    return this.service.getTasks(params);
  }

  getTaskById(id: number) {
    return this.service.getTaskById(id);
  }

  createTask(task: any) {
    return this.service.createTask(task);
  }

  updateTask(task: any) {
    return this.service.updateTask(task);
  }

  deleteTask(id: number) {
    return this.service.deleteTask(id);
  }
}
