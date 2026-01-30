import { useNavigation } from '@react-navigation/native';
import { Box, Button, Progress, ScrollView, Spinner, Text, VStack } from 'native-base';
import React, { useEffect, useState } from 'react';
import { Alert } from 'react-native';
import { Routes } from '../navigation/routes';
import { Session } from '../services/session';
import { useAuthViewModel } from '../viewmodels/AuthViewModel';
import { useStatsViewModel } from '../viewmodels/StatsViewModel';



export default function ProfileScreen() {
    const [user, setUser] = useState<any>(null);
    const { logout } = useAuthViewModel();
    const navigation = useNavigation();
    const { stats, loading, error, fetchStats } = useStatsViewModel();

    useEffect(() => {
        Session.getSession().then((session) => setUser(session.user));
        fetchStats();
    }, [fetchStats]);

    if (!user) {
        return (
            <Box flex={1} justifyContent="center" alignItems="center">
                <Text>Cargando perfil...</Text>
            </Box>
        );
    }

    return (
        <ScrollView contentContainerStyle={{ flexGrow: 1, justifyContent: 'center', alignItems: 'center', padding: 24 }}>
            <VStack space={4} w="100%" maxW="400px">
                <Text fontSize="lg">Nombre: {user.name}</Text>
                <Text fontSize="lg" mb={4}>Correo: {user.email}</Text>
                <Button colorScheme="danger" onPress={() => {
                    Alert.alert(
                        'Cerrar sesión',
                        '¿Deseas cerrar sesión?',
                        [
                            { text: 'Cancelar', style: 'cancel' },
                            {
                                text: 'Sí',
                                style: 'destructive',
                                onPress: async () => {
                                    await logout();
                                    navigation.reset({ index: 0, routes: [{ name: Routes.Login }] });
                                },
                            },
                        ]
                    );
                }}>
                    Cerrar sesión
                </Button>
                {loading && <Spinner />}
                {error && <Text color="red.500">{error}</Text>}
                {stats && (
                    <VStack w="100%" space={4} mt={4}>
                        <Text fontWeight="bold">Estadísticas de tareas</Text>
                        <VStack space={2}>
                            <Text>Pendientes ({stats.pending_count})</Text>
                            <Progress value={stats.pending_count} max={stats.total_tasks || 1} />
                        </VStack>
                        <VStack space={2}>
                            <Text>En progreso ({stats.in_progress_count})</Text>
                            <Progress value={stats.in_progress_count} max={stats.total_tasks || 1} colorScheme="info" />
                        </VStack>
                        <VStack space={2}>
                            <Text>Completadas ({stats.completed_count})</Text>
                            <Progress value={stats.completed_count} max={stats.total_tasks || 1} colorScheme="success" />
                        </VStack>
                        <VStack space={2}>
                            <Text>Canceladas ({stats.cancelled_count})</Text>
                            <Progress value={stats.cancelled_count} max={stats.total_tasks || 1} colorScheme="danger" />
                        </VStack>
                        <VStack space={2}>
                            <Text>Alta prioridad ({stats.high_priority_count})</Text>
                            <Progress value={stats.high_priority_count} max={stats.total_tasks || 1} colorScheme="warning" />
                        </VStack>
                        <VStack space={2}>
                            <Text>Vencidas ({stats.overdue_count})</Text>
                            <Progress value={stats.overdue_count} max={stats.total_tasks || 1} colorScheme="rose" />
                        </VStack>
                        <Text mt={2}>Total tareas: {stats.total_tasks}</Text>
                    </VStack>
                )}
            </VStack>
        </ScrollView>
    );
}

