import { useNavigation } from '@react-navigation/native';
import React, { useEffect } from 'react';
import { Image, View } from 'react-native';
import api from '../services/api';
import { Session } from '../services/session';

export default function SplashScreen() {
  const navigation = useNavigation();
  useEffect(() => {
    let isActive = true;
    const checkSession = async () => {
      await new Promise(res => setTimeout(res, 5000)); // Espera 5 segundos
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
    <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#fff' }}>
      <Image
        source={{ uri: 'https://media4.giphy.com/media/v1.Y2lkPTZjMDliOTUyZDBjamVjN3pnYmtwNjEzam1jeTRnOHVqejdobWIwMTR5Y2wyMnkyYyZlcD12MV9zdGlja2Vyc19zZWFyY2gmY3Q9cw/X0Vihujvo7FiEIXOsm/giphy.gif' }}
        style={{ width: 180, height: 180, resizeMode: 'contain' }}
      />
    </View>
  );
}
