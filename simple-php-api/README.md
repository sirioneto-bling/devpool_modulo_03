# API REST em PHP com padrão MVC

Este é um projeto de estudos para iniciantes que implementa uma **API REST** utilizando o padrão **MVC (Model-View-Controller)** em PHP puro, configurado para rodar em Docker.

## O que é uma API?

**API (Application Programming Interface)** é uma interface que permite a comunicação entre diferentes sistemas. Uma **API REST** utiliza o protocolo HTTP para receber requisições e retornar dados, geralmente no formato JSON.

Por exemplo, quando você acessa `/api/usuarios`, a API processa a requisição e retorna uma lista de usuários em formato JSON:

```json
{
    "status": "success",
    "data": [
        { "id": 1, "nome": "João" },
        { "id": 2, "nome": "Maria" }
    ]
}
```

## O que é o padrão MVC?

**MVC (Model-View-Controller)** é um padrão de arquitetura que separa a aplicação em três camadas:

| Camada | Responsabilidade | Neste projeto |
|--------|------------------|---------------|
| **Model** | Representa os dados e a lógica de acesso ao banco de dados | `app/Models/` |
| **View** | Apresenta os dados ao usuário (em APIs, é a resposta JSON) | Resposta JSON |
| **Controller** | Recebe as requisições, processa e coordena Model e View | `app/Controllers/` |

### Fluxo de uma requisição

```
Requisição HTTP → Controller → Model (busca dados) → Controller → Resposta JSON
```

**Exemplo prático:**

1. Usuário faz requisição `GET /api/exemplo/show/1`
2. O **Controller** (`ExemploController`) recebe a requisição
3. O Controller chama o **Model** (`Exemplo`) para buscar o registro com ID 1
4. O Model consulta o banco de dados e retorna os dados
5. O Controller formata e retorna a resposta JSON

## Requisitos

Antes de rodar o projeto, é necessário ter instalado:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Como Executar

### 1. Preparando o ambiente

Faça uma cópia do arquivo `.env.example` e renomeie para `.env`:

```bash
cp .env.example .env
```

### 2. Iniciar o Projeto

Execute o comando para iniciar os containers:

```bash
docker-compose up
```

### 3. Acessando a API

Após o build ser concluído, a API estará disponível em:

```
http://localhost:88
```

## Estrutura de Diretórios

```
projeto-mvc-php/
├── .docker/                # Configurações do Docker
├── app/                    # Código principal da aplicação
│   ├── Client/             # Clients para APIs externas
│   │   ├── BaseClient.php  # Classe abstrata (HTTP genérico com Guzzle)
│   │   ├── NominatimClient.php # Client do Nominatim (geocodificação)
│   │   └── TaskClient.php  # Client da API Go de tasks
│   ├── Controllers/        # Controllers da API
│   │   └── Api/            # Controllers de endpoints
│   ├── Core/               # Classes base (Controller, Core, Model)
│   ├── Middleware/         # Middlewares de autenticação
│   ├── Models/             # Models (entidades do banco)
│   └── Supports/           # Classes de suporte (Criptografia, Logs, etc)
├── config/                 # Configurações da aplicação
│   └── config.php          # Arquivo principal de configuração
├── docker-compose.yml      # Configuração do Docker Compose
├── index.php               # Ponto de entrada da aplicação
└── README.md               # Este arquivo
```

## Como Funciona o Roteamento

O projeto implementa um sistema de roteamento automático baseado em namespaces. A URL é mapeada diretamente para os controllers e métodos.

### Estrutura das Rotas

```
http://localhost:88/api/{controller}/{método}/{parâmetros}
```

### Exemplos

| URL | Controller | Método | Descrição |
|-----|------------|--------|-----------|
| `/api/exemplo` | ExemploController | index() | Lista todos |
| `/api/exemplo/show/1` | ExemploController | show(1) | Busca por ID |
| `/api/exemplo/store` | ExemploController | store() | Cria novo (POST) |
| `/api/exemplo/update/1` | ExemploController | update(1) | Atualiza (PUT) |
| `/api/exemplo/delete/1` | ExemploController | delete(1) | Remove (DELETE) |

## Criando um Novo Controller

Para criar um novo endpoint, basta criar um controller na pasta `app/Controllers/Api/`.

Exemplo: Criar um `UserController`:

