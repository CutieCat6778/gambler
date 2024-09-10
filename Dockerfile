FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gambler

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gambler .

COPY .env .env

EXPOSE 8080

CMD ["sh", "-c", "export $(cat .env | xargs) && ./gambler"]
