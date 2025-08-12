-- Criação da tabela users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Armazena o hash da senha
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- 'user' ou 'admin'
    subscription_level VARCHAR(20) NOT NULL DEFAULT 'free', -- 'free', 'tier', 'gold', 'platinum'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela subcontas
CREATE TABLE IF NOT EXISTS subcontas (
    id SERIAL PRIMARY KEY,
    descricao TEXT NOT NULL,
    valor DECIMAL(15, 2) NOT NULL,
    taxa_juros DECIMAL(10, 6) NOT NULL, -- Taxa de juros mensal
    parcelas INTEGER NOT NULL,
    data_inicio DATE NOT NULL,
    sistema VARCHAR(10) NOT NULL, -- 'price' ou 'sac'
    tipo VARCHAR(50) NOT NULL, -- Ex: 'financiamento', 'boleto'
    criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela parcelas
CREATE TABLE IF NOT EXISTS parcelas (
    id SERIAL PRIMARY KEY,
    subconta_id INTEGER NOT NULL REFERENCES subcontas(id) ON DELETE CASCADE,
    numero INTEGER NOT NULL, -- Número da parcela (1, 2, 3, ...)
    valor DECIMAL(15, 2) NOT NULL, -- Valor total da parcela (principal + juros)
    principal DECIMAL(15, 2) NOT NULL, -- Parte do valor que amortiza a dívida
    juros DECIMAL(15, 2) NOT NULL, -- Parte do valor referente aos juros
    data_venc DATE NOT NULL, -- Data de vencimento da parcela
    pago BOOLEAN NOT NULL DEFAULT FALSE, -- Indica se a parcela foi paga
    criado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela payments
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'completed', 'failed', 'cancelled'
    method VARCHAR(20) NOT NULL, -- 'credit_card', 'boleto', 'pix'
    transaction_id VARCHAR(255), -- ID da transação no gateway de pagamento
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela plans (para assinaturas)
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    interval VARCHAR(20) NOT NULL, -- 'monthly', 'yearly'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela subscriptions (para assinaturas)
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'inactive', 'cancelled'
    start_date DATE NOT NULL,
    end_date DATE, -- Pode ser NULL para assinaturas contínuas
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela investments (para tracking de investimentos)
CREATE TABLE IF NOT EXISTS investments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    initial_value DECIMAL(15, 2) NOT NULL,
    current_value DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela debts (para gestão de dívidas)
CREATE TABLE IF NOT EXISTS debts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    principal DECIMAL(15, 2) NOT NULL,
    interest_rate DECIMAL(10, 6) NOT NULL, -- Taxa de juros mensal
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);