```php
<?php

namespace App\Controllers\Api;

use App\Core\Controller;

class UserController extends Controller
{
    public function index()
    {
        return $this->jsonResponse(['message' => 'Lista de usuários']);
    }

    public function show(int $id)
    {
        return $this->jsonResponse(['message' => 'Usuário ' . $id]);
    }

    public function store()
    {
        $this->validateRequestMethods(['POST']);
        $data = $this->getRequestData();
        
        return $this->jsonResponse(['message' => 'Usuário criado', 'data' => $data]);
    }
}
```

Após criar o controller, as rotas já estarão disponíveis automaticamente:

- `GET /api/user` → `UserController::index()`
- `GET /api/user/show/1` → `UserController::show(1)`
- `POST /api/user/store` → `UserController::store()`

## Usando Middlewares

**Middlewares** são filtros que executam **antes** da ação do controller. São úteis para:

- Verificar se o usuário está autenticado
- Validar permissões de acesso
- Registrar logs de requisições

### Como usar um Middleware

No construtor do controller, chame o método `middleware()` passando a classe:

```php
<?php

namespace App\Controllers\Api;

use App\Core\Controller;
use App\Middleware\AuthApiMiddleware;

class MeuController extends Controller
{
    public function __construct()
    {
        // Executa o middleware antes de qualquer ação
        $this->middleware(AuthApiMiddleware::class);
    }

    public function index()
    {
        // Só chega aqui se o middleware permitir
        return $this->jsonResponse(['message' => 'Usuário autenticado!']);
    }
}
```

### Como criar um Middleware

Crie uma classe em `app/Middleware/` que estenda `Middleware`:

```php
<?php

namespace App\Middleware;

use App\Core\Middleware;

class MeuMiddleware extends Middleware
{
    public function handle()
    {
        // Sua lógica de validação aqui
        if (!$this->usuarioTemPermissao()) {
            return $this->jsonResponse([], 'Acesso negado', 403);
        }
        
        // Se não retornar nada, a requisição continua normalmente
    }

    private function usuarioTemPermissao()
    {
        // Implemente sua lógica
        return true;
    }
}
```

### Fluxo com Middleware

```
Requisição HTTP → Middleware (valida) → Controller → Resposta JSON
                      ↓
              Se falhar, retorna erro
```

## Autenticação

O projeto inclui um sistema de autenticação por **sessão** para proteger rotas. O `ExemploController` utiliza o middleware `AuthApiMiddleware` que verifica se o usuário está autenticado.

### Credenciais padrão

| Campo | Valor |
|-------|-------|
| Usuário | `devpool` |
| Senha | `asdf000` |

### Como fazer login

Faça uma requisição POST para `/api/auth/login`:

```bash
curl --location 'http://localhost:88/api/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "devpool",
    "password": "asdf000"
}'
```

**Resposta de sucesso:**

```json
{
    "data": {
        "email": "devpool@mail.com",
        "name": "DevPool",
        "token": "1234567890"
    },
    "message": "Autenticação efetuada com sucesso"
}
```

### Como fazer logout

```bash
curl --location --request POST 'http://localhost:88/api/auth/logout'
```

### Testando rotas protegidas

Após fazer login, a sessão fica ativa e você pode acessar as rotas protegidas:

```bash
# Listar todos os registros (rota protegida)
curl --location 'http://localhost:88/api/exemplo'
```

**Importante:** Se estiver usando ferramentas como Postman ou Insomnia, certifique-se de que os cookies estão sendo enviados nas requisições para manter a sessão ativa.

## Métodos Úteis do Controller Base

| Método | Descrição |
|--------|-----------|
| `$this->jsonResponse($data, $message, $status)` | Retorna resposta JSON |
| `$this->getRequestData()` | Obtém dados do body da requisição |
| `$this->validateRequestMethods(['GET', 'POST'])` | Valida método HTTP permitido |
| `$this->middleware(MeuMiddleware::class)` | Executa um middleware antes da ação |

## Busca de Endereço (Nominatim)

