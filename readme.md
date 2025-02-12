# Neoway

## Descrição
Aplicação REST para cadastrar e listar clientes de uma empresa fictícia. Essa empresa em questão fornece serviços B2C e B2B, ou seja, teremos cadastrosde clientes como pessoas físicas e como pessoas jurídicas

## Funcionalidades
- Cadastro de clientes PF (CPF) e PJ (CNPJ)
- Listagem de clientes
- Consulta de cliente por documento
- Bloqueio/desbloqueio de clientes (restrição?) 
- Métricas de requisições e uptime do sistema

## Requisitos
- Docker e Docker Compose
- Go 1.23 ou superior (para desenvolvimento)
- PostgreSQL (gerenciado via Docker)

## Executar
- make build
- make run

## Testar
- make test-all

## Endpoints
Aplicação inicia em localhost:8080, os endpoints são:

**POST /clients**

Ao realizar o cadastro de um cliente, o sistema irá validar se o documento é válido e se o tipo de cliente corresponde ao formato do documento. 
(CPF precisa ter os 11 dígitos válidos conforme calculo do dígito verificador, e CNPJ precisa ter os 14 dígitos válidos conforme calculo do dígito verificador tambem)
A validação é realiza utilizando o package util

```sh
curl -X POST http://localhost:8080/clients \
-H "Content-Type: application/json" \
-d '{
    "name": "Gionon Wright",
    "document": "018.729.949-89",
    "type": "PERSON"
}'
```

Resposta (201 Created)
```json
{
"id": 1,
"name": "Nome do Cliente",
"document": "12345678901",
"type": "PERSON",
"blocked": false,
"created_at": "2024-03-20T10:00:00Z",
"updated_at": "2024-03-20T10:00:00Z"
}
```

**GET /clients**

Listará todos os clientes cadastrados, independente de estado ou tipo, em ordem alfabética, conforme query realizada no banco.
```sql
	result := r.db.Order("name ASC").Find(&clients)
```
```sh
curl http://localhost:8080/clients
```

Resposta (200 OK)
```json
[
{
"id": 1,
"name": "Nome do Cliente",
"document": "12345678901",
"type": "PERSON",
"blocked": false,
"created_at": "2024-03-20T10:00:00Z",
"updated_at": "2024-03-20T10:00:00Z"
}
```
]
### Consulta por Documento

Consulta por documento irá funcionar utilizando tanto o documento com pontos, vírgula e barra, quanto
o número sem caracteres. 

**GET /clients/{document}**

```sh
curl http://localhost:8080/clients/12345678901
```
ou 
```sh
curl http://localhost:8080/clients/123.456.789-01
```


Resposta (200 OK)
```json
{
"id": 1,
"name": "Nome do Cliente",
"document": "12345678901",
"type": "PERSON",
"blocked": false,
"created_at": "2024-03-20T10:00:00Z",
"updated_at": "2024-03-20T10:00:00Z"
}
```

### Bloqueio de Cliente

Booleano básico que valida se o cliente está bloqueado ou não, é possível bloquear e desbloquear um cliente sem nenhuma restrição.

```http
PUT /clients/{document}/block
Exemplo: PUT /clients/12345678901/block
```

Resposta (200 OK)
```json
{
"message": "Cliente bloqueado com sucesso"
}
```

```http
PUT /clients/{document}/unblock
Exemplo: PUT /clients/12345678901/unblock
```
Resposta (200 OK)
```json
{
"message": "Cliente desbloqueado com sucesso"
}
```

### Métricas do Sistema

**Uptime**

Métricas básicas do docker, irá calcular o uptime conforme a data de início do container e a data atual.

```http
GET /metrics/uptime
```

Resposta (200 OK)
```json
{
"uptime": "24h 13m 5s",
"start_time": "2024-03-19T10:00:00Z",
"last_restart": "2024-03-20T00:00:00Z",
"uptime_string": "1 dia, 0 horas, 13 minutos e 5 segundos"
}
```
**Requests**

Métricas de requisições, listara todas as reqs realizadas, contabilizando quantas vezes cada endpoint foi acessado separadamente.

```http
GET /metrics/requests
```

Resposta (200 OK)
```json
[
{
"method": "POST",
"path": "/clients",
"count": 150,
"created_at": "2024-03-19T10:00:00Z",
"updated_at": "2024-03-20T10:00:00Z"
}
]
```

### Códigos de Erro
- `400 Bad Request`: Requisição inválida (dados incorretos)
- `404 Not Found`: Recurso não encontrado
- `409 Conflict`: Conflito (ex: documento já cadastrado)
- `500 Internal Server Error`: Erro interno do servidor

### Validações
- CPF deve conter 11 dígitos numéricos e calcúlo deve estar correto;
- CNPJ deve conter 14 dígitos numéricos e calcúlo deve estar correto;
- O tipo de cliente deve corresponder ao formato do documento;
- Documentos são únicos no sistema;
- Nome do cliente e docs são obrigatórios;

### Arquitetura e Patterns

Baseada em clean architecture, utilizando o padrão de camadas de domínio, aplicação, infra e interface.

```sh
internal/
  ├── domain/         # Regras de negócio
  ├── application/    # use cases
  ├── infrastructure/ # Implementações externas (DB, etc)
  └── interfaces/     # Res e respostas
```
**Design Patterns**

Utilização de repository pattern para acesso ao db, dependency injection para lidar com as dependencias de cada camada, separação da interface e service layer para facilitar a testabilidade e manutenção do código.

**Database**
Utilizado um banco postgres tanto para o desenvolimento da aplicação principal quanto para os testes,
para facilitar execuções e complexidade, utilizei GORM como ORM.

A ideia foi utilizar postgres pela familariedade que possuo, porém ao implementar por exemplo o repository pattern para abstrair a interface é possível substituir o db para mongodb com algumas alterações, visto que para os modelos utilizei tags específicas para o GORM.

**Testes**

Realizados com intuito de cobirr os principais pontos da aplicação: Validação dos dados, regras de negócio, endpoints e operações no db.

A validação de dados testa a validação de formato dos documentos (CPF e CNPJ) e a validação de dados obrigatórios (nome e documento), além de validar o tipo de cliente (PERSON ou BUSINESS) e a unicidade do documento.

Já a validação das regras de negócio, testa a validação de documento, bloqueio e desbloqueio de cliente, além de testar a listagem de clientes e a consulta por documento, além dos erros em diferentes cenários.

Os endepoits são testados através de reqs de validação, de formato da resposta, respostas para os erros e status HTTP esperados.

Por fim, as operações no db são testadas através de testes de integração, utilizando o mock para simular o db e testar a lógiga de operações CRUD, como ele está se comportando nas transações e a execução das queries.
