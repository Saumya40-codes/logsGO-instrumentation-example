FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:3.22 AS final
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]

LABEL rebuilt="true"
