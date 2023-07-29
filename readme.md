# Setting Up

```
docker-compose up -d --build
```

# Checking Database

```
docker-compose exec sqlite3 bash
sqlite3 currency_exchange.db
select * from quote;
```

# Running Server

```
docker-compose exec goapp bash --UPDATE LATER
go mod tidy
go run cmd/server.go
```

# Running Client

```
docker-compose exec goapp bash --UPDATE LATER
go mod tidy
go run cmd/client.go
```