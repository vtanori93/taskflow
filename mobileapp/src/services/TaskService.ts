// Dummy TaskService
export class TaskService {
  async getTasks() {
    return [
      { id: 1, title: 'Tarea 1', status: 'pending', priority: 'high' },
      { id: 2, title: 'Tarea 2', status: 'completed', priority: 'low' },
    ];
  }

  async getTaskById(id: number) {
    return { id, title: `Tarea ${id}`, status: 'pending', priority: 'medium' };
  }

  async createTask(task: any) {
    return { ...task, id: Math.random() };
  }

  async updateTask(task: any) {
    return task;
  }

  async deleteTask(id: number) {
    return true;
  }
}
