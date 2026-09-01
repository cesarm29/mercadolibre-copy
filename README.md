# Marketplace - Clone estilo MercadoLibre

Marketplace completo con carrito de compras estilo MercadoLibre.

## Stack

- **Backend**: Go (Gin, GORM, JWT, Swagger) — API REST para PostgreSQL
- **Frontend**: Astro — interfaz responsive estilo MercadoLibre
- **Base de datos**: PostgreSQL en Neon
- **Deploy**: Vercel

## URLs en producción

- **Frontend (web)**: https://mercadolibre-copy.vercel.app
- **Backend (API)**: https://mercadolibre-copy-backend.vercel.app
- **Swagger API**: https://mercadolibre-copy-backend.vercel.app/swagger/

## Estructura

```
backend/   → API REST en Go (JWT auth, roles, productos, carrito, ordenes)
frontend/  → Web en Astro (interfaz estilo MercadoLibre)
```

## Backend (Go)

Endpoints principales en `/api`:

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/auth/register` | Crear cuenta |
| POST | `/api/auth/login` | Iniciar sesión (JWT) |
| GET | `/api/auth/profile` | Perfil del usuario |
| GET | `/api/products` | Listar productos (filtros + paginación) |
| POST | `/api/products` | Crear producto (seller/admin) |
| GET | `/api/products/:id` | Detalle de producto |
| GET/POST | `/api/cart` | Ver / agregar al carrito |
| POST | `/api/orders` | Crear orden de compra |
| GET | `/api/orders` | Mis compras |
| GET | `/api/categories` | Categorías |
| GET | `/api/admin/users` | Usuarios (admin) |

Swagger disponible en `/swagger/`.

### Variables de entorno del backend

```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE, JWT_SECRET, SERVER_PORT
```

## Frontend (Astro)

Páginas: inicio, explorar productos, detalle de producto, carrito, checkout, mis compras, perfil, vender.

La constante `API` en cada página apunta a la URL del backend Go desplegado.

## Desarrollo local

```bash
# Backend
cd backend && go run main.go

# Frontend
cd frontend && npm install && npm run dev
```

## Producción

- Backend Go: desplegado como función serverless / contenedor en Vercel
- Frontend: build estático de Astro en Vercel
- PostgreSQL: Neon (gratuito)