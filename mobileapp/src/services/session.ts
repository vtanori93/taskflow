import AsyncStorage from '@react-native-async-storage/async-storage';

export const Session = {
  async setSession({ access_token, refresh_token, user }: { access_token: string, refresh_token: string, user: any }) {
    await AsyncStorage.setItem('access_token', access_token);
    await AsyncStorage.setItem('refresh_token', refresh_token);
    await AsyncStorage.setItem('user', JSON.stringify(user));
  },
  async clearSession() {
    await AsyncStorage.removeItem('access_token');
    await AsyncStorage.removeItem('refresh_token');
    await AsyncStorage.removeItem('user');
  },
  async getSession() {
    const access_token = await AsyncStorage.getItem('access_token');
    const refresh_token = await AsyncStorage.getItem('refresh_token');
    const user = await AsyncStorage.getItem('user');
    return {
      access_token,
      refresh_token,
      user: user ? JSON.parse(user) : null,
    };
  },
};
