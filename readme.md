**Observações**

- Este código foi testado em MacOC (M1) e Windows 10.
- Embora existam arquivos de Docker neste projeto, não é necessário utilizar o Docker para executar o programa. Basta executar os comandos listados abaixo.
- O banco de dados SQLite será criado automaticamente em ambiante local após a primeira chamada que for completada com sucesso.

**Observações Importantes**

- Importante lembrar que conforme o enunciado, o timeout máximo para conseguir persistir os dados no banco é de 10ms. No entanto, em algumas vezes, 10ms não é suficiente para persistir os dados no banco. Por exemplo, abaixo o endpoint precisou ser chamado 3 vezes até que a operação completou em menos de 10ms. Esse endpoint é chamado pelo cliente, então é possível que o cliente precise ser chamado mais de uma vez.
```
C:\Users\andrewsjuchem>curl -X GET http://localhost:8080/cotacao
curl: (52) Empty reply from server

C:\Users\andrewsjuchem>curl -X GET http://localhost:8080/cotacao
curl: (52) Empty reply from server

C:\Users\andrewsjuchem>curl -X GET http://localhost:8080/cotacao
{"bid":"4.957"}
```
- O biblioteca do Go SQLite precisa que a opção CGO esteja ativada ("go-sqlite3 requires cgo to work"). Nesse caso, certifique-se que a variável CGO_ENABLED está preenchida com "1" na hora de executar o servidor.
- Quando o CGO está ativado, é necessário ter um compilador C instalado no ambiente. Caso você não tenha, os comandos abaixos mostram como instalar.
Para Linux, execute o comando abaixo:
```
apt-get install build-essential
```
Para Windows, instale o choco (https://chocolatey.org/install) e execute o comando abaixo.
```
choco install mingw -y
```

**Executar Servidor**

```
export CGO_ENABLED=1 (set CGO_ENABLED=1 para windows)
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