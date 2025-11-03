# Go CRUD API

A simple REST API for product management built with Go, Gin, and PostgreSQL.

## Features

- Complete CRUD operations for products
- REST API endpoints
- PostgreSQL database integration
- Docker containerization
- Database migrations

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products` | Get all products |
| GET | `/products/:id` | Get product by ID |
| POST | `/products` | Create new product |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |

## Quick Start

1. Clone the repository
```bash
git clone <repository-url>
cd go_crud
```

2. Start with Docker Compose
```bash
docker-compose up -d
```

3. API will be available at `http://localhost:5001`

## Tech Stack

- **Go** 1.21+
- **Gin** - Web framework
- **PostgreSQL** - Database
- **pgx** - PostgreSQL driver
- **Docker** - Containerization

## Environment Variables

```env
DB_USER=usuario
DB_PASSWORD=senha
DB_NAME=crud_example
DB_HOST=localhost
DB_PORT=5433
API_PORT=5001
```

## License

MIT License
