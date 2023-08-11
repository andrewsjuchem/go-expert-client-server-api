# Checking Database

```
sqlite3 currency_exchange.db
select * from quote;
```

# Running Server (Docker)

```
docker-compose -f docker-compose-server.yml up -d --build
or
docker-compose -f docker-compose-server.yml up
```

# Stopping Server (Docker)

```
docker-compose -f docker-compose-server.yml down
```

# Calling the Server's Endpoint

```
curl -X GET http://localhost:8080/cotacao
```

# Running Client (Docker)

```
docker-compose exec goapp bash --UPDATE LATER
go mod tidy
go run cmd/client.go
```