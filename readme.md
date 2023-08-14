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

# Running Server (Local)

```
go run ./server/main.go
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
docker-compose -f docker-compose-client.yml up -d --build
or
docker-compose -f docker-compose-client.yml up
```

# Running Client (Local)

```
go run ./client/main.go
```