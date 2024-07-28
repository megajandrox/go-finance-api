# Usa la imagen oficial de PostgreSQL
FROM postgres:latest

# Establece las variables de entorno necesarias para PostgreSQL
ENV POSTGRES_DB=yourdb
ENV POSTGRES_USER=youruser
ENV POSTGRES_PASSWORD=yourpassword

# Copia cualquier script SQL en la carpeta de inicialización
# COPY init.sql /docker-entrypoint-initdb.d/
