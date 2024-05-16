FROM postgres:latest as psqldb

COPY start.sql /docker-entrypoint-initdb.d/