O projeto inclui uma integração com o [Nominatim](https://nominatim.openstreetmap.org/) (OpenStreetMap) para busca e geocodificação de endereços. Não requer autenticação.

### Endpoints disponíveis

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/searchAddress/search` | Busca endereço por texto |
| `GET` | `/api/searchAddress/reverse` | Geocodificação reversa (coordenadas → endereço) |
| `GET` | `/api/searchAddress/lookup` | Busca por OSM IDs |

### Buscar endereço por texto

```bash
curl 'http://localhost:88/api/searchAddress/search?q=Rua+Augusta,+São+Paulo&limit=5'
```

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `q` | string | Sim | Texto do endereço a buscar |
| `limit` | int | Não | Máximo de resultados (padrão: 10) |

### Geocodificação reversa

```bash
curl 'http://localhost:88/api/searchAddress/reverse?lat=-23.5614&lon=-46.6558'
```

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `lat` | float | Sim | Latitude |
| `lon` | float | Sim | Longitude |

### Buscar por OSM IDs

```bash
curl 'http://localhost:88/api/searchAddress/lookup?osm_ids=R146656,W104393803'
```

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `osm_ids` | string | Sim | IDs do OpenStreetMap separados por vírgula (ex: R146656,W104393803) |

### Arquitetura do Client

A integração com APIs externas segue uma arquitetura em camadas reutilizável:

```
Controller → Client concreto → BaseClient (Guzzle) → API externa
```

- **`BaseClient`** (`app/Client/BaseClient.php`): Classe abstrata que encapsula o Guzzle HTTP Client com métodos `get()`, `post()`, `put()`, `delete()` e tratamento de erros padronizado.
- **`NominatimClient`** (`app/Client/NominatimClient.php`): Client concreto que implementa os endpoints do Nominatim.
- **`TaskClient`** (`app/Client/TaskClient.php`): Client concreto que consome a API Go de tasks.

Para integrar uma nova API externa, basta criar um novo client estendendo `BaseClient`:

```php
<?php

namespace App\Client;

class MeuNovoClient extends BaseClient
{
    protected string $baseUrl = 'https://api.exemplo.com';
    protected int $timeout = 10;

    public function buscarDados(string $param): array
    {
        return $this->get('/endpoint', ['param' => $param]);
    }
}
```

## Tasks (API Go)

O projeto inclui uma integração com a **API Go de tasks** (`golang-web-api`), que roda localmente e compartilha o mesmo banco de dados MySQL. Isso demonstra a comunicação entre dois serviços em linguagens diferentes (PHP → Go) via HTTP. Não requer autenticação.

> **Pré-requisito:** A API Go deve estar rodando em `http://localhost:8080` (via debug no Cursor ou `go run`).

### Endpoints disponíveis

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/task/index` | Lista todas as tasks |
| `POST` | `/api/task/store` | Cria uma nova task |

### Listar tasks

```bash
curl http://localhost:88/api/task/index
```

**Resposta de sucesso:**

```json
{
    "data": [
        {
            "id": 1,
            "title": "Estudar Go",
            "description": "Completar o tour of Go",
            "status": "pending",
            "created_at": "2026-04-22T12:00:00Z",
            "updated_at": "2026-04-22T12:00:00Z"
        }
    ],
    "message": "Tasks listadas com sucesso"
}
```

### Criar task

```bash
curl -X POST http://localhost:88/api/task/store \
  -H "Content-Type: application/json" \
  -d '{"title": "Estudar Go", "description": "Completar o tour of Go"}'
```

**Parâmetros (body JSON):**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `title` | string | Sim | Título da task |
| `description` | string | Não | Descrição da task |

**Resposta de sucesso (201):**

```json
{
    "data": {
        "id": 1,
        "title": "Estudar Go",
        "description": "Completar o tour of Go",
        "status": "pending",
        "created_at": "2026-04-22T12:00:00Z",
        "updated_at": "2026-04-22T12:00:00Z"
    },
    "message": "Task criada com sucesso"
}
```

### Arquitetura da integração

```
Requisição HTTP → PHP (TaskController) → TaskClient (Guzzle) → API Go (:8080) → MySQL
```

Os dois projetos compartilham o mesmo banco de dados (`devpool_erp`) no container MySQL do Docker. O PHP acessa via rede interna do Docker (`devpool-mysql:3306`) e o Go acessa pela porta exposta na máquina host (`localhost:3312`).

## Configuração

O arquivo `config/config.php` contém as configurações principais:

```php
define('DEFAULT_CONTROLLER', 'Api\\Exemplo');  // Controller padrão
define('URL_BASE', 'http://localhost:88');     // URL base da API
define('MAINTENANCE', 0);                       // Modo manutenção (0 = desligado)
```

## Dicas para Iniciantes

1. **Comece pelo ExemploController**: Analise o arquivo `app/Controllers/Api/ExemploController.php` para entender como funciona um CRUD completo.

2. **Entenda o Model**: Veja `app/Models/Exemplo.php` para entender como interagir com o banco de dados.

3. **Use o Postman ou Insomnia**: Ferramentas como Postman facilitam testar os endpoints da API.

4. **Leia os arquivos Core**: Os arquivos em `app/Core/` são a base do framework e ajudam a entender como tudo funciona.

---

Qualquer dúvida, estou à disposição! :)
