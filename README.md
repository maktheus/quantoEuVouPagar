# Quanto Eu Vou Pagar?

[![Go Report Card](https://goreportcard.com/badge/github.com/matheus-uchoa/quanto-eu-vou-pagar)](https://goreportcard.com/report/github.com/matheus-uchoa/quanto-eu-vou-pagar)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Test](https://github.com/matheus-uchoa/quanto-eu-vou-pagar/actions/workflows/test.yml/badge.svg)](https://github.com/matheus-uchoa/quanto-eu-vou-pagar/actions/workflows/test.yml)
[![Docker](https://github.com/matheus-uchoa/quanto-eu-vou-pagar/actions/workflows/docker.yml/badge.svg)](https://github.com/matheus-uchoa/quanto-eu-vou-pagar/actions/workflows/docker.yml)
[![Audit](https://github.com/matheus-uchoa/quanto-eu-vou-pagar/actions/workflows/audit.yml/badge.svg)](https://github.com/matheus-uchoa/quanto-eu-vou-pagar/actions/workflows/audit.yml)

## Descrição

**Quanto Eu Vou Pagar?** é um microserviço em Go projetado para ajudar você a gerenciar e prever os custos relacionados a grandes compras, como a aquisição de uma casa. Ele permite o cálculo detalhado de subcontas com diferentes tipos de juros (tratados internamente como juros mensais para simplificação), simulações de ganhos/perdas de juros, tracking de investimentos com aportes flexíveis, e gestão de dívidas com amortizações.

O objetivo é fornecer uma visão clara e precisa de quanto você vai pagar por mês, bem como simular cenários futuros (projeções) e gerar relatórios (incluindo PDFs) para um acompanhamento eficaz. O serviço também visa facilitar análises futuras com modelos de IA.

## Funcionalidades Implementadas

- **Cálculo de Subcontas**: Suporte a múltiplas subcontas com diferentes regimes de juros (Sistema Price e SAC).
- **Simulações de Juros**: Simule cenários de ganhos ou perdas com base em diferentes taxas de juros.
- **Tracking de Investimentos**: Estrutura criada (modelos, tabelas), lógica de negócio em desenvolvimento.
- **Gestão de Dívidas**: Estrutura criada (modelos, tabelas), lógica de negócio em desenvolvimento.
- **Visualização Temporal**: Ainda não implementada.
- **Exportação de Relatórios**: Ainda não implementada.
- **Facilitação para Análise de Dados/IA**: Estrutura do banco de dados preparada.
- **Sistema de Pagamento**: Estrutura criada (modelos, tabelas, handlers básicos), integração com gateways de pagamento em desenvolvimento.
- **Sistema de Autenticação e Autorização**: Sistema de login/registro com JWT e níveis de permissão (`user`, `admin`).
- **Sistema de Assinatura com Níveis**: Níveis de assinatura (Free, Tier, Gold, Platinum) com funcionalidades diferenciadas.

## Documentação

- **Documentação do Projeto**: Veja a documentação detalhada do projeto em [docs/README.md](docs/README.md).
- **Documentação da API**: Veja a documentação detalhada da API em [web/api-docs.md](web/api-docs.md).
- **Documentação Interativa da API**: A documentação interativa da API está disponível em `http://localhost:8080/swagger/index.html` (quando rodando localmente).

## Como Rodar

### Pré-requisitos

- [Docker](https://www.docker.com/products/docker-desktop) e [Docker Compose](https://docs.docker.com/compose/install/)

### Configuração das Variáveis de Ambiente

1. Crie um arquivo `.env` para desenvolvimento (opcional, valores padrão serão usados se não existir):
   ```bash
   cp .env.production .env
   # Edite .env para configurar as variáveis de ambiente para desenvolvimento
   ```

2. Para produção, edite `.env.production` com as variáveis de ambiente para produção.

### Executando com Docker Compose (Recomendado)

#### Desenvolvimento

1. **Clone o repositório**:
   ```bash
   git clone https://github.com/matheus-uchoa/quanto-eu-vou-pagar.git
   cd quanto-eu-vou-pagar
   ```

2. **Inicie os serviços**:
   ```bash
   docker-compose up -d
   ```
   Isso iniciará os seguintes containers:
   - `quanto-eu-vou-pagar-db`: PostgreSQL.
   - `quanto-eu-vou-pagar-app`: A aplicação Go.
   - `quanto-eu-vou-pagar-nginx`: Nginx (proxy reverso).
   - `quanto-eu-vou-pagar-prometheus`: Prometheus (monitoramento).
   - `quanto-eu-vou-pagar-grafana`: Grafana (visualização de métricas).
   - `quanto-eu-vou-pagar-fluentd`: Fluentd (logging).
   - `quanto-eu-vou-pagar-backup`: Serviço de backup automático.

3. **Acesse os serviços**:
   - **API**: `http://localhost` (através do Nginx)
   - **Swagger UI**: `http://localhost/swagger/index.html`
   - **Métricas Prometheus**: `http://localhost/metrics` (protegido, acessível apenas internamente)
   - **Prometheus**: `http://localhost:9090`
   - **Grafana**: `http://localhost:3000` (usuário: `admin`, senha: `admin`)
   - **Banco de Dados**: `localhost:5432` (se a porta for mapeada no docker-compose.yml)

#### Produção

1. **Edite `.env.production`** com as variáveis de ambiente para produção.
2. **Edite `nginx.conf`** para configurar o domínio e SSL/TLS (se necessário).
3. **Inicie os serviços**:
   ```bash
   docker-compose up -d
   ```

### Parando os Serviços

Para parar todos os serviços:

```bash
docker-compose down
```

Para parar todos os serviços e remover os volumes (incluindo dados do banco de dados):

```bash
docker-compose down -v
```

## Testes

O projeto utiliza testes unitários para garantir a qualidade e a confiabilidade do código. Os testes são escritos utilizando o pacote `testing` padrão do Go.

### Rodando os Testes Localmente

**Pré-requisitos**:
- [Go](https://golang.org/dl/) (versão 1.22 ou superior)

1. Clone o repositório.
2. Execute os testes:
   ```bash
   go test ./...
   ```

3. Para rodar os testes com cobertura:
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

### Rodando os Testes com Docker

1. Clone o repositório.
2. Execute os testes:
   ```bash
   docker-compose run --rm app go test ./...
   ```

## CI/CD

Este projeto utiliza GitHub Actions para integração contínua e entrega contínua.

### Workflows

- **Test**: Executado em cada push ou pull request para a branch `main`. Roda `go build` e `go test`.
  - Arquivo: `.github/workflows/test.yml`
- **Docker**: Executado em cada push ou pull request para a branch `main`. Faz o build da imagem Docker e a publica no Docker Hub (se não for um pull request).
  - Arquivo: `.github/workflows/docker.yml`
  - **Requer**: Secrets `DOCKER_USERNAME` e `DOCKER_PASSWORD` configurados no repositório.
- **Audit**: Executado em cada push ou pull request para a branch `main`. Roda `go vet` e `govulncheck` para auditoria de segurança.
  - Arquivo: `.github/workflows/audit.yml`
- **Release Please**: Executado em cada push para a branch `main`. Cria releases automatizadas.
  - Arquivo: `.github/workflows/release-please.yaml`
  - **Requer**: (Opcional) Secret `RELEASE_PLEASE_TOKEN` configurado no repositório.

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

## Como Contribuir

Este projeto está em fase inicial. Contribuições são bem-vindas! Por favor, abra uma issue para discutir mudanças importantes antes de enviar um pull request.

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.