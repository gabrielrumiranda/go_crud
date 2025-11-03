FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download
COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts ./scripts

RUN chmod +x ./scripts/init-db.sh 2>/dev/null || true

EXPOSE 5001

CMD ["./main"]
