import { useNavigation } from '@react-navigation/native';
import React, { useEffect } from 'react';
import { Text, View } from 'react-native';
import { Session } from '../services/session';
import api from '../services/api';

export default function SplashScreen() {
  const navigation = useNavigation();
  useEffect(() => {
    let isActive = true;
    const checkSession = async () => {
      const session = await Session.getSession();
      if (session && session.access_token) {
        // Verificar token llamando a /auth/profile
        try {
          const res = await api.get('/auth/profile', {
            headers: { Authorization: `Bearer ${session.access_token}` },
          });
          if (isActive && res.data && res.data.data) {
            navigation.reset({ index: 0, routes: [{ name: 'Main' }] });
            return;
          }
        } catch (e) {
          // Token inválido o expirado
        }
      }
      // Si no hay sesión o token inválido, limpiar y navegar a Login
      await Session.clearSession();
      if (isActive) navigation.reset({ index: 0, routes: [{ name: 'Login' }] });
    };
    checkSession();
    return () => { isActive = false; };
  }, [navigation]);
  return (
    <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
      <Text>Splash Screen</Text>
    </View>
  );
}
