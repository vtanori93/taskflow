import { useNavigation } from '@react-navigation/native';
import type { NativeStackNavigationProp } from '@react-navigation/native-stack';
import React, { useState } from 'react';
import { ActivityIndicator, Button, Text, TextInput, View } from 'react-native';
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
    <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', padding: 16 }}>
      <Text style={{ fontSize: 24, marginBottom: 16 }}>Login</Text>
      <TextInput
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        autoCapitalize="none"
        keyboardType="email-address"
        style={{ borderWidth: 1, borderColor: '#ccc', width: '100%', marginBottom: 12, padding: 8, borderRadius: 6 }}
      />
      <TextInput
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        style={{ borderWidth: 1, borderColor: '#ccc', width: '100%', marginBottom: 12, padding: 8, borderRadius: 6 }}
      />
      {error && <Text style={{ color: 'red', marginBottom: 8 }}>{error}</Text>}
      {loading ? <ActivityIndicator /> : <Button title="Iniciar sesión" onPress={handleLogin} />}
      {success && <Text style={{ color: 'green', marginTop: 12 }}>¡Login exitoso!</Text>}
    </View>
  );
}
