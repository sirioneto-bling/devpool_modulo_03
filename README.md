# Módulo 03: Decomposição de Domínio em Microsserviço

Projeto educacional que demonstra a integracao entre uma API monolitica em PHP e um microsservico em Go, simulando um cenario real de decomposicao de dominio.

---

## Objetivo

No modulo anterior (modulo 2), implementamos uma API de pedidos em PHP onde os dados do cliente -- como nome e documento -- eram armazenados diretamente no registro da venda. Essa abordagem funciona para um MVP, mas gera duplicacao de dados, inconsistencias entre vendas e impede a evolucao independente do cadastro de clientes.

Neste modulo, extraimos o dominio de **Clientes** para um microsservico Go dedicado. A API PHP passa a gravar apenas o `id_cliente` na venda e consome o microsservico para todas as operacoes de consulta e manutencao de clientes -- incluindo busca por termo (lookup), carregamento por ID e filtro de vendas por cliente.

### Beneficios

- **Dominio isolado** -- o cadastro de clientes evolui, escala e deploya independentemente
- **Fonte unica de verdade** -- elimina duplicacao de dados de cliente entre tabelas
- **Tecnologia adequada** -- Go oferece alta performance e baixo consumo de recursos para um servico de consulta intensiva
- **Contrato claro** -- a comunicacao via REST com JSON forca um contrato explicito entre os sistemas

---

## Arquitetura

```mermaid
graph LR
    frontend["Vue.js (Frontend)"]

    subgraph phpApi ["PHP Web API (monolito)"]
        vendas["Vendas / Pedidos"]
    end

    subgraph goApi ["Go Microservice"]
        clientes["Clientes CRUD"]
    end

    subgraph databases ["Databases"]
        mysqlPhp["MySQL - PHP"]
        mysqlGo["MySQL - Go"]
    end

    blingApi["API Bling (externa)"]

    frontend -->|"HTTP"| vendas
    frontend -->|"HTTP REST"| blingApi
    vendas -->|"HTTP REST"| clientes
    vendas --> mysqlPhp
    clientes --> mysqlGo
```

A API PHP mantem seu banco MySQL proprio para pedidos. O microsservico Go possui seu banco isolado para clientes. Produtos sao consultados pelo frontend (Vue.js) diretamente na API externa do Bling, sem passar pelo webserver PHP. A comunicacao entre PHP e Go e feita via HTTP REST -- sem acoplamento a nivel de banco de dados.

---

## Fluxos Principais

### Busca de clientes por termo (autocomplete na venda)

Quando o usuario esta criando uma venda e digita parte do nome do cliente, a API PHP consulta o microsservico Go para exibir sugestoes.

```mermaid
sequenceDiagram
    participant User as Usuario
    participant PhpAPI as PHP Web API
    participant GoAPI as Go Microservice
    participant MySQL as MySQL Go

    User->>PhpAPI: Digita termo de busca
    PhpAPI->>GoAPI: GET /v1/clientes?q=termo
    GoAPI->>MySQL: SELECT WHERE nome LIKE '%termo%'
    MySQL-->>GoAPI: Lista de clientes
    GoAPI-->>PhpAPI: JSON com clientes encontrados
    PhpAPI-->>User: Exibe sugestoes para selecao
```

### Criacao de venda com vinculo de cliente (lookup + selecao)

Fluxo completo: o usuario busca o cliente por termo, seleciona da lista retornada, a API PHP carrega os dados completos e grava a venda vinculada ao ID do cliente.

