import { Routes } from './routes';

export type RootStackParamList = {
  [Routes.Splash]: undefined;
  [Routes.Login]: undefined;
  [Routes.Register]: undefined;
  [Routes.Main]: undefined;
  [Routes.Home]: undefined;
  [Routes.Profile]: undefined;
  [Routes.TaskDetail]: { taskId?: string } | undefined;
  [Routes.AddTask]: undefined;
};
