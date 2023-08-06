# First stage: build the application
FROM golang:1.20.4 AS builder
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=1 go build -o server-app -ldflags="-w -s" server/main.go
CMD [ "/app/server-app" ]

# NOTE: CGO has to be enabled for sqlite3 to work, but the scratch container does not work well with CGO enabled
#       because it results in dynamic links to libc/libmusl. For example, if we build the Go binary with
#       the command below, the binary will run but will fail the sqlite3 command.
# RUN CGO_ENABLED=0 go build -o server-app -ldflags="-w -s" server/main.go

# Second stage: create a minimal runtime image
# FROM scratch
# WORKDIR /app
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /app/server-app .
# CMD [ "/app/server-app" ]