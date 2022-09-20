#!/bin/sh
PGPASSWORD=postgres psql -wq -h localhost -p 5437 -U postgres -d devbox -f ./remove_test_migrations.sql