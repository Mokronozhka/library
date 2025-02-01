FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# COPY wait-for-it.sh ./wait-for-it.sh
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o library cmd/library/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/library .
COPY --from=builder /app/library wait-for-it.sh
# COPY --from=builder /app/library/wait-for-it.sh ./wait-for-it.sh
RUN chmod +x wait-for-it.sh
CMD ["./wait-for-it.sh", "db:5432", "--timeout=30", "./library", "-debug"]
EXPOSE 8080