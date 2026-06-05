FROM golang:1.26-alpine AS builder

RUN apk add --no-cache build-base ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN touch .env
RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cmd/app
RUN CGO_ENABLED=1 GOOS=linux go build -o migrate ./cmd/migrate

FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrate .

COPY ./migrations ./migrations

COPY --from=builder /app/.env .

COPY ./scripts/entrypoint.sh .

RUN chmod +x entrypoint.sh