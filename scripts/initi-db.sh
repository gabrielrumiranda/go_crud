ash#!/bin/sh
# Script para executar migrations manualmente se necessário

set -e

echo "🔧 Aguardando PostgreSQL estar pronto..."
until PGPASSWORD=$DB_PASSWORD psql -h "postgres" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "PostgreSQL não está pronto - aguardando..."
  sleep 1
done

echo "✅ PostgreSQL está pronto!"

echo "📝 Executando migrations..."
for migration in ./migrations/*.sql; do
  if [ -f "$migration" ]; then
    echo "Executando: $migration"
    PGPASSWORD=$DB_PASSWORD psql -h "postgres" -U "$DB_USER" -d "$DB_NAME" -f "$migration"
  fi
done

echo "✅ Migrations executadas com sucesso!"
