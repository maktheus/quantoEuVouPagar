# Documentação da API - Quanto Eu Vou Pagar?

## Introdução

Esta documentação descreve os endpoints disponíveis na API do serviço "Quanto Eu Vou Pagar?".

## Servidor

O servidor está rodando em `http://localhost:8080` (quando rodando localmente com `docker-compose`).

## Documentação Swagger

A documentação interativa da API está disponível em:
- **URL**: `http://localhost:8080/swagger/index.html`

## Métricas com Prometheus

As métricas da aplicação podem ser coletadas pelo Prometheus no endpoint:
- **Endpoint**: `http://localhost:8080/metrics`

### Métricas Disponíveis

- `http_requests_total`: Contador para o número total de requisições HTTP, com labels para método, endpoint e código de status.
- `http_request_duration_seconds`: Histograma para a duração das requisições HTTP, com labels para método e endpoint.

## Autenticação

A API utiliza tokens JWT (JSON Web Tokens) para autenticação.

### Como obter um token

1. Registre um novo usuário em `POST /auth/register`.
2. Faça login em `POST /auth/login` com as credenciais do usuário.
3. A resposta do login conterá um `token`. Use este token nas requisições protegidas.

### Como usar o token

Inclua o token no header `Authorization` das suas requisições:

```
Authorization: Bearer SEU_TOKEN_JWT
```

## Permissões

Os endpoints da API são protegidos por níveis de permissão:

- **`user`**: Permissão básica para usuários registrados.
- **`admin`**: Permissão elevada para administradores.

Alguns endpoints podem exigir uma permissão específica. Se você não tiver a permissão necessária, receberá um erro `403 Forbidden`.

## Níveis de Assinatura

O serviço oferece diferentes níveis de assinatura, cada um com um conjunto específico de funcionalidades:

### Free

- **Funcionalidades**:
  - Registrar e fazer login.
  - Visualizar parcelas de subcontas já criadas (apenas leitura).
  - Simular juros com valores limitados (ex: até R$ 10.000,00).
  - Criar até 3 subcontas.
  - Acesso a documentação básica.
- **Limitações**:
  - Limite de 3 subcontas.
  - Limite de R$ 10.000,00 para simulações de juros.

### Tier

- **Funcionalidades**:
  - Todas as features do **Free**.
  - Criar subcontas ilimitadas.
  - Simular juros com valores ilimitados.
  - Tracking de investimentos (criar, visualizar, atualizar).
  - Gestão de dívidas básicas (criar, visualizar).
  - Exportar relatórios em PDF (limite de 5 relatórios por mês).
- **Limitações**:
  - Limite de 5 relatórios PDF por mês.

### Gold

- **Funcionalidades**:
  - Todas as features do **Tier**.
  - Gestão de dívidas avançada (amortizações, projeções).
  - Tracking de investimentos avançado (projeções, comparação).
  - Visualização temporal avançada (gráficos).
  - Exportar relatórios em PDF ilimitados.
  - Suporte prioritário.
- **Limitações**:
  - Nenhuma limitação explícita mencionada.

### Platinum

- **Funcionalidades**:
  - Todas as features do **Gold**.
  - Acesso a modelos de IA para análise de dados (ex: previsão de gastos, sugestões de investimento).
  - Integração com calendário (sincronização de vencimentos).
  - Relatórios personalizados.
  - Suporte 24/7.
- **Limitações**:
  - Nenhuma limitação explícita mencionada.

## Endpoints

### Health Check

- **URL**: `/health`
- **Método**: `GET`
- **Descrição**: Verifica se o serviço está funcionando.
- **Autenticação**: Não requerida
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "status": "ok",
      "message": "Quanto Eu Vou Pagar? service is up and running"
    }
    ```

### Autenticação

#### Registrar Usuário

- **URL**: `/auth/register`
- **Método**: `POST`
- **Descrição**: Registra um novo usuário no sistema.
- **Autenticação**: Não requerida
- **Dados necessários**:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string"
  }
  ```
