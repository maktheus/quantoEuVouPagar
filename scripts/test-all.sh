#!/bin/bash

# Script para testar todos os serviços do projeto

echo "Iniciando testes..."

# Verificar se o docker-compose está instalado
if ! command -v docker-compose &> /dev/null
then
    echo "docker-compose não encontrado. Por favor, instale o Docker e o Docker Compose."
    exit 1
fi

# Iniciar os serviços
echo "Iniciando serviços com docker-compose..."
docker-compose up -d

# Esperar um pouco para os serviços iniciarem
echo "Aguardando serviços iniciarem..."
sleep 30

# Testar Health Check
echo "Testando Health Check..."
curl -s -o /dev/null -w "%{http_code}" http://localhost/health | grep -q "200"
if [ $? -eq 0 ]; then
    echo "✓ Health Check: OK"
else
    echo "✗ Health Check: Falhou"
fi

# Testar Swagger
echo "Testando Swagger..."
curl -s -o /dev/null -w "%{http_code}" http://localhost/swagger/index.html | grep -q "200"
if [ $? -eq 0 ]; then
    echo "✓ Swagger: OK"
else
    echo "✗ Swagger: Falhou"
fi

# Testar Prometheus
echo "Testando Prometheus..."
curl -s -o /dev/null -w "%{http_code}" http://localhost:9090/graph | grep -q "200"
if [ $? -eq 0 ]; then
    echo "✓ Prometheus: OK"
else
    echo "✗ Prometheus: Falhou"
fi

# Testar Grafana
echo "Testando Grafana..."
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/login | grep -q "200"
if [ $? -eq 0 ]; then
    echo "✓ Grafana: OK"
else
    echo "✗ Grafana: Falhou"
fi

# Testar API - Registro de Usuário
echo "Testando API - Registro de Usuário..."
curl -s -X POST http://localhost/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}' | grep -q '"message":"Usuário registrado com sucesso"'
if [ $? -eq 0 ]; then
    echo "✓ Registro de Usuário: OK"
else
    echo "✗ Registro de Usuário: Falhou"
fi

# Testar API - Login de Usuário
echo "Testando API - Login de Usuário..."
TOKEN=$(curl -s -X POST http://localhost/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo "✓ Login de Usuário: OK"
else
    echo "✗ Login de Usuário: Falhou"
fi

# Testar API - Criar Subconta
if [ -n "$TOKEN" ]; then
    echo "Testando API - Criar Subconta..."
    curl -s -X POST http://localhost/api/subcontas \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer $TOKEN" \
         -d '{
           "descricao": "Financiamento casa",
           "valor": 300000.0,
           "taxa_juros": 0.01,
           "parcelas": 120,
           "data_inicio": "2024-01-01T00:00:00Z",
           "sistema": "price",
           "tipo": "financiamento"
         }' | grep -q '"message":"Subconta criada com sucesso"'
    if [ $? -eq 0 ]; then
        echo "✓ Criar Subconta: OK"
    else
        echo "✗ Criar Subconta: Falhou"
    fi
fi

# Testar API - Obter Parcelas
if [ -n "$TOKEN" ]; then
    echo "Testando API - Obter Parcelas..."
    curl -s http://localhost/api/subcontas/1/parcelas \
         -H "Authorization: Bearer $TOKEN" | grep -q '"parcelas"'
    if [ $? -eq 0 ]; then
        echo "✓ Obter Parcelas: OK"
    else
        echo "✗ Obter Parcelas: Falhou"
    fi
fi

# Parar os serviços
echo "Parando serviços..."
docker-compose down

echo "Testes finalizados."