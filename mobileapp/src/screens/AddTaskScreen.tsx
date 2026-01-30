// import { CheckIcon, Select } from 'native-base';
import { Button, Input, ScrollView, Text } from 'native-base';
import React, { useState } from 'react';
import { Alert, KeyboardAvoidingView, Platform, StyleSheet } from 'react-native';
import { TaskService } from '../services/TaskService';

interface AddTaskScreenProps {
    navigation: any;
}

const AddTaskScreen: React.FC<AddTaskScreenProps> = ({ navigation }) => {
    const [form, setForm] = useState({
        title: '',
        description: '',
        dueDate: '',
        priority: '',
    });
    const [errors, setErrors] = useState<{ [key: string]: string }>({});


    const validate = () => {
        const newErrors: { [key: string]: string } = {};
        if (!form.title.trim()) newErrors.title = 'El título es obligatorio.';
        if (!form.description.trim()) newErrors.description = 'La descripción es obligatoria.';
        if (!form.dueDate.trim()) {
            newErrors.dueDate = 'La fecha de vencimiento es obligatoria.';
        } else if (!/^\d{4}-\d{2}-\d{2}$/.test(form.dueDate)) {
            newErrors.dueDate = 'Formato de fecha inválido (YYYY-MM-DD).';
        } else {
            // Validar que sea una fecha real
            const [yyyy, mm, dd] = form.dueDate.split('-').map(Number);
            const date = new Date(yyyy, mm - 1, dd);
            if (
                date.getFullYear() !== yyyy ||
                date.getMonth() + 1 !== mm ||
                date.getDate() !== dd
            ) {
                newErrors.dueDate = 'La fecha no es válida.';
            }
        }
        if (!form.priority.trim()) {
            newErrors.priority = 'La prioridad es obligatoria.';
        } else if (!['low', 'medium', 'high', 'urgent'].includes(form.priority.trim().toLowerCase())) {
            newErrors.priority = 'Prioridad inválida (low, medium, high, urgent).';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleAddTask = async () => {
        if (!validate()) return;
        try {
            const taskService = new TaskService();
            await taskService.createTask({
                title: form.title,
                description: form.description,
                due_date: form.dueDate,
                priority: form.priority,
            });
            Alert.alert('Tarea agregada', 'La tarea se ha agregado correctamente.');
            navigation.goBack();
        } catch (error: any) {
            Alert.alert('Error', error?.message || 'No se pudo agregar la tarea.');
        }
    };


    const styles = StyleSheet.create({
        error: {
            color: 'red',
            marginTop: 2,
            marginBottom: 4,
            fontSize: 13,
        },
        container: {
            flex: 1,
            padding: 20,
            backgroundColor: '#fff',
        },
        header: {
            flexDirection: 'row',
            alignItems: 'center',
            marginBottom: 24,
        },
        backButton: {
            paddingVertical: 8,
            paddingHorizontal: 8,
            marginRight: 8,
        },
        backButtonText: {
            color: '#007AFF',
            fontSize: 16,
        },
        headerTitle: {
            fontSize: 18,
            fontWeight: 'bold',
            flex: 1,
            textAlign: 'center',
            marginRight: 32, // para centrar el título visualmente
        },
        label: {
            fontWeight: 'bold',
            marginTop: 16,
        }
    });

    return (
        <KeyboardAvoidingView
            style={{ flex: 1 }}
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
            keyboardVerticalOffset={80}
        >
            <ScrollView flex={1} p={4} bg="#fff" keyboardShouldPersistTaps="handled">
                <Text style={styles.label}>Título</Text>
                <Input
                    value={form.title}
                    onChangeText={v => setForm(f => ({ ...f, title: v }))}
                    placeholder="Título de la tarea"
                    placeholderTextColor="#B0B0B0"
                    autoCapitalize="sentences"
                    mb={errors.title ? 0 : 4}
                    fontSize={18}
                />
                {errors.title ? <Text style={styles.error}>{errors.title}</Text> : null}

                <Text style={styles.label}>Descripción</Text>
                <Input
                    value={form.description}
                    onChangeText={v => setForm(f => ({ ...f, description: v }))}
                    placeholder="Descripción detallada"
                    placeholderTextColor="#B0B0B0"
                    multiline
                    mb={errors.description ? 0 : 4}
                    fontSize={18}
                    textAlignVertical="top"
                />
                {errors.description ? <Text style={styles.error}>{errors.description}</Text> : null}

                <Text style={styles.label}>Fecha de vencimiento</Text>
                <Input
                    value={form.dueDate}
                    onChangeText={v => setForm(f => ({ ...f, dueDate: v }))}
                    placeholder="YYYY-MM-DD"
                    placeholderTextColor="#B0B0B0"
                    keyboardType="numbers-and-punctuation"
                    maxLength={10}
                    mb={errors.dueDate ? 0 : 4}
                    fontSize={18}
                />
                {errors.dueDate ? <Text style={styles.error}>{errors.dueDate}</Text> : null}

                <Text style={styles.label}>Prioridad</Text>
                <Input
                    value={form.priority}
                    onChangeText={v => setForm(f => ({ ...f, priority: v }))}
                    placeholder="Prioridad (low, medium, high, urgent)"
                    placeholderTextColor="#B0B0B0"
                    autoCapitalize="none"
                    mb={errors.priority ? 0 : 4}
                    fontSize={18}
                />
                {errors.priority ? <Text style={styles.error}>{errors.priority}</Text> : null}

                <Button mt={4} onPress={handleAddTask}>Agregar Tarea</Button>
            </ScrollView>
        </KeyboardAvoidingView>
    );
};

export default AddTaskScreen;
