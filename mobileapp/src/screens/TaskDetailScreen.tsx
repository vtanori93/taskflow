import { useNavigation, useRoute } from '@react-navigation/native';
import { Box, Button, HStack, Input, Spinner, Text, VStack } from 'native-base';
import React, { useEffect, useState } from 'react';
import { TaskService } from '../services/TaskService';

export default function TaskDetailScreen() {
  const route = useRoute();
  const navigation = useNavigation();
  const { taskId } = route.params as { taskId: string };
  const [task, setTask] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [form, setForm] = useState<any>({});
  const [assignUser, setAssignUser] = useState('');
  const [assignLoading, setAssignLoading] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTask = async () => {
      setLoading(true);
      try {
        const service = new TaskService();
        const data = await service.getTaskById(taskId);
        setTask(data);
        setForm({ ...data });
      } catch (e: any) {
        setError(e.message || 'Error al cargar la tarea');
      } finally {
        setLoading(false);
      }
    };
    fetchTask();
  }, [taskId]);

  const handleUpdate = async () => {
    setLoading(true);
    try {
      const service = new TaskService();
      await service.updateTask({ ...form, id: taskId });
      setEditMode(false);
      setTask({ ...form });
    } catch (e: any) {
      setError(e.message || 'Error al actualizar');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    setDeleteLoading(true);
    try {
      const service = new TaskService();
      await service.deleteTask(taskId);
      navigation.goBack();
    } catch (e: any) {
      setError(e.message || 'Error al eliminar');
    } finally {
      setDeleteLoading(false);
    }
  };

  const handleAssign = async () => {
    setAssignLoading(true);
    try {
      await TaskService.prototype.assignTask(taskId, assignUser);
      setAssignUser('');
    } catch (e: any) {
      setError(e.message || 'Error al asignar');
    } finally {
      setAssignLoading(false);
    }
  };

  if (loading) {
    return (
      <Box flex={1} justifyContent="center" alignItems="center">
        <Spinner />
      </Box>
    );
  }

  if (error) {
    return (
      <Box flex={1} justifyContent="center" alignItems="center">
        <Text color="red.500">{error}</Text>
      </Box>
    );
  }

  return (
    <Box flex={1} p={4} bg="#fff">
      <VStack space={4}>
        <Text fontSize={22} fontWeight="bold">Detalle de Tarea</Text>
        {!editMode ? (
          <>
            <Text fontSize={18} fontWeight="bold">{task.title}</Text>
            <Text>{task.description}</Text>
            <Text color="gray.500">Prioridad: {task.priority}</Text>
            <Text color="gray.500">Estado: {task.status}</Text>
            <Text color="gray.400">Vence: {task.due_date}</Text>
            <HStack space={2} mt={4}>
              <Button onPress={() => setEditMode(true)} colorScheme="primary">Editar</Button>
              <Button onPress={handleDelete} colorScheme="danger" isLoading={deleteLoading}>Eliminar</Button>
            </HStack>
          </>
        ) : (
          <>
            <Text fontWeight="bold" mt={2}>Título</Text>
            <Input
              placeholder="Ej: Implementar login"
              value={form.title}
              onChangeText={v => setForm((f: any) => ({ ...f, title: v }))}
              fontSize={16}
            />
            <Text fontWeight="bold" mt={2}>Descripción</Text>
            <Input
              placeholder="Ej: Crear pantalla y lógica de login"
              value={form.description}
              onChangeText={v => setForm((f: any) => ({ ...f, description: v }))}
              fontSize={16}
            />
            <Text fontWeight="bold" mt={2}>Prioridad (low, medium, high, urgent)</Text>
            <Input
              placeholder="Ej: medium"
              value={form.priority}
              onChangeText={v => setForm((f: any) => ({ ...f, priority: v }))}
              fontSize={16}
            />
            <Text fontWeight="bold" mt={2}>Estado (pending, in_progress, completed, cancelled)</Text>
            <Input
              placeholder="Ej: pending"
              value={form.status}
              onChangeText={v => setForm((f: any) => ({ ...f, status: v }))}
              fontSize={16}
            />
            <Text fontWeight="bold" mt={2}>Fecha de vencimiento</Text>
            <Input
              placeholder="Ej: 2024-12-31"
              value={form.due_date}
              onChangeText={v => setForm((f: any) => ({ ...f, due_date: v }))}
              fontSize={16}
            />
            <HStack space={2} mt={4}>
              <Button onPress={handleUpdate} colorScheme="primary" isLoading={loading}>Guardar</Button>
              <Button onPress={() => setEditMode(false)} colorScheme="muted">Cancelar</Button>
            </HStack>
          </>
        )}
        <Box mt={6}>
          <Text fontWeight="bold" mb={2}>Asignar tarea a usuario</Text>
          <HStack space={2}>
            <Input
              placeholder="ID usuario"
              value={assignUser}
              onChangeText={setAssignUser}
              fontSize={16}
              flex={1}
            />
            <Button onPress={handleAssign} isLoading={assignLoading} colorScheme="primary">Asignar</Button>
          </HStack>
        </Box>
      </VStack>
    </Box>
  );
}
