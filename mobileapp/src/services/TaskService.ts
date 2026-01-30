import api from './api';

export class TaskService {
  async getTasks(params?: { page?: number; page_size?: number }) {
    const response = await api.get('/tasks', { params });
    // El backend responde con { data: { tasks: [...] } }
    return response.data.data.tasks;
  }

  async getTaskById(id: string) {
    const response = await api.get(`/tasks/${id}`);
    return response.data.data;
  }

  async createTask(task: any) {
    const response = await api.post('/tasks', task);
    return response.data.data;
  }

  async updateTask(task: any) {
    const response = await api.put(`/tasks/${task.id}`, task);
    return response.data.data;
  }

  async deleteTask(id: string) {
    await api.delete(`/tasks/${id}`);
    return true;
  }
}
