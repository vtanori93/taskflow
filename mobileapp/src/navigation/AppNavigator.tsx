import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import React from 'react';


import { Ionicons } from '@expo/vector-icons';
import HomeScreen from '../screens/HomeScreen';
import LoginScreen from '../screens/LoginScreen';
import ProfileScreen from '../screens/ProfileScreen';
import RegisterScreen from '../screens/RegisterScreen';
import SplashScreen from '../screens/SplashScreen';
import TaskDetailScreen from '../screens/TaskDetailScreen';
import { Routes } from './routes';
import { RootStackParamList } from './types';

const Stack = createNativeStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator();

function MainTabs() {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ color, size }) => {
          let iconName: keyof typeof Ionicons.glyphMap | undefined = undefined;
          if (route.name === Routes.Home) {
            iconName = 'list';
          } else if (route.name === Routes.Profile) {
            iconName = 'person';
          }
          return iconName ? <Ionicons name={iconName} size={size} color={color} /> : null;
        },
      })}
    >
      <Tab.Screen name={Routes.Home} component={HomeScreen} options={{ title: 'Tareas' }} />
      <Tab.Screen name={Routes.Profile} component={ProfileScreen} options={{ title: 'Perfil' }} />
    </Tab.Navigator>
  );
}

export default function AppNavigator() {
  return (
    <NavigationContainer>
      <Stack.Navigator initialRouteName={Routes.Splash} screenOptions={{ headerShown: false }}>
        <Stack.Screen name={Routes.Splash} component={SplashScreen} />
        <Stack.Screen name={Routes.Login} component={LoginScreen} />
        <Stack.Screen name={Routes.Register} component={RegisterScreen} />
        <Stack.Screen name={Routes.Main} component={MainTabs} />
        <Stack.Screen name={Routes.TaskDetail} component={TaskDetailScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}
