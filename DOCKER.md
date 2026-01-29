# Configuración de Docker para TaskFlow

## Requisitos
- Docker
- Docker Compose

## Inicio rápido

### 1. Iniciar los servicios
```bash
docker-compose up -d
```

Esto levantará:
- **PostgreSQL** en puerto `5432`
- **pgAdmin** en puerto `5050` (admin@taskflow.local / admin)

### 3. Verificar estado
```bash
docker-compose ps
```

## Comandos útiles

### Ver logs
```bash
docker-compose logs postgres      # Logs de PostgreSQL
docker-compose logs               # Logs de todos los servicios
```

### Acceder a PostgreSQL desde CLI
```bash
docker-compose exec postgres psql -U postgres -d taskflow
```

### Detener servicios
```bash
docker-compose down
```

### Eliminar volúmenes (cuidado - elimina datos)
```bash
docker-compose down -v
```

## pgAdmin
- URL: http://localhost:5050
- Email: admin@taskflow.local
- Contraseña: admin

Para conectar a PostgreSQL en pgAdmin:
- Host: postgres
- Usuario: postgres
- Contraseña: postgres
- Base de datos: taskflow
