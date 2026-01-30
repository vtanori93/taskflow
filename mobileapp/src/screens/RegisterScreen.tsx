
import { useNavigation } from '@react-navigation/native';
import { Box, Button, Input, KeyboardAvoidingView, ScrollView, Text, VStack } from 'native-base';
import React, { useState } from 'react';
import { Alert, Platform } from 'react-native';
import { AuthService } from '../services/AuthService';

export default function RegisterScreen() {
  const navigation = useNavigation();
  const [form, setForm] = useState<{ email: string; name: string; password: string }>({
    email: '',
    name: '',
    password: '',
  });
  const [errors, setErrors] = useState<{ email?: string; name?: string; password?: string }>({});
  const [loading, setLoading] = useState(false);

  const validate = () => {
    const newErrors: { email?: string; name?: string; password?: string } = {};
    if (!form.email.trim()) newErrors.email = 'El email es obligatorio.';
    else if (!/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(form.email)) newErrors.email = 'Email inválido.';
    if (!form.name.trim()) newErrors.name = 'El nombre es obligatorio.';
    if (!form.password.trim()) newErrors.password = 'La contraseña es obligatoria.';
    else if (form.password.length < 6) newErrors.password = 'Mínimo 6 caracteres.';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleRegister = async () => {
    if (!validate()) return;
    setLoading(true);
    try {
      const authService = new AuthService();
      await authService.register(form.email, form.name, form.password);
      Alert.alert('Registro exitoso', '¡Usuario registrado!');
      navigation.goBack();
    } catch (e: any) {
      Alert.alert('Error', e.message || 'No se pudo registrar.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      keyboardVerticalOffset={60}
      style={{ flex: 1 }}
    >
      <ScrollView
        style={{ flex: 1 }}
        contentContainerStyle={{ flexGrow: 1 }}
        keyboardShouldPersistTaps="handled"
      >
        <Box flex={1} justifyContent="center" alignItems="center" p={4} bg="#fff">
          <VStack space={4} width="100%" maxW="100%">
            <Text fontSize={24} fontWeight="bold" textAlign="center">Registro</Text>
            <Input
              placeholder="Email"
              value={form.email}
              onChangeText={v => setForm(f => ({ ...f, email: v }))}
              autoCapitalize="none"
              keyboardType="email-address"
              fontSize={16}
            />
            {errors.email && <Text color="red.500">{errors.email}</Text>}
            <Input
              placeholder="Nombre"
              value={form.name}
              onChangeText={v => setForm(f => ({ ...f, name: v }))}
              fontSize={16}
            />
            {errors.name && <Text color="red.500">{errors.name}</Text>}
            <Input
              placeholder="Contraseña"
              value={form.password}
              onChangeText={v => setForm(f => ({ ...f, password: v }))}
              type="password"
              fontSize={16}
            />
            {errors.password && <Text color="red.500">{errors.password}</Text>}
            <Button mt={2} isLoading={loading} onPress={handleRegister} colorScheme="primary">Registrarse</Button>
          </VStack>
        </Box>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
