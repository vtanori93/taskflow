# TaskFlow - Guía de Configuración

## Requisitos Previos
- **Go 1.21+** instalado ([Descargar Go](https://go.dev/dl/))
- **Docker Desktop** instalado y ejecutándose
- **Docker Compose**
- (Opcional) PostgreSQL Client (psql) o cualquier cliente SQL

## Pasos para Levantar la Base de Datos

### 1. Iniciar los Contenedores Docker
```bash
docker-compose up -d
```

Este comando levantará:
- PostgreSQL en el puerto `5432`
- pgAdmin en el puerto `5050`

### 2. Verificar que los Contenedores Estén Ejecutándose
```bash
docker-compose ps
```

Deberías ver:
- `taskflow_postgres` (status: Up)
- `taskflow_pgadmin` (status: Up)

### 3. Ejecutar el Schema SQL

#### Opción A: Usando Docker Exec
```bash
docker exec -i taskflow_postgres psql -U postgres -d taskflow < backend/database/schema.sql
```

**Windows PowerShell:**
```powershell
Get-Content backend/database/schema.sql | docker exec -i taskflow_postgres psql -U postgres -d taskflow
```

#### Opción B: Usando psql (si lo tienes instalado)
```bash
psql -h localhost -U postgres -d taskflow -f backend/database/schema.sql
```
Contraseña: `postgres`

#### Opción C: Usando pgAdmin (Interfaz Web)
1. Abre tu navegador en `http://localhost:5050`
2. Inicia sesión con:
   - Email: `admin@taskflow.local`
   - Contraseña: `admin`
3. Agrega un nuevo servidor:
   - Name: `TaskFlow`
   - Host: `postgres` (nombre del servicio en Docker)
   - Port: `5432`
   - Username: `postgres`
   - Password: `postgres`
   - Database: `taskflow`
4. Abre Query Tool y pega el contenido de `backend/database/schema.sql`
5. Ejecuta el script (F5 o botón Execute)

### 4. Verificar que las Tablas se Crearon
```bash
docker exec -it taskflow_postgres psql -U postgres -d taskflow -c "\dt"
```

Deberías ver las tablas:
- `users`
- `tasks`

## Información de Conexión

**PostgreSQL:**
- Host: `localhost`
- Port: `5432`
- Database: `taskflow`
- User: `postgres`
- Password: `postgres`

**pgAdmin:**
- URL: `http://localhost:5050`
- Email: `admin@taskflow.local`
- Password: `admin`

## Comandos Útiles

### Detener los Contenedores
```bash
docker-compose down
```

### Detener y Eliminar Volúmenes (CUIDADO: Elimina los datos)
```bash
docker-compose down -v
```

### Ver Logs de PostgreSQL
```bash
docker-compose logs -f postgres
```

### Reiniciar los Contenedores
```bash
docker-compose restart
```

### Conectarse Directamente a PostgreSQL
```bash
docker exec -it taskflow_postgres psql -U postgres -d taskflow
```

## Ejecutar la Aplicación Backend (Go)

### 1. Verificar Instalación de Go
```bash
go version
```
Deberías ver `go version go1.21` o superior.

### 2. Navegar al Directorio Backend
```bash
cd backend
```

### 3. Instalar Dependencias
```bash
go mod download
```

Esto instalará todas las dependencias necesarias:
- Gin (Framework web)
- PostgreSQL driver (lib/pq)
- JWT (golang-jwt/jwt)
- Swagger (documentación API)
- Bcrypt (encriptación de contraseñas)
- UUID (generación de IDs)
- Y más...

### 4. Configurar Variables de Entorno

Asegúrate de que el archivo `.env` en la raíz del proyecto exista con:
```env
# PostgreSQL Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=taskflow
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

# Application Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/taskflow?sslmode=disable

# JWT Configuration (opcional - valores por defecto)
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRATION=24h

# Server Configuration (opcional - valores por defecto)
SERVER_PORT=8080
SERVER_ENV=development
```

### 5. Generar Documentación Swagger (Opcional)
```bash
swag init
```

Si no tienes `swag` instalado:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 6. Ejecutar la Aplicación

#### Modo Desarrollo
```bash
go run main.go
```

#### Compilar y Ejecutar
```bash
go build -o taskflow.exe
./taskflow.exe
```

#### Ejecutar con Air (Hot Reload - Opcional)
Instalar Air:
```bash
go install github.com/cosmtrek/air@latest
```

Ejecutar:
```bash
air
```

### 7. Verificar que el API está Funcionando

El servidor debería iniciarse en `http://localhost:8080`

**Endpoints disponibles:**
- `GET /health` - Health check
- `GET /swagger/index.html` - Documentación Swagger
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/tasks` - Listar tareas (requiere auth)
- `POST /api/v1/tasks` - Crear tarea (requiere auth)
- Y más...

**Probar el health check:**
```bash
curl http://localhost:8080/health
```

### 8. Acceder a la Documentación Swagger
Abre tu navegador en: `http://localhost:8080/swagger/index.html`

## Ejecutar Tests

Ejecuta las pruebas unitarias desde la carpeta `backend`.

### Tests Unitarios
```bash
cd backend
go test ./... -v
```

### Tests con Coverage
```bash
go test ./... -cover
```

### Tests de un Paquete Específico
```bash
go test ./internal/service -v
```

## Estructura del Proyecto

```
backend/
├── main.go                 # Punto de entrada
├── go.mod                  # Dependencias
├── database/
│   └── schema.sql          # Schema de base de datos
├── docs/                   # Documentación Swagger generada
├── internal/
│   ├── config/             # Configuración
│   ├── domain/             # Interfaces de dominio
│   ├── handler/            # Handlers HTTP
│   ├── middleware/         # Middlewares (auth, CORS, etc)
│   ├── models/             # Modelos de datos
│   ├── repository/         # Capa de datos (PostgreSQL)
│   ├── service/            # Lógica de negocio
│   └── utils/              # Utilidades (JWT, password, validation)
```

## Solución de Problemas

### Error: "go: command not found"
Instala Go desde https://go.dev/dl/ y asegúrate de que esté en tu PATH.

### Error: "cannot find package"
```bash
cd backend
go mod tidy
go mod download
```

### Puerto 8080 ya está en uso
Cambia el puerto en el archivo `.env`:
```env
SERVER_PORT=8081
```

### Error de conexión a la base de datos
- Verifica que Docker esté ejecutándose: `docker ps`
- Verifica que PostgreSQL esté corriendo: `docker-compose ps`
- Verifica las credenciales en `.env`
- Verifica la conexión: 
   ```bash
   docker exec -it taskflow_postgres psql -U postgres -d taskflow -c "SELECT 1;"
   ```

### Puerto 5432 ya está en uso
Si ya tienes PostgreSQL corriendo localmente:
- Detén tu PostgreSQL local, o
- Cambia el puerto en `docker-compose.yml` (ej: `"5433:5432"`)
- Actualiza `POSTGRES_PORT` y `DATABASE_URL` en `.env`

### Contenedor no inicia
```bash
docker-compose logs postgres
```

### Resetear la Base de Datos
```bash
docker-compose down -v
docker-compose up -d
# Luego vuelve a ejecutar schema.sql
docker exec -i taskflow_postgres psql -U postgres -d taskflow < backend/database/schema.sql
```

### Error: "Swagger not found"
Genera la documentación Swagger:
```bash
cd backend
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### Hot Reload no funciona
Verifica que Air esté instalado correctamente:
```bash
which air  # Linux/Mac
where air  # Windows
```
