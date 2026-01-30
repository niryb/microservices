FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o shipping cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/shipping .

EXPOSE 3002

CMD ["./shipping"]