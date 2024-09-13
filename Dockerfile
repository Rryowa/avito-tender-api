FROM golang:alpine AS builder
RUN apk --no-cache add bash git make curl tar
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh
RUN /usr/local/bin/goose -version || echo "Goose not found"

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main ./cmd/main.go

ENTRYPOINT ["./bin/main"]