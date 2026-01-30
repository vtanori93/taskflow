import { useNavigation } from '@react-navigation/native';
import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { Box, Button, Input, KeyboardAvoidingView, ScrollView, Spinner, Text, VStack } from 'native-base';
import React, { useState } from 'react';
import { Platform } from 'react-native';
import { useAuthViewModel } from '../viewmodels/AuthViewModel';

export default function LoginScreen() {
  const { loading, error, user, login } = useAuthViewModel();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [success, setSuccess] = useState(false);
  const navigation = useNavigation<NativeStackNavigationProp<any>>();

  const handleLogin = async () => {
    const res = await login(email, password);
    if (res && res.user) {
      setSuccess(true);
      navigation.reset({ index: 0, routes: [{ name: 'Main' }] });
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
            <Input
              placeholder="Email"
              value={email}
              onChangeText={setEmail}
              autoCapitalize="none"
              keyboardType="email-address"
              returnKeyType="next"
              fontSize={16}
            />
            <Input
              placeholder="Password"
              value={password}
              onChangeText={setPassword}
              type="password"
              fontSize={16}
              returnKeyType="done"
            />
            {error && <Text color="red.500" textAlign="center">{error}</Text>}
            {loading ? (
              <Spinner accessibilityLabel="Cargando" />
            ) : (
              <Button onPress={handleLogin} colorScheme="primary">Iniciar sesión</Button>
            )}
            {success && <Text color="green.500" textAlign="center">¡Login exitoso!</Text>}
            <Button variant="outline" mt={2} onPress={() => navigation.navigate('Register')}>Registrarse</Button>
          </VStack>
        </Box>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
