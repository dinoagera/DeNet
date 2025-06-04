FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY .env . 
COPY . .
COPY ./migrations ./migrations

RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.22/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache curl unzip && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/bin/migrate && \
    chmod +x /usr/bin/migrate

RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /usr/bin/migrate /usr/bin/migrate
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

CMD ["sh", "-c", "migrate -path /app/migrations -database \"$DB_URL\" up && ./main"]
