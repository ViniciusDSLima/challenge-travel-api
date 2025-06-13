# API de Viagens

API RESTful para gerenciamento de viagens, desenvolvida em Go utilizando Clean Architecture.

## 🚀 Tecnologias

- Go 1.24+
- PostgreSQL
- Docker
- Docker Compose
- Make

## 📋 Pré-requisitos

- Go 1.24 ou superior
- Docker e Docker Compose
- Make

## 🔧 Instalação

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/challenge-travel-api.git
cd challenge-travel-api
```

3. Execute o projeto usando Docker:

Para iniciar a aplicação, execute o seguinte comando:

```bash
docker compose up -d
```

## 🏗️ Estrutura do Projeto

```
.
├── cmd/                    # Ponto de entrada da aplicação
├── config/                 # Configurações da aplicação
├── docs/                   # Documentação adicional
├── internal/              # Código fonte da aplicação
│   ├── domain/           # Entidades e regras de negócio
│   ├── infrastructure/   # Implementações concretas (repositórios, etc)
│   ├── interface/        # Handlers HTTP, gRPC, etc
│   ├── usecase/          # Casos de uso da aplicação
│   └── utils/            # Utilitários
├── migrations/           # Migrações do banco de dados
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

## 📚 Documentação da API

A documentação completa da API está disponível através do Swagger UI. Para acessá-la:

1. Acesse a URL: `http://localhost:8080/swagger/index.html`

## 👥 Usuários de Teste

A aplicação vem com alguns usuários pré-configurados para testes:

### Administrador
- Email: admin@example.com
- Senha: admin123
- Role: ADMIN

### Usuário Comum
- Email: user@example.com
- Senha: user123
- Role: USER

Para fazer login, use o endpoint `/api/v1/auth/login` com as credenciais acima.

### Endpoints

#### Viagens

##### Listar Viagens
```http
GET /api/v1/travels
```

**Resposta**
```json
{
    "data": [
        {
            "id": "uuid",
            "origin": "São Paulo",
            "destination": "Rio de Janeiro",
            "departure_time": "2024-03-20T10:00:00Z",
            "arrival_time": "2024-03-20T12:00:00Z",
            "price": 150.00,
            "available_seats": 45,
            "status": "SCHEDULED"
        }
    ]
}
```

##### Criar Viagem
```http
POST /api/v1/travels
```

**Request Body**
```json
{
    "origin": "São Paulo",
    "destination": "Rio de Janeiro",
    "departure_time": "2024-03-20T10:00:00Z",
    "arrival_time": "2024-03-20T12:00:00Z",
    "price": 150.00,
    "available_seats": 45
}
```

##### Buscar Viagem por ID
```http
GET /api/v1/travels/{id}
```

##### Atualizar Viagem
```http
PUT /api/v1/travels/{id}
```

##### Cancelar Viagem
```http
DELETE /api/v1/travels/{id}
```

## 🛠️ Comandos Make

### Migrações do Banco de Dados

- `make migrate-create name=nome_da_migration`: Cria uma nova migration
- `make migrate-up`: Executa todas as migrations pendentes
- `make migrate-down`: Reverte a última migration
- `make migrate-down-steps steps=N`: Reverte N migrations
- `make migrate-status`: Mostra o status das migrations
- `make migrate-force version=YYYYMMDDHHMMSS`: Força a versão da migration
