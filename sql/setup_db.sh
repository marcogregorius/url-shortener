#!/bin/bash
PGPASSWORD=${DB_PASSWORD} psql --host ${DB_HOST} --port ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -f create_table.sql
