import { TaskService } from '../services/TaskService';

export class TaskRepository {
  private service = new TaskService();

  getTasks() {
    return this.service.getTasks();
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
