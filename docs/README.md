# Documentação do Projeto - Quanto Eu Vou Pagar?

## Visão Geral

Esta documentação fornece uma visão geral da arquitetura, estrutura do código e guias para desenvolvimento e manutenção do projeto "Quanto Eu Vou Pagar?".

## Arquitetura

O projeto segue uma arquitetura monolítica com uma API RESTful escrita em Go. A estrutura do código é organizada dentro do diretório `internal/` seguindo princípios de separação de responsabilidades.

### Componentes Principais

- **API**: Interface RESTful para interação com o serviço.
- **Banco de Dados**: PostgreSQL para armazenamento de dados.
- **Autenticação**: JWT (JSON Web Tokens) para autenticação de usuários.
- **Autorização**: Níveis de permissão (`user`, `admin`) e níveis de assinatura (`free`, `tier`, `gold`, `platinum`).
- **Monitoramento**: Métricas exportadas para Prometheus.
- **Logging**: Centralizado com Fluentd.
- **Documentação**: Documentação da API gerada automaticamente com Swagger.

### Estrutura de Diretórios

```
quanto-eu-vou-pagar/
├── cmd/
│   └── server/          # Ponto de entrada da aplicação
├── internal/            # Código fonte principal
│   ├── auth/            # Lógica de autenticação
│   ├── config/          # Configuração da aplicação
│   ├── database/        # Conexão com o banco de dados
│   ├── handler/         # Handlers HTTP
│   ├── logger/          # Configuração de logging
│   ├── middleware/      # Middlewares
│   ├── model/           # Modelos de dados
│   ├── repository/      # Acesso ao banco de dados
│   ├── service/         # Lógica de negócio
│   ├── subscription/    # Lógica de assinaturas
│   └── utils/           # Funções utilitárias
├── scripts/             # Scripts auxiliares
├── web/                 # Documentação da API
├── .github/workflows/   # Workflows do GitHub Actions
├── docs/                # Documentação do projeto
├── Dockerfile           # Dockerfile para build da imagem
├── docker-compose.yml   # Docker Compose para desenvolvimento e produção
├── nginx.conf           # Configuração do Nginx
├── prometheus.yml       # Configuração do Prometheus
├── fluentd.conf         # Configuração do Fluentd
├── .env.production      # Variáveis de ambiente para produção
├── go.mod               # Dependências do Go
├── go.sum               # Soma de verificação das dependências
├── Makefile             # Comandos úteis
├── README.md            # README principal
└── LICENSE              # Licença do projeto
```

## Desenvolvimento

### Pré-requisitos

