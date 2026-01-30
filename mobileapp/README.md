# TaskFlow Mobile App

TaskFlow es una aplicación móvil desarrollada con React Native y Expo para la gestión de tareas personales. Permite a los usuarios crear, visualizar, editar y eliminar tareas, así como gestionar su autenticación y perfil.

## Características principales
- Registro e inicio de sesión de usuarios
- Creación, edición y eliminación de tareas
- Visualización de detalles de tareas
- Pantalla de inicio y perfil de usuario
- Navegación entre pantallas
- Manejo de sesiones y autenticación

## Estructura del proyecto
```
app/                # Componentes principales y layout de la app
assets/             # Recursos estáticos (imágenes, fuentes)
src/
  handlers/         # Lógica de manejo de autenticación y tareas
  navigation/       # Navegación y rutas de la app
  repositories/     # Acceso a datos y lógica de repositorios
  screens/          # Pantallas principales de la app
  services/         # Servicios de API y lógica de negocio
  stores/           # (Opcional) Gestión de estado
  viewmodels/       # Lógica de presentación (MVVM)
```

## Instalación
1. Clona el repositorio:
   ```bash
   git clone <url-del-repositorio>
   ```
2. Instala las dependencias:
   ```bash
   npm install
   ```
3. Inicia el proyecto con Expo:
   ```bash
   npx expo start
   ```

## Scripts útiles
- `npm start` – Inicia el servidor de desarrollo de Expo
- `npm run android` – Ejecuta la app en un emulador/dispositivo Android
- `npm run ios` – Ejecuta la app en un emulador/dispositivo iOS (solo MacOS)

## Tecnologías utilizadas
- React Native
- Expo
- TypeScript
- Context API / MVVM

## Estructura de carpetas clave
- **app/**: Layout y punto de entrada de la app
- **src/screens/**: Pantallas principales (Login, Home, Perfil, etc.)
- **src/handlers/**: Lógica de manejo de autenticación y tareas
- **src/services/**: Servicios de API y lógica de negocio
- **src/navigation/**: Configuración de navegación y rutas

## Contribución
1. Haz un fork del repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-feature`)
3. Realiza tus cambios y haz commit (`git commit -am 'Agrega nueva feature'`)
4. Haz push a tu rama (`git push origin feature/nueva-feature`)
5. Abre un Pull Request

## Licencia
Este proyecto está bajo la licencia MIT.

- [Expo documentation](https://docs.expo.dev/): Learn fundamentals, or go into advanced topics with our [guides](https://docs.expo.dev/guides).
- [Learn Expo tutorial](https://docs.expo.dev/tutorial/introduction/): Follow a step-by-step tutorial where you'll create a project that runs on Android, iOS, and the web.

## Join the community

Join our community of developers creating universal apps.

- [Expo on GitHub](https://github.com/expo/expo): View our open source platform and contribute.
- [Discord community](https://chat.expo.dev): Chat with Expo users and ask questions.
