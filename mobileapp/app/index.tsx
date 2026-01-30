import { Inter_400Regular, Inter_700Bold, useFonts } from '@expo-google-fonts/inter';
import { registerRootComponent } from 'expo';
import AppLoading from 'expo-app-loading';
import React from 'react';
import AppNavigator from '../src/navigation/AppNavigator';

function Main() {
  const [fontsLoaded] = useFonts({
    Inter_400Regular,
    Inter_700Bold,
  });
  if (!fontsLoaded) {
    return <AppLoading />;
  }
  return <AppNavigator />;
}

registerRootComponent(Main);
