FROM postgres:12

# Open the database and create the UUID extension.
RUN echo "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";" > /docker-entrypoint-initdb.d/uuid-ossp.sql

COPY ./compose/production/postgres/maintenance /usr/local/bin/maintenance
RUN chmod -x /usr/local/bin/maintenance/*
RUN mv /usr/local/bin/maintenance/* /usr/local/bin \
    && rm -rf /usr/local/bin/maintenance
