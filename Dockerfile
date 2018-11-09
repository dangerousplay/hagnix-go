FROM golang:stretch

WORKDIR /app

COPY main.go .

RUN go get k8s.io/client-go/...

RUN env GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o kubego.so main.go