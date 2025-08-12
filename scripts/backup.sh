#!/bin/bash

# Script para fazer backup do banco de dados PostgreSQL

# Configurações
DB_HOST="db" # Nome do serviço do banco de dados no docker-compose
DB_PORT="5432"
DB_NAME="quanto_eu_vou_pagar"
DB_USER="quanto_eu_vou_pagar_user"
DB_PASSWORD="strong_password_here" # Substitua pela senha real

# Diretório para salvar os backups
BACKUP_DIR="/backups"
DATE=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql"

# Criar diretório de backups se não existir
mkdir -p $BACKUP_DIR

# Fazer backup
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME > $BACKUP_FILE

if [ $? -eq 0 ]; then
  echo "Backup realizado com sucesso: $BACKUP_FILE"
else
  echo "Erro ao realizar backup"
  exit 1
fi

# Manter apenas os últimos 7 backups
find $BACKUP_DIR -type f -name "backup_*.sql" -mtime +7 -delete