- [Go](https://golang.org/dl/) (versão 1.22 ou superior)
- [Docker](https://www.docker.com/products/docker-desktop) e [Docker Compose](https://docs.docker.com/compose/install/)

### Configuração do Ambiente

1. Clone o repositório:
   ```bash
   git clone https://github.com/matheus-uchoa/quanto-eu-vou-pagar.git
   cd quanto-eu-vou-pagar
   ```

2. Crie um arquivo `.env` para desenvolvimento (opcional, valores padrão serão usados se não existir):
   ```bash
   cp .env.production .env
   # Edite .env para configurar as variáveis de ambiente para desenvolvimento
   ```

3. Inicie os serviços dependentes (banco de dados):
   ```bash
   docker-compose up -d db
   ```

4. (Opcional) Instale as dependências do Go:
   ```bash
   go mod tidy
   ```

### Executando a Aplicação

#### Localmente

```bash
go run cmd/server/main.go
```

Ou, compile e execute:

```bash
go build -o quanto-eu-vou-pagar cmd/server/main.go
./quanto-eu-vou-pagar
```

A API estará disponível em `http://localhost:8080`.

#### Com Docker Compose (Desenvolvimento)

```bash
docker-compose up -d
```

A API estará disponível em `http://localhost:8080`.

### Testes

Para rodar os testes:

```bash
go test ./...
```

Para rodar os testes com cobertura:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### CI/CD

O projeto utiliza GitHub Actions para CI/CD.

- **Test**: Executado em cada push ou pull request para a branch `main`. Roda `go build` e `go test`.
- **Docker**: Executado em cada push ou pull request para a branch `main`. Faz o build da imagem Docker e a publica no Docker Hub (se não for um pull request).
- **Audit**: Executado em cada push ou pull request para a branch `main`. Roda `go vet` e `govulncheck` para auditoria de segurança.
- **Release Please**: Executado em cada push para a branch `main`. Cria releases automatizadas.

## API

Veja a documentação detalhada da API em [web/api-docs.md](../web/api-docs.md).

A documentação interativa da API está disponível em:
- **URL**: `http://localhost:8080/swagger/index.html` (quando rodando localmente)

## Banco de Dados

### Esquema

O esquema do banco de dados é definido em `scripts/init.sql`.

### Migrations

Atualmente, o projeto não utiliza um sistema de migrations. O esquema é criado a partir do script `init.sql` quando o container do PostgreSQL é iniciado.

Para futuras alterações no esquema, considere a implementação de um sistema de migrations.

## Monitoramento

As métricas da aplicação podem ser coletadas pelo Prometheus no endpoint:
- **Endpoint**: `http://localhost:8080/metrics` (quando rodando localmente)

O projeto inclui configurações para Prometheus e Grafana. Para acessar o Grafana:
- **URL**: `http://localhost:3000`
- **Usuário**: `admin`
- **Senha**: `admin` (deve ser alterada em produção)

## Segurança

- As senhas são armazenadas como hashes bcrypt.
- Os tokens JWT são assinados com uma chave secreta.
- A chave secreta deve ser configurada como uma variável de ambiente em produção.

## Deployment

### Com Docker Compose (Produção)

1. Crie um arquivo `.env.production` com as variáveis de ambiente para produção (veja `.env.production` como exemplo).
2. Edite `docker-compose.yml` para configurar volumes, networks e outros parâmetros para produção.
3. Edite `nginx.conf` para configurar o domínio e SSL/TLS (se necessário).
4. Execute:
   ```bash
   docker-compose up -d
   ```

A API estará disponível em `http://localhost` (através do Nginx).

### Configurações de Produção

#### Variáveis de Ambiente

As variáveis de ambiente são carregadas de `.env.production` quando `ENV=production`.

#### Proxy Reverso (Nginx)

O Nginx é usado como proxy reverso para servir a API. A configuração está em `nginx.conf`.

#### SSL/TLS

Para configurar SSL/TLS, edite `nginx.conf` para usar certificados SSL/TLS. Isso requer um domínio válido e acesso externo ao servidor.

#### Backups

Um serviço de backup é configurado para fazer backup do banco de dados diariamente às 02:00. Os backups são armazenados no diretório `./backups` e os backups mais antigos que 7 dias são excluídos.

#### Monitoramento e Logging

O projeto inclui configurações para Prometheus, Grafana e Fluentd.

- **Prometheus**: Coleta métricas da aplicação.
- **Grafana**: Visualiza métricas coletadas pelo Prometheus.
- **Fluentd**: Centraliza logs de todos os serviços.

## Como Testar a API

A API pode ser testada usando ferramentas como `curl`, `Postman` ou `Swagger UI`.

### Endpoints Principais

1. **Health Check**:
   ```bash
   curl http://localhost/health
   ```

2. **Registro de Usuário**:
   ```bash
   curl -X POST http://localhost/auth/register \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}'
   ```

3. **Login de Usuário**:
   ```bash
   curl -X POST http://localhost/auth/login \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "password": "password123"}'
   ```
   Anote o `token` retornado.

4. **Criar Subconta** (requer token):
   ```bash
   curl -X POST http://localhost/api/subcontas \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer SEU_TOKEN_JWT" \
        -d '{
          "descricao": "Financiamento casa",
          "valor": 300000.0,
          "taxa_juros": 0.01,
          "parcelas": 120,
          "data_inicio": "2024-01-01T00:00:00Z",
          "sistema": "price",
          "tipo": "financiamento"
        }'
   ```

5. **Obter Parcelas de uma Subconta** (requer token):
   ```bash
   curl http://localhost/api/subcontas/1/parcelas \
        -H "Authorization: Bearer SEU_TOKEN_JWT"
   ```

## Contribuindo

1. Faça um fork do projeto.
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`).
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`).
4. Push para a branch (`git push origin feature/AmazingFeature`).
5. Abra um Pull Request.

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](../LICENSE) para mais detalhes.