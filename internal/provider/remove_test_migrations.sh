#!/bin/sh
PGPASSWORD=secret psql -wq -h localhost -p 5437 -U aidbox -d aidbox -f ./remove_test_migrations.sql