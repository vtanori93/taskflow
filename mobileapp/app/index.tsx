import { Inter_400Regular, Inter_700Bold, useFonts } from '@expo-google-fonts/inter';
import { registerRootComponent } from 'expo';
import * as SplashScreen from 'expo-splash-screen';
import { NativeBaseProvider } from 'native-base';
import React from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import AppNavigator from '../src/navigation/AppNavigator';

SplashScreen.preventAutoHideAsync();

function Main() {
  const [fontsLoaded] = useFonts({
    Inter_400Regular,
    Inter_700Bold,
  });
  React.useEffect(() => {
    if (fontsLoaded) {
      SplashScreen.hideAsync();
    }
  }, [fontsLoaded]);
  if (!fontsLoaded) {
    return null;
  }
  return (
    <SafeAreaProvider>
      <NativeBaseProvider>
        <AppNavigator />
      </NativeBaseProvider>
    </SafeAreaProvider>
  );
}

registerRootComponent(Main);