```mermaid
sequenceDiagram
    participant User as Usuario
    participant PhpAPI as PHP Web API
    participant GoAPI as Go Microservice
    participant MySQL as MySQL Go

    User->>PhpAPI: Digita termo de busca do cliente
    PhpAPI->>GoAPI: GET /v1/clientes?q=termo
    GoAPI->>MySQL: SELECT WHERE nome LIKE '%termo%'
    MySQL-->>GoAPI: Lista de clientes
    GoAPI-->>PhpAPI: JSON com clientes encontrados
    PhpAPI-->>User: Exibe sugestoes para selecao

    User->>PhpAPI: Seleciona cliente (id=42)
    PhpAPI->>GoAPI: GET /v1/clientes/42
    GoAPI->>MySQL: SELECT WHERE id = 42
    MySQL-->>GoAPI: Dados do cliente
    GoAPI-->>PhpAPI: JSON com cliente completo
    PhpAPI->>PhpAPI: Grava venda com id_cliente = 42
    PhpAPI-->>User: Venda criada com sucesso
```

### Filtro de vendas por cliente (lookup + listagem)

O usuario busca o cliente por termo, seleciona da lista e a API PHP filtra as vendas vinculadas ao ID daquele cliente.

```mermaid
sequenceDiagram
    participant User as Usuario
    participant PhpAPI as PHP Web API
    participant GoAPI as Go Microservice
    participant MySQLPhp as MySQL PHP
    participant MySQLGo as MySQL Go

    User->>PhpAPI: Digita termo de busca do cliente
    PhpAPI->>GoAPI: GET /v1/clientes?q=termo
    GoAPI->>MySQLGo: SELECT WHERE nome LIKE '%termo%'
    MySQLGo-->>GoAPI: Lista de clientes
    GoAPI-->>PhpAPI: JSON com clientes encontrados
    PhpAPI-->>User: Exibe sugestoes para selecao

    User->>PhpAPI: Seleciona cliente (id=42)
    PhpAPI->>MySQLPhp: SELECT vendas WHERE id_cliente = 42
    MySQLPhp-->>PhpAPI: Lista de vendas
    PhpAPI->>GoAPI: GET /v1/clientes/42
    GoAPI->>MySQLGo: SELECT WHERE id = 42
    MySQLGo-->>GoAPI: Dados do cliente
    GoAPI-->>PhpAPI: JSON com nome e dados do cliente
    PhpAPI->>PhpAPI: Combina vendas + dados do cliente
    PhpAPI-->>User: Vendas do cliente exibidas
```

### CRUD de clientes no microsservico

O microsservico Go expoe um CRUD completo. O cadastro pode ser acessado tanto pela API PHP quanto diretamente.

```mermaid
sequenceDiagram
    participant Client as Cliente HTTP
    participant GoAPI as Go Microservice
    participant MySQL as MySQL Go

    Client->>GoAPI: POST /v1/clientes (criar)
    GoAPI->>MySQL: INSERT INTO clientes
    MySQL-->>GoAPI: ID gerado
    GoAPI-->>Client: 201 Created + JSON

    Client->>GoAPI: GET /v1/clientes (listar)
    GoAPI->>MySQL: SELECT * FROM clientes
    MySQL-->>GoAPI: Resultados
    GoAPI-->>Client: 200 OK + JSON

    Client->>GoAPI: PUT /v1/clientes/42 (atualizar)
    GoAPI->>MySQL: UPDATE clientes SET ...
    MySQL-->>GoAPI: Rows affected
    GoAPI-->>Client: 200 OK + JSON

    Client->>GoAPI: DELETE /v1/clientes/42 (remover)
    GoAPI->>MySQL: DELETE FROM clientes
    MySQL-->>GoAPI: Rows affected
    GoAPI-->>Client: 204 No Content
```

---

## Estrutura do Repositorio

```
devpool-base-web-api/
├── golang-web-api/          # Microsservico Go -- CRUD de clientes
│   ├── cmd/                 # Ponto de entrada da aplicacao
│   ├── internal/            # Camadas DDD (domain, app, infrastructure, presentation)
│   ├── scripts_db/          # Scripts DDL do banco
│   ├── docs/                # Collection Bruno para testes de API
│   ├── docker-compose.yml   # MySQL do microsservico
│   ├── Makefile             # Comandos de desenvolvimento
│   └── README.md            # Setup, arquitetura DDD e testes
│
├── simple-php-api/          # API PHP -- pedidos, produtos e integracao com clientes
│   └── README.md            # Setup e documentacao da API PHP
│
└── README.md                # Este arquivo -- visao geral e arquitetura
```

