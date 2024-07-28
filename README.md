# go-finance-api


### How to create a Postgres database in docker
Use the following command in a terminal
```bash
docker run --name go-finance -e POSTGRES_DB=gofinance -e POSTGRES_USER=finance -e POSTGRES_PASSWORD=finance -d -p 5432:5432 postgres
```