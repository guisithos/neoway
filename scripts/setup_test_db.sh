#!/bin/bash
set -e

# Para o container de teste
docker-compose -f docker-compose.test.yml down -v || true

# Inicia db
docker-compose -f docker-compose.test.yml up -d

# Aguarda um pouco
echo "Aguardando db pronta..."
for i in {1..30}; do
    if docker-compose -f docker-compose.test.yml exec -T postgres_test pg_isready -U postgres; then
        echo "db pronta!"
        break
    fi
    echo "Aguardando postgres... (tentativa $i/30)"
    sleep 2
    if [ $i -eq 30 ]; then
        echo "Timeout aguardando postgres"
        exit 1
    fi
done

# Cria se não exist db
docker-compose -f docker-compose.test.yml exec -T postgres_test psql -U postgres -c "CREATE DATABASE neoway_test;" || true

# Exporta as env 
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres123
export TEST_DB_NAME=neoway_test

# Printa as env para debug
echo "Test database configuration:"
echo "Host: $TEST_DB_HOST"
echo "Port: $TEST_DB_PORT"
echo "User: $TEST_DB_USER"
echo "Database: $TEST_DB_NAME"

echo "ambiente de teste pronto!" 