- **Resposta de Sucesso**:
  - **Código**: `201 Created`
  - **Conteúdo**:
    ```json
    {
      "message": "Usuário registrado com sucesso",
      "user": {
        "id": "integer",
        "username": "string",
        "email": "string",
        "role": "string",
        "subscription_level": "string",
        "created_at": "string (ISO 8601)",
        "updated_at": "string (ISO 8601)"
      }
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

#### Login de Usuário

- **URL**: `/auth/login`
- **Método**: `POST`
- **Descrição**: Autentica um usuário e retorna um token JWT.
- **Autenticação**: Não requerida
- **Dados necessários**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "token": "string"
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

### Subcontas

#### Criar Subconta

- **URL**: `/api/subcontas`
- **Método**: `POST`
- **Descrição**: Cria uma nova subconta e calcula suas parcelas.
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `user` (qualquer usuário autenticado)
- **Limitações por Nível**:
  - **Free**: Limite de 3 subcontas e valor máximo de R$ 10.000,00.
  - **Tier, Gold, Platinum**: Sem limites.
- **Dados necessários**:
  ```json
  {
    "descricao": "string",
    "valor": "number (float)",
    "taxa_juros": "number (float) - Taxa de juros mensal",
    "parcelas": "integer",
    "data_inicio": "string (ISO 8601)",
    "sistema": "string - 'price' ou 'sac'",
    "tipo": "string"
  }
  ```
- **Resposta de Sucesso**:
  - **Código**: `201 Created`
  - **Conteúdo**:
    ```json
    {
      "message": "Subconta criada com sucesso",
      "subconta": {
        "id": "integer",
        "descricao": "string",
        "valor": "number",
        "taxa_juros": "number",
        "parcelas": "integer",
        "data_inicio": "string (ISO 8601)",
        "sistema": "string",
        "tipo": "string",
        "criado_em": "string (ISO 8601)",
        "atualizado_em": "string (ISO 8601)"
      },
      "parcelas": [
        {
          "id": "integer",
          "subconta_id": "integer",
          "numero": "integer",
          "valor": "number",
          "principal": "number",
          "juros": "number",
          "data_venc": "string (ISO 8601)",
          "pago": "boolean"
        }
        // ... mais parcelas
      ]
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `403 Forbidden`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

#### Obter Parcelas de uma Subconta

- **URL**: `/api/subcontas/:id/parcelas`
- **Método**: `GET`
- **Descrição**: Obtém todas as parcelas associadas a uma subconta específica.
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `user` (qualquer usuário autenticado)
- **Limitações por Nível**:
  - **Free**: Acesso básico.
  - **Tier, Gold, Platinum**: Acesso completo.
- **Parâmetros da URL**:
  - `id` (integer): ID da subconta.
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "parcelas": [
        {
          "id": "integer",
          "subconta_id": "integer",
          "numero": "integer",
          "valor": "number",
          "principal": "number",
          "juros": "number",
          "data_venc": "string (ISO 8601)",
          "pago": "boolean"
        }
        // ... mais parcelas
      ]
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `403 Forbidden`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

### Pagamentos

#### Criar Pagamento

- **URL**: `/api/payments`
- **Método**: `POST`
- **Descrição**: Cria um novo pagamento para o usuário autenticado.
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `user` (qualquer usuário autenticado)
- **Dados necessários**:
  ```json
  {
    "amount": "number (float)",
    "currency": "string (opcional, padrão: BRL)",
    "method": "string (ex: 'credit_card', 'boleto', 'pix')",
    "description": "string (opcional)"
  }
  ```
- **Resposta de Sucesso**:
  - **Código**: `201 Created`
  - **Conteúdo**:
    ```json
    {
      "payment": {
        "id": "integer",
        "user_id": "integer",
        "amount": "number",
        "currency": "string",
        "status": "string",
        "method": "string",
        "transaction_id": "string (opcional)",
        "description": "string (opcional)",
        "created_at": "string (ISO 8601)",
        "updated_at": "string (ISO 8601)"
      }
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

#### Obter Pagamento

- **URL**: `/api/payments/:id`
- **Método**: `GET`
- **Descrição**: Obtém os detalhes de um pagamento específico.
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `user` (qualquer usuário autenticado)
- **Parâmetros da URL**:
  - `id` (integer): ID do pagamento.
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "payment": {
        "id": "integer",
        "user_id": "integer",
        "amount": "number",
        "currency": "string",
        "status": "string",
        "method": "string",
        "transaction_id": "string (opcional)",
        "description": "string (opcional)",
        "created_at": "string (ISO 8601)",
        "updated_at": "string (ISO 8601)"
      }
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `404 Not Found`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

### Assinaturas (Administração)

#### Atualizar Assinatura de Usuário

- **URL**: `/api/admin/users/:id/subscription`
- **Método**: `PUT`
- **Descrição**: Atualiza o nível de assinatura de um usuário (requer permissão de admin).
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `admin`
- **Parâmetros da URL**:
  - `id` (integer): ID do usuário.
- **Dados necessários**:
  ```json
  {
    "level": "string (ex: 'free', 'tier', 'gold', 'platinum')"
  }
  ```
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "message": "Nível de assinatura atualizado com sucesso"
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `403 Forbidden`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `404 Not Found`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

### Pagamentos

#### Criar Pagamento

- **URL**: `/api/payments`
- **Método**: `POST`
- **Descrição**: Cria um novo pagamento para o usuário autenticado.
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `user` (qualquer usuário autenticado)
- **Dados necessários**:
  ```json
  {
    "amount": "number (float)",
    "currency": "string (opcional, padrão: BRL)",
    "method": "string (ex: 'credit_card', 'boleto', 'pix')",
    "description": "string (opcional)"
  }
  ```
- **Resposta de Sucesso**:
  - **Código**: `201 Created`
  - **Conteúdo**:
    ```json
    {
      "payment": {
        "id": "integer",
        "user_id": "integer",
        "amount": "number",
        "currency": "string",
        "status": "string",
        "method": "string",
        "transaction_id": "string (opcional)",
        "description": "string (opcional)",
        "created_at": "string (ISO 8601)",
        "updated_at": "string (ISO 8601)"
      }
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

#### Obter Pagamento

- **URL**: `/api/payments/:id`
- **Método**: `GET`
- **Descrição**: Obtém os detalhes de um pagamento específico.
- **Autenticação**: Requerida (Bearer Token)
- **Permissões**: `user` (qualquer usuário autenticado)
- **Parâmetros da URL**:
  - `id` (integer): ID do pagamento.
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "payment": {
        "id": "integer",
        "user_id": "integer",
        "amount": "number",
        "currency": "string",
        "status": "string",
        "method": "string",
        "transaction_id": "string (opcional)",
        "description": "string (opcional)",
        "created_at": "string (ISO 8601)",
        "updated_at": "string (ISO 8601)"
      }
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `401 Unauthorized`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `404 Not Found`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`

### Obter Parcelas de uma Subconta

- **URL**: `/subcontas/:id/parcelas`
- **Método**: `GET`
- **Descrição**: Obtém todas as parcelas associadas a uma subconta específica.
- **Parâmetros da URL**:
  - `id` (integer): ID da subconta.
- **Resposta de Sucesso**:
  - **Código**: `200 OK`
  - **Conteúdo**:
    ```json
    {
      "parcelas": [
        {
          "id": "integer",
          "subconta_id": "integer",
          "numero": "integer",
          "valor": "number",
          "principal": "number",
          "juros": "number",
          "data_venc": "string (ISO 8601)",
          "pago": "boolean"
        }
        // ... mais parcelas
      ]
    }
    ```
- **Respostas de Erro**:
  - **Código**: `400 Bad Request`
    - **Conteúdo**: `{"error": "mensagem de erro"}`
  - **Código**: `500 Internal Server Error`
    - **Conteúdo**: `{"error": "mensagem de erro"}`