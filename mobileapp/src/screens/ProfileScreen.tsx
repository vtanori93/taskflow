import { useFocusEffect, useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import React, { useCallback, useEffect, useState } from 'react';
import { ActivityIndicator, Alert, Button, Dimensions, RefreshControl, ScrollView, Text, View } from 'react-native';
import { BarChart } from 'react-native-chart-kit';
import { Routes } from '../navigation/routes';
import { RootStackParamList } from '../navigation/types';
import { Session } from '../services/session';
import { useAuthViewModel } from '../viewmodels/AuthViewModel';
import { useStatsViewModel } from '../viewmodels/StatsViewModel';


export default function ProfileScreen() {
    const [user, setUser] = useState<any>(null);
    const [refreshing, setRefreshing] = useState(false);
    const { logout } = useAuthViewModel();
    const navigation = useNavigation<NativeStackNavigationProp<RootStackParamList>>();
    const { stats, loading, error, fetchStats } = useStatsViewModel();

    // Cargar usuario de sesión al enfocar pantalla
    useFocusEffect(
        useCallback(() => {
            let isActive = true;
            Session.getSession().then((session) => {
                if (isActive) setUser(session.user);
            });
            return () => { isActive = false; };
        }, [])
    );

    useEffect(() => {
        fetchStats();
    }, [fetchStats]);

    const onRefresh = async () => {
        setRefreshing(true);
        await Session.getSession().then((session) => setUser(session.user));
        await fetchStats();
        setRefreshing(false);
    };

    if (!user) {
        return (
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
                <Text>Cargando perfil...</Text>
            </View>
        );
    }

    return (
        <ScrollView
            contentContainerStyle={{ flexGrow: 1, justifyContent: 'center', alignItems: 'center', padding: 24 }}
            refreshControl={
                <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
            }
        >
            <Text style={{ fontSize: 18, marginBottom: 8 }}>Nombre: {user.name}</Text>
            <Text style={{ fontSize: 18, marginBottom: 24 }}>Correo: {user.email}</Text>
            <Button
                title="Cerrar sesión"
                color="#d32f2f"
                onPress={() => {
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
                }}
            />
            {loading && <ActivityIndicator />}
            {error && <Text style={{ color: 'red' }}>{error}</Text>}
            {stats && (
                <View style={{ width: '100%', marginTop: 16 }}>
                    <Text style={{ fontWeight: 'bold', marginBottom: 8 }}>Estadísticas de tareas</Text>
                    <BarChart
                        data={{
                            labels: [
                                'Pen', // Pendientes
                                'Prog', // En progreso
                                'Comp', // Completadas
                                'Canc', // Canceladas
                                'Pri', // Alta prioridad
                                'Venc', // Vencidas
                            ],
                            datasets: [
                                {
                                    data: [
                                        stats.pending_count,
                                        stats.in_progress_count,
                                        stats.completed_count,
                                        stats.cancelled_count,
                                        stats.high_priority_count,
                                        stats.overdue_count,
                                    ],
                                },
                            ],
                        }}
                        width={Dimensions.get('window').width - 48}
                        height={220}
                        yAxisLabel={''}
                        yAxisSuffix={''}
                        chartConfig={{
                            backgroundColor: '#fff',
                            backgroundGradientFrom: '#fff',
                            backgroundGradientTo: '#fff',
                            decimalPlaces: 0,
                            color: (opacity = 1) => `rgba(33, 150, 243, ${opacity})`,
                            labelColor: (opacity = 1) => `rgba(0, 0, 0, ${opacity})`,
                            style: { borderRadius: 16 },
                            propsForDots: { r: '6', strokeWidth: '2', stroke: '#2196f3' },
                        }}
                        style={{ marginVertical: 8, borderRadius: 16 }}
                    />
                    <View style={{ flexDirection: 'row', flexWrap: 'wrap', justifyContent: 'center', marginTop: 8 }}>
                        <Text style={{ marginHorizontal: 4 }}>Pen: Pendientes</Text>
                        <Text style={{ marginHorizontal: 4 }}>Prog: En progreso</Text>
                        <Text style={{ marginHorizontal: 4 }}>Comp: Completadas</Text>
                        <Text style={{ marginHorizontal: 4 }}>Canc: Canceladas</Text>
                        <Text style={{ marginHorizontal: 4 }}>Pri: Alta prioridad</Text>
                        <Text style={{ marginHorizontal: 4 }}>Venc: Vencidas</Text>
                    </View>
                    <Text style={{ marginTop: 8 }}>Total tareas: {stats.total_tasks}</Text>
                </View>
            )}
        </ScrollView>
    );
}

