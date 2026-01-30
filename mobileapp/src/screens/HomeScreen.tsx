import { Box, FlatList, Spinner, Text, VStack } from 'native-base';
import React, { useEffect } from 'react';
import { useTaskViewModel } from '../viewmodels/TaskViewModel';

export default function HomeScreen() {
  const { tasks, loading, error, fetchTasks } = useTaskViewModel();

  useEffect(() => {
    // page_size grande para traer todas las tareas (ajusta según tu backend)
    fetchTasks({ page: 1, page_size: 100 });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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
      <Text fontSize="2xl" mb={4} textAlign="center">Tareas</Text>
      <FlatList
        data={tasks}
        keyExtractor={(item) => item.id?.toString()}
        renderItem={({ item }) => (
          <Box borderWidth={1} borderColor="#ccc" borderRadius={8} p={4} mb={3}>
            <VStack space={1}>
              <Text fontWeight="bold">{item.title}</Text>
              <Text>{item.description}</Text>
              <Text color="gray.500">Estado: {item.status}</Text>
              <Text color="gray.500">Prioridad: {item.priority}</Text>
              {item.due_date && <Text color="gray.400">Vence: {item.due_date}</Text>}
            </VStack>
          </Box>
        )}
        ListEmptyComponent={<Text textAlign="center">No hay tareas.</Text>}
      />
    </Box>
  );
}
