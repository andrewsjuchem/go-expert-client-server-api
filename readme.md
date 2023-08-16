# Observações

- Embora existam arquivos de Docker neste projeto, não é necessário utilizar o Docker para executar o programa. Basta executar os comandos listados abaixo.
- O banco de dados SQLite será criado automaticamente em ambiante local.


**Executar Servidor**

```
go run ./server/main.go
```

**Chamar Endpoint**

```
curl -X GET http://localhost:8080/cotacao
```

**Executar Cliente**

```
go run ./client/main.go
```