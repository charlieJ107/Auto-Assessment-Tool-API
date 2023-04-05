FROM golang:1.19.3 AS builder
LABEL authors="Charlie"
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
ENV GIN_MODE=release
EXPOSE 80
CMD ["./main"]