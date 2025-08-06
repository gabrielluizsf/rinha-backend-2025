FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

WORKDIR /app
RUN go build -o /rinha_api

FROM alpine:latest
COPY --from=builder /rinha_api /rinha_api

EXPOSE 8080

CMD ["/rinha_api"]
