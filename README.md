# Desafio To Do List API

Este projeto é uma API RESTful de gerenciamento de tarefas (To Do List) desenvolvida em **Go** com **MongoDB**, como parte do processo seletivo para a vaga de desenvolvedora backend jr.  
O projeto implementa operações completas de **CRUD** , validações de regras de negócio e documentação automática da API via **Swagger/OpenAPI**.

---

## Tecnologias
- Go
- MongoDB
- Gorilla Mux
- Swagger (swaggo)

---

## Pré-requisitos 

Antes de executar o projeto é necessário ter instalado:

- **Go 1.22 ou superior** 
- **MongoDB** rodando localmente em:  `mongodb://localhost:27017`

Configuração usada no projeto:
- Database: `taskdb`
- Collection: `tasks`
- Porta da API: `8080`

---

## Dependências principais 
As dependências são gerenciadas via **Go Modules (go.mod)**.

Principais libs utilizadas:

- `github.com/google/uuid` **v1.6.0**
- `github.com/gorilla/mux` **v1.8.1**
- `go.mongodb.org/mongo-driver/v2` **v2.5.0**
- `github.com/swaggo/swag` **v1.16.6**
- `github.com/swaggo/http-swagger/v2` **v2.0.2**

> As dependências indiretas são resolvidas automaticamente pelo Go e estão listadas no `go.mod` como `// indirect`.

---

## Como executar o projeto:

### 1) Garantir que o MongoDB está rodando
A API espera o Mongo em `mongodb://localhost:27017`.

Formas de confirmar:
- **MongoDB Compass**: conectar em `mongodb://localhost:27017`
- **Terminal** (se tiver `mongosh`): executar `mongosh` e ver se abre sem erro

### 2) Instalar dependências
Na raiz do projeto execute:
```bash
go mod tidy
```
### 3) Rodar a API

Execute:
```bash
go run main.go
```

Se tudo estiver correto, deve aparecer o seguinte no terminal:

Conectado ao MongoDB com sucesso
Servidor rodando na porta 8080

Dessa forma a API está disponível em:

`http://localhost:8080`

> Observação importante: mantenha o terminal aberto enquanto estiver testando a API.

---

## Como testar a API:

Opção principal: Insomnia

Com a API rodando, abrir o Insomnia e utilizar a seuguinte base URL:
`http://localhost:8080`

## Endpoints da API

| Método | Endpoint | Descrição |
|------|------|------|
| POST | /tasks | Criar nova tarefa |
| GET | /tasks | Listar todas as tarefas |
| GET | /tasks/{id} | Buscar tarefa por ID |
| PUT | /tasks/{id} | Atualizar tarefa |
| DELETE | /tasks/{id} | Remover tarefa |

### Fluxo recomendado de testes:

#### Criar tarefa (POST /tasks)

Body de exemplo:
```json
{
  "title": "Estudar Golang",
  "description": "Revisar conceitos de goroutines",
  "priority": "high",
  "due_date": "2026-04-10"
}
```

Resposta esperada: 
```
201 Created
```
#### Listar tarefas (GET /tasks)

Resposta esperada:
```
 200 OK
```
#### Buscar tarefa por ID (GET /tasks/{id})

Substitua {id} pelo identificador retornado na criação da tarefa.

Resposta esperada:
```
 200 OK
```
#### Atualizar tarefa (PUT /tasks/{id})

Body de exemplo:
```json
{
  "status": "in_progress"
}
```
Resposta esperada:
```
 200 OK
```

#### Deletar tarefa (DELETE /tasks/{id})

Resposta esperada:
```
 204 No Content
```
---

#### Observação: A API também possui documentação interativa utilizando Swagger.

Com a aplicação rodando, acessar no navegador:

`http://localhost:8080/swagger/index.html`

> O Swagger permite visualizar todos os endpoints e executar requisições diretamente pela interface web.

---

#### Regras de negócio implementadas

Status permitidos:
```
pending
in_progress
completed
cancelled
```
Prioridades permitidas:
```
low
medium
high
```
---

#### Validações implementadas:

- título obrigatório (3 a 100 caracteres)
- prioridade válida
- status válido
- data de vencimento não pode estar no passado
- tasks com status "completed" não podem ser editadas

Códigos de resposta HTTP utilizados
- 200 → sucesso
- 201 → criado
- 204 → removido
- 400 → erro de validação
- 404 → recurso não encontrado
- 500 → erro interno