---

## Setup

Cada subprojeto possui seu proprio ambiente, dependencias e instrucoes de execucao:

| Subprojeto | Stack | Documentacao |
|------------|-------|--------------|
| **golang-web-api** | Go 1.24+, MySQL 8, Docker | [golang-web-api/README.md](golang-web-api/README.md) |
| **simple-php-api** | PHP, MySQL, Docker | [simple-php-api/README.md](simple-php-api/README.md) |

---

## Como rodar o fluxo completo (PHP -> Go)

Sequencia para subir os dois servicos integrados e testar a chamada PHP -> Go.

### 1. Subir o ambiente PHP (que ja levanta o MySQL compartilhado)

```bash
cd simple-php-api
cp .env.example .env
docker compose up -d
```

Isso sobe dois containers: `devpool-webserver` (PHP em `localhost:88`) e `devpool-mysql` (MySQL exposto em `localhost:3312`).

### 2. Garantir que a tabela `tasks` existe no MySQL

A API Go grava em `devpool_erp` (mesmo banco do PHP). Crie a tabela se ainda nao existir:

```bash
docker exec -i devpool-mysql mysql -uroot -pasdf000 devpool_erp \
  < ../golang-web-api/scripts_db/01-database.sql
```

### 3. Configurar o `.env` da API Go

Crie `golang-web-api/.env` apontando para o MySQL do container PHP:

```env
APP_NAME=devpool-base-web-api
API_PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=3312
DB_USER=root
DB_PASSWORD=asdf000
DB_NAME=devpool_erp
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
```

### 4. Subir a API Go via debug

Abra o workspace `devpool_modulo_03` no Cursor/VS Code, va em **Run and Debug** (Ctrl+Shift+D), selecione **"Go: web_api (debug)"** e pressione **F5**.

A API Go fica disponivel em `http://localhost:8080`.

> A configuracao de debug ja esta versionada em `.vscode/launch.json`. Pre-requisitos: Go 1.24+ instalado e Delve (`go install github.com/go-delve/delve/cmd/dlv@latest`).

### 5. Testar a integracao PHP -> Go

```bash
# Cria task via PHP, que internamente chama a API Go
curl -X POST http://localhost:88/api/task/store \
  -H "Content-Type: application/json" \
  -d '{"title": "Estudar Go", "description": "Completar o tour of Go"}'

# Lista tasks via PHP
curl http://localhost:88/api/task/index
```

### Como o PHP enxerga a API Go

O container do PHP usa o DNS especial `host.docker.internal` para acessar a maquina host (onde a API Go roda em modo debug). Isso ja esta configurado em `simple-php-api/docker-compose.yml` via `extra_hosts: ["host.docker.internal:host-gateway"]`.

> Se voce alterar o `docker-compose.yml`, recrie o container com `docker compose up -d` (o `restart` nao aplica mudancas no compose).

---

## Comunicacao entre os servicos

| Origem | Destino | Metodo | Endpoint | Descricao |
|--------|---------|--------|----------|-----------|
| PHP API | Go Microservice | `GET` | `/v1/clientes?q={termo}` | Busca clientes por nome (autocomplete) |
| PHP API | Go Microservice | `GET` | `/v1/clientes/{id}` | Carrega dados completos do cliente |
| PHP API | Go Microservice | `POST` | `/v1/clientes` | Cadastra novo cliente |
| PHP API | Go Microservice | `PUT` | `/v1/clientes/{id}` | Atualiza dados do cliente |
| PHP API | Go Microservice | `DELETE` | `/v1/clientes/{id}` | Remove cliente |
