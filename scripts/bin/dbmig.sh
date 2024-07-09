#!/bin/zsh

goose -allow-missing -dir ../dbmigrations postgres "postgres://$DATABASE_USERNAME:$DATABASE_PASSWORD@$DATABASE_URL:5432/$DATABASE_NAME?sslmode=disable" $*
