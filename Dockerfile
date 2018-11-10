FROM golang:stretch as Builder

WORKDIR /app

COPY /src ./src

COPY vendor $HOME/go/src

WORKDIR /app/src/main

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o kubego main.go

FROM alpine

COPY --from=Builder /app/src/main/kubego /app/kubego

WORKDIR /app

ENTRYPOINT ["./kubego"]